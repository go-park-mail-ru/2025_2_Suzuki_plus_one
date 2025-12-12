#!/usr/bin/env bash

# Simple sqlmap test runner for PopFilms API based on the OpenAPI spec.
#
# USAGE EXAMPLES:
#   chmod +x test_sqlmap_popfilms.sh
#   ./test_sqlmap_popfilms.sh
#
#   BASE_URL="http://localhost:8080/api/v1" ./test_sqlmap_popfilms.sh
#
#   BASE_URL="https://popfilms.ru/api/v1" ACCESS_TOKEN="YOUR_JWT_HERE" ./test_sqlmap_popfilms.sh
#
# NOTE:
# - Only run this against systems you own / are allowed to test.
# - Script is intentionally conservative: risk=1, level=1, small delay.

set -u  # treat unset vars as error, but keep running on sqlmap failures via "|| true"

# Default to production API, can be overridden:
BASE_URL="${BASE_URL:-https://popfilms.ru/api/v1}"

# You can override common sqlmap options via SQLMAP_COMMON_OPTS if you want.
SQLMAP_COMMON_OPTS="${SQLMAP_COMMON_OPTS:---batch --risk=1 --level=1 --random-agent --delay=0.5}"

# Optional: JWT access token for authenticated endpoints.
ACCESS_TOKEN="${ACCESS_TOKEN:-}"

if ! command -v sqlmap >/dev/null 2>&1; then
  echo "[-] sqlmap is not installed or not in PATH" >&2
  exit 1
fi

auth_headers() {
  if [[ -n "$ACCESS_TOKEN" ]]; then
    echo "Authorization: Bearer $ACCESS_TOKEN"
  else
    echo ""
  fi
}

test_get() {
  local url="$1"
  local desc="$2"
  local auth_required="${3:-no}"

  echo
  echo "=================================================================="
  echo "[*] Testing: $desc"
  echo "    URL: $url"
  echo "    Auth required: $auth_required"

  local headers=""
  if [[ "$auth_required" == "yes" ]]; then
    if [[ -z "$ACCESS_TOKEN" ]]; then
      echo "    Skipping (ACCESS_TOKEN not set)"
      return
    fi
    headers="$(auth_headers)"
  fi

  if [[ -n "$headers" ]]; then
    sqlmap -u "$url" $SQLMAP_COMMON_OPTS --headers="$headers" || true
  else
    sqlmap -u "$url" $SQLMAP_COMMON_OPTS || true
  fi
}

echo "[+] Using BASE_URL: $BASE_URL"
if [[ -n "$ACCESS_TOKEN" ]]; then
  echo "[+] ACCESS_TOKEN is set; authenticated endpoints will be tested."
else
  echo "[!] ACCESS_TOKEN not set; authenticated endpoints will be skipped."
fi

# -------------------------------------------------------------------
# PUBLIC ENDPOINTS (no auth)
# -------------------------------------------------------------------

# /object?key=&bucket_name=
test_get "$BASE_URL/object?key=1*&bucket_name=posters" \
  "/object (S3 link by key)"

# /media/recommendations?type=&genre_ids=&limit=&offset=
test_get "$BASE_URL/media/recommendations?type=movie&genre_ids=1*&limit=3&offset=0" \
  "/media/recommendations (type + genre_ids + limit/offset)"

# /media/{media_id}
test_get "$BASE_URL/media/1*" \
  "/media/{media_id} (media details by id)"

# /media/{media_id}/episodes
test_get "$BASE_URL/media/1*/episodes" \
  "/media/{media_id}/episodes"

# /media/{media_id}/actor
test_get "$BASE_URL/media/1*/actor" \
  "/media/{media_id}/actor"

# /actor/{actor_id}
test_get "$BASE_URL/actor/1*" \
  "/actor/{actor_id}"

# /actor/{actor_id}/media
test_get "$BASE_URL/actor/1*/media" \
  "/actor/{actor_id}/media"

# /genre/{genre_id}
test_get "$BASE_URL/genre/1*" \
  "/genre/{genre_id}"

# /genre/{genre_id}/media?limit=&offset=
test_get "$BASE_URL/genre/1*/media?limit=10&offset=0" \
  "/genre/{genre_id}/media (with pagination)"

# /genre/all (no params, but quick sanity check)
test_get "$BASE_URL/genre/all" \
  "/genre/all (no parameters; quick sanity check)"

# /appeal/all?tag=&status=&limit=&offset=
test_get "$BASE_URL/appeal/all?tag=bug&status=open&limit=10*&offset=0" \
  "/appeal/all (filters + limit)"

# /search?query=&type=&limit=&offset=
test_get "$BASE_URL/search?query=1*&type=any&limit=10&offset=0" \
  "/search (query + pagination)"

# -------------------------------------------------------------------
# AUTHENTICATED ENDPOINTS (require ACCESS_TOKEN)
# -------------------------------------------------------------------

# /media/my?limit=&offset=&is_dislike=
test_get "$BASE_URL/media/my?limit=10*&offset=0&is_dislike=false" \
  "/media/my (limit/offset/is_dislike)" \
  "yes"

# /media/watch?media_id=
test_get "$BASE_URL/media/watch?media_id=1*" \
  "/media/watch (media_id query)" \
  "yes"

# /media/{media_id}/like (GET)
test_get "$BASE_URL/media/1*/like" \
  "/media/{media_id}/like (check like status)" \
  "yes"

# /appeal/my
test_get "$BASE_URL/appeal/my" \
  "/appeal/my (no params; sanity check)" \
  "yes"

# /appeal/{appeal_id}
test_get "$BASE_URL/appeal/1*" \
  "/appeal/{appeal_id}" \
  "yes"

# /appeal/{appeal_id}/message (GET)
test_get "$BASE_URL/appeal/1*/message" \
  "/appeal/{appeal_id}/message (GET messages)" \
  "yes"

echo
echo "[+] Done. Review sqlmap output above for any reported vulnerabilities."
