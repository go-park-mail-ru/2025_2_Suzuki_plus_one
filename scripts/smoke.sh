#!/usr/bin/env bash

# PopFilms API smoke test
# BASE_URL example:
#   export BASE_URL="http://localhost:8080"
#   # or for prod:
#   # export BASE_URL="https://popfilms.ru/api/v1"
BASE_URL="${BASE_URL:-https://popfilms.ru/api/v1}"
# BASE_URL="${BASE_URL:-http://localhost:8080}"

# Credentials for /auth/signin
EMAIL="${EMAIL:-test@example.com}"
PASSWORD="${PASSWORD:-Password123!}"

# IDs used in path params – adjust to match your data
MEDIA_ID="${MEDIA_ID:-1}"
GENRE_ID="${GENRE_ID:-1}"
ACTOR_ID="${ACTOR_ID:-1}"
APPEAL_ID="${APPEAL_ID:-1}"

# Optional: avatar file for /user/me/update/avatar
# e.g. export AVATAR_FILE="./avatar.jpg"
AVATAR_FILE="${AVATAR_FILE:-}"

COOKIES_FILE="cookies.txt"

########################################
# Colors
########################################

GREEN="\033[0;32m"
RED="\033[0;31m"
BLUE="\033[0;34m"
YELLOW="\033[0;33m"
NC="\033[0m"

########################################
# Helpers
########################################

print_result() {
  local name="$1"
  local code="$2"

  if [[ "$code" =~ ^[0-9]+$ ]] && [[ "$code" -ge 200 && "$code" -lt 300 ]]; then
    printf "%b[OK]%b   %s (%s)\n" "$GREEN" "$NC" "$name" "$code"
  else
    printf "%b[FAIL]%b %s (%s)\n" "$RED" "$NC" "$name" "$code"
  fi
}

call_endpoint() {
  local method="$1"
  local path="$2"
  local use_token="$3"   # 0 or 1
  local label="$4"

  local auth_header=()
  if [[ "$use_token" -eq 1 && -n "${ACCESS_TOKEN:-}" ]]; then
    auth_header=("-H" "Authorization: Bearer $ACCESS_TOKEN")
  fi

  local code
  code=$(curl -s -o /dev/null -w "%{http_code}" \
    -X "$method" \
    "${auth_header[@]}" \
    "$BASE_URL$path" || echo "000")

  local suffix="no-auth"
  if [[ "$use_token" -eq 1 ]]; then
    suffix="with-auth"
  fi

  print_result "$method $path [$suffix] - $label" "$code"
}

call_endpoint_with_body() {
  local method="$1"
  local path="$2"
  local use_token="$3"   # 0 or 1
  local label="$4"
  local body="$5"

  local auth_header=()
  if [[ "$use_token" -eq 1 && -n "${ACCESS_TOKEN:-}" ]]; then
    auth_header=("-H" "Authorization: Bearer $ACCESS_TOKEN")
  fi

  local code
  code=$(curl -s -o /dev/null -w "%{http_code}" \
    -X "$method" \
    -H "Content-Type: application/json" \
    "${auth_header[@]}" \
    -d "$body" \
    "$BASE_URL$path" || echo "000")

  local suffix="no-auth"
  if [[ "$use_token" -eq 1 ]]; then
    suffix="with-auth"
  fi

  print_result "$method $path [$suffix] - $label" "$code"
}

test_both() {
  local method="$1"
  local path="$2"
  local label="$3"
  call_endpoint "$method" "$path" 0 "$label"
  call_endpoint "$method" "$path" 1 "$label"
}

test_both_with_body() {
  local method="$1"
  local path="$2"
  local label="$3"
  local body="$4"
  call_endpoint_with_body "$method" "$path" 0 "$label" "$body"
  call_endpoint_with_body "$method" "$path" 1 "$label" "$body"
}

########################################
# Sign in to get token + cookies
########################################

