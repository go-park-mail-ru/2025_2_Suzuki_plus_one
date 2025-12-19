#!/usr/bin/env bash
set -euo pipefail

# ----------------------------
# Vegeta perf runner for PopFilms
# ----------------------------
# Usage:
#   ./vegeta.sh --get-reco [--type movie|series] [--rate N] [--duration 30s]
#   ./vegeta.sh --post-signup --users N [--rate N] [--duration 30s]
#
# Env overrides:
#   ENDPOINT=http://localhost:8080
#   REPORT_DIR=reports
#   TARGETS_DIR=targets
#   RATE=50
#   DURATION=30s
#   TIMEOUT=10s
#   CONNECTIONS=0
#   KEEPALIVE=true
#
# Notes:
# - GET /media/recommendations requires query param "type"
# - POST /auth/signup needs unique email/username for each request:
#   we generate JSONL targets with base64 bodies and run vegeta with -format=json

# ---- Config ----
ENDPOINT="${ENDPOINT:-https://popfilms.ru/api/v1}"
REPORT_DIR="${REPORT_DIR:-reports}"
TARGETS_DIR="${TARGETS_DIR:-targets}"

RATE="${RATE:-50}"              # requests/sec
DURATION="${DURATION:-30s}"     # e.g. 30s, 2m
TIMEOUT="${TIMEOUT:-10s}"       # per request
CONNECTIONS="${CONNECTIONS:-0}" # 0 = default
KEEPALIVE="${KEEPALIVE:-true}"  # true/false

mkdir -p "$REPORT_DIR" "$TARGETS_DIR"

need() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "ERROR: '$1' not found. Install it and retry." >&2
    exit 1
  }
}

ts() { date +"%Y%m%d_%H%M%S"; }

usage() {
  cat <<EOF
Usage:
  $0 --get-reco [--type movie|series] [--rate N] [--duration 30s] [--timeout 10s]
  $0 --post-signup --users N [--rate N] [--duration 30s] [--timeout 10s]

Examples:
  ENDPOINT=http://localhost:8080 $0 --get-reco --type movie --rate 200 --duration 60s
  ENDPOINT=http://localhost:8080 $0 --post-signup --users 100000 --rate 100 --duration 1000s

Tips:
  - Keep USERS >= RATE * duration_seconds, otherwise you may run out of unique users and start getting 400 conflicts.
EOF
}

# ----------------------------
# Reporting
# ----------------------------
run_attack() {
  local name="$1"
  local targets_file="$2"
  local format="${3:-http}" # http|json

  local stamp; stamp="$(ts)"
  local out_bin="$REPORT_DIR/${name}_${stamp}.bin"
  local out_stderr="$REPORT_DIR/${name}_${stamp}_stderr.log"
  local out_txt="$REPORT_DIR/${name}_${stamp}_report.txt"
  local out_json="$REPORT_DIR/${name}_${stamp}_report.json"
  local out_hist="$REPORT_DIR/${name}_${stamp}_hist.txt"
  local out_plot="$REPORT_DIR/${name}_${stamp}_plot.html"

  echo "== Running: $name"
  echo "   endpoint: $ENDPOINT"
  echo "   targets:  $targets_file"
  echo "   format:   $format"
  echo "   rate:     $RATE"
  echo "   duration: $DURATION"
  echo "   timeout:  $TIMEOUT"
  echo

  # Build optional args
  local opt_conn=()
  local opt_keepalive=()
  if [[ "$CONNECTIONS" != "0" ]]; then opt_conn=(-connections="$CONNECTIONS"); fi
  if [[ "$KEEPALIVE" == "false" ]]; then opt_keepalive=(-keepalive=false); fi

  vegeta attack \
    -format="$format" \
    -rate="$RATE" \
    -duration="$DURATION" \
    -timeout="$TIMEOUT" \
    "${opt_conn[@]}" \
    "${opt_keepalive[@]}" \
    -targets="$targets_file" \
    > "$out_bin" 2> "$out_stderr" || true

  if [[ ! -s "$out_bin" ]]; then
    echo "ERROR: vegeta produced empty results: $out_bin" >&2
    echo "---- stderr ($out_stderr) ----" >&2
    sed -n '1,200p' "$out_stderr" >&2 || true
    exit 1
  fi

  vegeta report -type=text "$out_bin" > "$out_txt"
  vegeta report -type=json "$out_bin" > "$out_json"
  vegeta report -type=hist[0,50ms,100ms,200ms,500ms,1s,2s,5s] "$out_bin" > "$out_hist"
  vegeta plot "$out_bin" > "$out_plot"

  echo "== Done: $name"
  echo "   bin:    $out_bin"
  echo "   stderr: $out_stderr"
  echo "   report: $out_txt"
  echo "   json:   $out_json"
  echo "   hist:   $out_hist"
  echo "   plot:   $out_plot"
  echo
}

