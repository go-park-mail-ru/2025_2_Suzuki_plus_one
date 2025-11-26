#!/usr/bin/env bash
set -euo pipefail

#########################
# Configurable settings #
#########################

BASE_URL="${BASE_URL:-https://popfilms.ru/api/v1}"
# BASE_URL="${BASE_URL:-http://localhost:8080}"
EMAIL="${EMAIL:-test@example.com}"
PASSWORD="${PASSWORD:-Password123!}"

# How hard to hit each endpoint
TOTAL_PUBLIC="${TOTAL_PUBLIC:-4000}"
CONC_PUBLIC="${CONC_PUBLIC:-100}"

TOTAL_AUTH="${TOTAL_AUTH:-4000}"
CONC_AUTH="${CONC_AUTH:-100}"

# IDs used in paths. Adjust to real IDs in your DB if you want fewer 404s.
MEDIA_ID="${MEDIA_ID:-1}"
APPEAL_ID="${APPEAL_ID:-1}"
ACTOR_ID="${ACTOR_ID:-1}"
GENRE_ID="${GENRE_ID:-1}"

##################################
# Check dependencies: hey & jq   #
##################################

for cmd in hey jq curl; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "ERROR: '$cmd' not found in PATH" >&2
    exit 1
  fi
done

##################################
# Step 1: Login to get JWT       #
##################################

echo "==> Logging in to $BASE_URL/auth/signin as $EMAIL"

LOGIN_RESPONSE="$(curl -sS -X POST \
  "$BASE_URL/auth/signin" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")"

echo "Login response: $LOGIN_RESPONSE"

ACCESS_TOKEN="$(echo "$LOGIN_RESPONSE" | jq -r '.access_token // empty')"

if [[ -z "$ACCESS_TOKEN" || "$ACCESS_TOKEN" == "null" ]]; then
  echo "ERROR: could not extract access_token from response" >&2
  exit 1
fi

AUTH_HEADER="Authorization: Bearer $ACCESS_TOKEN"
echo "==> Got access token. Starting load..."

##################################
# Helper functions               #
##################################

bomb_get_public() {
  local path="$1"
  echo ">> [PUBLIC] GET $path (n=$TOTAL_PUBLIC, c=$CONC_PUBLIC)"
  hey -n "$TOTAL_PUBLIC" -c "$CONC_PUBLIC" "$BASE_URL$path"
}

bomb_get_auth() {
  local path="$1"
  echo ">> [AUTH] GET $path (n=$TOTAL_AUTH, c=$CONC_AUTH)"
  hey -n "$TOTAL_AUTH" -c "$CONC_AUTH" \
    -H "$AUTH_HEADER" \
    "$BASE_URL$path"
}

bomb_write_auth() {
  local method="$1"
  local path="$2"
  local data="${3:-}"
  local total="${4:-$TOTAL_AUTH}"

  echo ">> [AUTH] $method $path (n=$total, c=$CONC_AUTH)"
  if [[ -n "$data" ]]; then
    hey -n "$total" -c "$CONC_AUTH" \
      -m "$method" \
      -T "application/json" \
      -H "$AUTH_HEADER" \
      -d "$data" \
      "$BASE_URL$path"
  else
    hey -n "$total" -c "$CONC_AUTH" \
      -m "$method" \
      -H "$AUTH_HEADER" \
      "$BASE_URL$path"
  fi
}

bomb_signin() {
  # Some extra load on signin endpoint itself
  local total="${1:-2000}"
  local conc="${2:-100}"
  echo ">> [PUBLIC] POST /auth/signin (n=$total, c=$conc)"
  hey -n "$total" -c "$conc" \
    -m POST \
    -T "application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" \
    "$BASE_URL/auth/signin"
}

##################################
# Step 2: Fire the cannons       #
##################################

# Run several hey jobs in parallel to really stress the service.
# Adjust TOTAL_* and CONC_* at the top if you want more/less pain.

# Public endpoints
bomb_get_public "/media/recommendations" &
bomb_get_public "/genre/all" &
bomb_get_public "/genre/$GENRE_ID" &
bomb_get_public "/actor/$ACTOR_ID" &
bomb_get_public "/actor/$ACTOR_ID/media" &
bomb_get_public "/search?query=test" &
bomb_get_public "/object" &
bomb_get_public "/media/watch?media_id=$MEDIA_ID" &
bomb_get_public "/appeal/all" &

# Auth endpoints
bomb_get_auth "/user/me" &
bomb_get_auth "/appeal/my" &
bomb_get_auth "/media/$MEDIA_ID/like" &
bomb_get_auth "/appeal/$APPEAL_ID" &
bomb_get_auth "/appeal/$APPEAL_ID/message" &

# Auth writes (keep n a bit smaller to not totally trash the DB)
bomb_write_auth "PUT" "/media/$MEDIA_ID/like" "" 1000 &
bomb_write_auth "DELETE" "/media/$MEDIA_ID/like" "" 1000 &
bomb_write_auth "POST" "/appeal/new" '{"title":"Load test appeal","description":"Just testing load"}' 500 &
bomb_write_auth "POST" "/appeal/'"$APPEAL_ID"'/message" '{"message":"Hello from load test"}' 500 &

# Some signin load too
bomb_signin 1500 150 &

echo "==> All hey jobs started. Waiting for them to finish..."
wait
echo "==> Load test finished."