echo -e "${BLUE}=== Signing in at /auth/signin ===${NC}"
signin_response=$(curl -s -c "$COOKIES_FILE" -X POST "$BASE_URL/auth/signin" \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" || true)

if [[ -z "$signin_response" ]]; then
  echo -e "${RED}Signin request failed (no response).${NC}"
else
  if command -v jq >/dev/null 2>&1; then
    ACCESS_TOKEN=$(echo "$signin_response" | jq -r '.access_token // .token // empty')
    if [[ -n "$ACCESS_TOKEN" && "$ACCESS_TOKEN" != "null" ]]; then
      echo -e "${GREEN}Access token obtained.${NC}"
    else
      echo -e "${YELLOW}Signin response received, but token not found in JSON.${NC}"
    fi
  else
    echo -e "${YELLOW}jq not installed; cannot parse token from signin response.${NC}"
    ACCESS_TOKEN=""
  fi
fi

########################################
# Auth
########################################

echo -e "\n${BLUE}=== Auth ===${NC}"

# /auth/refresh (needs refresh_token cookie)
code=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/auth/refresh" || echo "000")
print_result "GET /auth/refresh [no-cookie]" "$code"

code=$(curl -s -b "$COOKIES_FILE" -o /dev/null -w "%{http_code}" "$BASE_URL/auth/refresh" || echo "000")
print_result "GET /auth/refresh [with-cookie]" "$code"

# /auth/signup - use random email so you can rerun the script
if command -v date >/dev/null 2>&1; then
  RAND_SUFFIX=$(date +%s)
else
  RAND_SUFFIX=$RANDOM
fi
SIGNUP_EMAIL="user_${RAND_SUFFIX}@example.com"
SIGNUP_BODY=$(cat <<EOF
{"email":"$SIGNUP_EMAIL","password":"securepassword123","username":"user_$RAND_SUFFIX"}
EOF
)
test_both_with_body "POST" "/auth/signup" "PostAuthSignUp" "$SIGNUP_BODY"

# /auth/signout – protected
test_both "GET" "/auth/signout" "GetAuthSignOut"

########################################
# Object
########################################

echo -e "\n${BLUE}=== Object ===${NC}"
# /object?key=...&bucket_name=...
test_both "GET" "/object?key=test_key&bucket_name=posters" "GetObjectMedia"

########################################
# Content: Media
########################################

echo -e "\n${BLUE}=== Content: Media ===${NC}"

# /media/my?is_dislike=...&limit=...&offset=...
test_both "GET" "/media/my?is_dislike=false&limit=10&offset=0" "GetMediaMy"

# /media/recommendations?type=movie&limit=3&offset=0
test_both "GET" "/media/recommendations?type=movie&limit=3&offset=0" "GetMediaRecommendations"

# /media/{media_id}
test_both "GET" "/media/${MEDIA_ID}" "GetMedia"

# /media/{media_id}/actor
test_both "GET" "/media/${MEDIA_ID}/actor" "GetMediaActor"

# /media/watch?media_id=...
test_both "GET" "/media/watch?media_id=${MEDIA_ID}" "GetMediaWatch"

# /media/{media_id}/like (GET/PUT/DELETE) – secured
test_both "GET" "/media/${MEDIA_ID}/like" "GetMediaLike"
test_both "PUT" "/media/${MEDIA_ID}/like" "PutMediaLike"
test_both "DELETE" "/media/${MEDIA_ID}/like" "DeleteMediaLike"

########################################
# Content: Actor & Search
########################################

echo -e "\n${BLUE}=== Content: Actor & Search ===${NC}"

# /actor/{actor_id}
test_both "GET" "/actor/${ACTOR_ID}" "GetActor"

# /actor/{actor_id}/media
test_both "GET" "/actor/${ACTOR_ID}/media" "GetActorMedia"

# /search?query=...&type=any&limit=...&offset=...
test_both "GET" "/search?query=test&type=any&limit=5&offset=0" "GetSearch"

########################################
# Genre
########################################

echo -e "\n${BLUE}=== Genre ===${NC}"

# /genre/all
test_both "GET" "/genre/all" "GetGenreAll"

# /genre/{genre_id}?media_limit=...&media_offset=...
test_both "GET" "/genre/${GENRE_ID}?media_limit=5&media_offset=0" "GetGenre"

########################################
# User
########################################

echo -e "\n${BLUE}=== User ===${NC}"

# /user/me
test_both "GET" "/user/me" "GetUserMe"

# /user/me/update (JSON body)
USER_UPDATE_BODY='{
  "username": "johndoe",
  "email": "test@example.com",
  "date_of_birth": "1990-01-01",
  "phone_number": "+1234567890"
}'
test_both_with_body "POST" "/user/me/update" "PostUserMeUpdate" "$USER_UPDATE_BODY"

# /user/me/update/avatar (multipart/form-data)
if [[ -n "$AVATAR_FILE" && -f "$AVATAR_FILE" ]]; then
  for use_token in 0 1; do
    auth_header=()
    if [[ "$use_token" -eq 1 && -n "${ACCESS_TOKEN:-}" ]]; then
      auth_header=(-H "Authorization: Bearer $ACCESS_TOKEN")
    fi

    code=$(curl -s -o /dev/null -w "%{http_code}" \
      -X POST \
      "${auth_header[@]}" \
      -F "avatar=@${AVATAR_FILE}" \
      "$BASE_URL/user/me/update/avatar" || echo "000")

    suffix="no-auth"
    [[ "$use_token" -eq 1 ]] && suffix="with-auth"
    print_result "POST /user/me/update/avatar [$suffix] - PostUserMeUpdateAvatar" "$code"
  done
else
  echo -e "${YELLOW}Skipping /user/me/update/avatar (AVATAR_FILE not set or not found).${NC}"
fi

# /user/me/update/password
USER_PASSWORD_UPDATE_BODY='{
  "current_password": "oldpassword123",
  "new_password": "newsecurepassword456"
}'
test_both_with_body "POST" "/user/me/update/password" "PostUserMeUpdatePassword" "$USER_PASSWORD_UPDATE_BODY"

########################################
# Appeal
########################################

echo -e "\n${BLUE}=== Appeal ===${NC}"

# /appeal/all?limit=...&offset=...&tag=&status=
test_both "GET" "/appeal/all?limit=10&offset=0&tag=bug&status=open" "GetAppealAll"

# /appeal/my (secured)
test_both "GET" "/appeal/my" "GetAppealMy"

# /appeal/new (secured)
APPEAL_NEW_BODY='{
  "tag": "bug",
  "name": "Issue title",
  "message": "Detailed description of the issue..."
}'
test_both_with_body "POST" "/appeal/new" "PostAppealNew" "$APPEAL_NEW_BODY"

# /appeal/{appeal_id}
test_both "GET" "/appeal/${APPEAL_ID}" "GetAppeal"

# /appeal/{appeal_id}/resolve (no body defined in OpenAPI)
test_both "PUT" "/appeal/${APPEAL_ID}/resolve" "PutAppealResolve"

# /appeal/{appeal_id}/message (POST + GET)
APPEAL_MESSAGE_BODY='{
  "message": "Thank you for your feedback!"
}'
test_both_with_body "POST" "/appeal/${APPEAL_ID}/message" "PostAppealMessage" "$APPEAL_MESSAGE_BODY"
test_both "GET" "/appeal/${APPEAL_ID}/message" "GetAppealMessage"

echo -e "\n${GREEN}Done.${NC} Adjust IDs, credentials, and AVATAR_FILE as needed and rerun."