# ----------------------------
# Targets generators
# ----------------------------

# HTTP format targets for GET
make_targets_get_reco() {
  local reco_type="${1:-movie}" # movie|series
  local limit="${2:-3}"
  local offset="${3:-0}"

  local f="$TARGETS_DIR/get_media_recommendations_${reco_type}.http"
  cat > "$f" <<EOF
GET ${ENDPOINT}/media/recommendations?type=${reco_type}&limit=${limit}&offset=${offset}
Accept: application/json
User-Agent: vegeta

EOF
  echo "$f"
}

# JSON format targets for POST /auth/signup
generate_signup_targets_jsonl() {
  local n="$1"
  local out="$TARGETS_DIR/post_auth_signup_${n}.jsonl"

  # Check if we can reuse existing targets
  if [[ -f "$out" ]]; then
    local lines
    lines="$(wc -l < "$out" | tr -d ' ')"
    if [[ "$lines" == "$n" ]]; then
      echo "Reusing existing targets: $out ($lines lines)" >&2
      echo "$out"
      return 0
    fi
    echo "Targets file exists but has $lines lines (expected $n). Regenerating..." >&2
  else
    echo "Generating signup targets for $n users ..." >&2
  fi

  : >"$out"
  local base; base="$(date +%s)"

  for i in $(seq 1 "$n"); do
    local email="load_${base}_${i}@example.com"
    local username="load_${base}_${i}"
    local password="Password123!"
    local b64
    b64="$(printf '{"email":"%s","password":"%s","username":"%s"}' \
      "$email" "$password" "$username" | base64 | tr -d '\n')"

    printf '{"method":"POST","url":"%s/auth/signup","header":{"Content-Type":["application/json"],"Accept":["application/json"]},"body":"%s"}\n' \
      "$ENDPOINT" "$b64" >>"$out"
  done

  echo "$out"
}

# ----------------------------
# Args parsing
# ----------------------------
DO_GET_RECO=false
DO_POST_SIGNUP=false
USERS=1000
RECO_TYPE="movie"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --get-reco) DO_GET_RECO=true; shift ;;
    --post-signup) DO_POST_SIGNUP=true; shift ;;
    --users) USERS="$2"; shift 2 ;;
    --type) RECO_TYPE="$2"; shift 2 ;;
    --rate) RATE="$2"; shift 2 ;;
    --duration) DURATION="$2"; shift 2 ;;
    --timeout) TIMEOUT="$2"; shift 2 ;;
    --connections) CONNECTIONS="$2"; shift 2 ;;
    --keepalive) KEEPALIVE="$2"; shift 2 ;;
    -h|--help) usage; exit 0 ;;
    *)
      echo "Unknown arg: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if ! $DO_GET_RECO && ! $DO_POST_SIGNUP; then
  usage
  exit 1
fi

# ---- Dependencies ----
need vegeta
need base64
need date
need sed

# Basic validation
if [[ "$RECO_TYPE" != "movie" && "$RECO_TYPE" != "series" ]]; then
  echo "ERROR: --type must be 'movie' or 'series' (got '$RECO_TYPE')" >&2
  exit 1
fi

# ----------------------------
# Run one test per invocation (as you wanted)
# ----------------------------
if $DO_GET_RECO; then
  tfile="$(make_targets_get_reco "$RECO_TYPE")"
  run_attack "get_media_recommendations_${RECO_TYPE}" "$tfile" "http"
  exit 0
fi

if $DO_POST_SIGNUP; then
  tfile="$(generate_signup_targets_jsonl "$USERS")"
  run_attack "post_auth_signup_users_${USERS}" "$tfile" "json"
  exit 0
fi
