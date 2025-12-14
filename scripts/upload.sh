#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  upload_media.sh --endpoint URL --trailers-bucket NAME --posters-bucket NAME [options] [FILES...]

Options:
  --endpoint URL            S3 endpoint, e.g. https://s3.cloud.ru
  --trailers-bucket NAME    Bucket for .mp4 files (e.g. trailers)
  --posters-bucket NAME     Bucket for image files (e.g. posters)
  --dir PATH                Scan this directory when no FILES are passed (default: .)
  --map FILE                Optional mapping file: "local_filename|object_key"
  --dry-run                 Print commands without running them
  --profile NAME            Optional AWS CLI profile
  --region REGION           Optional AWS region

Mapping file format examples:
  Finding_Nemo.mp4|Finding Nemo.mp4
  InceptionMovie.mp4|Inception.mp4
  MatrixMovie.mp4|The Matrix.mp4
  Matrix. Reloaded(2003).mp4|The Matrix Reloaded.mp4
EOF
}

ENDPOINT=""
TRAILERS_BUCKET=""
POSTERS_BUCKET=""
DIR="."
MAP_FILE=""
DRY_RUN=0
PROFILE=""
REGION=""

# --- args ---
while [[ $# -gt 0 ]]; do
  case "$1" in
    --endpoint) ENDPOINT="$2"; shift 2;;
    --trailers-bucket) TRAILERS_BUCKET="$2"; shift 2;;
    --posters-bucket) POSTERS_BUCKET="$2"; shift 2;;
    --dir) DIR="$2"; shift 2;;
    --map) MAP_FILE="$2"; shift 2;;
    --dry-run) DRY_RUN=1; shift;;
    --profile) PROFILE="$2"; shift 2;;
    --region) REGION="$2"; shift 2;;
    -h|--help) usage; exit 0;;
    --) shift; break;;
    -*) echo "Unknown option: $1" >&2; usage; exit 2;;
    *) break;;
  esac
done

if [[ -z "$ENDPOINT" || -z "$TRAILERS_BUCKET" || -z "$POSTERS_BUCKET" ]]; then
  echo "ERROR: --endpoint, --trailers-bucket, and --posters-bucket are required." >&2
  usage
  exit 2
fi

# --- load mapping (local -> key) ---
declare -A KEYMAP
if [[ -n "$MAP_FILE" ]]; then
  if [[ ! -f "$MAP_FILE" ]]; then
    echo "ERROR: map file not found: $MAP_FILE" >&2
    exit 2
  fi
  while IFS='|' read -r local key; do
    [[ -z "${local:-}" || -z "${key:-}" ]] && continue
    KEYMAP["$local"]="$key"
  done < "$MAP_FILE"
fi

aws_base=(aws)
[[ -n "$PROFILE" ]] && aws_base+=(--profile "$PROFILE")
[[ -n "$REGION"  ]] && aws_base+=(--region "$REGION")

guess_key() {
  local file="$1"
  local base ext
  base="$(basename "$file")"
  ext="${base##*.}"

  # mapping overrides everything
  if [[ -n "${KEYMAP[$base]+x}" ]]; then
    printf '%s\n' "${KEYMAP[$base]}"
    return 0
  fi

  # default: replace underscores with spaces, keep extension
  local stem="${base%.*}"
  stem="${stem//_/ }"
  printf '%s.%s\n' "$stem" "$ext"
}

bucket_for() {
  local file="$1"
  local lower="${file,,}"
  case "$lower" in
    *.mp4) echo "$TRAILERS_BUCKET" ;;
    *.jpg|*.jpeg|*.png|*.webp) echo "$POSTERS_BUCKET" ;;
    *) echo "" ;;
  esac
}

run_cmd() {
  if [[ "$DRY_RUN" -eq 1 ]]; then
    printf 'DRY-RUN: '
    printf '%q ' "$@"
    printf '\n'
  else
    "$@"
  fi
}

# --- build file list ---
files=()
if [[ $# -gt 0 ]]; then
  files=("$@")
else
  while IFS= read -r -d '' f; do
    files+=("$f")
  done < <(find "$DIR" -maxdepth 1 -type f \( -iname '*.mp4' -o -iname '*.jpg' -o -iname '*.jpeg' -o -iname '*.png' -o -iname '*.webp' \) -print0)
fi

if [[ ${#files[@]} -eq 0 ]]; then
  echo "No media files found." >&2
  exit 0
fi

# --- upload ---
for f in "${files[@]}"; do
  if [[ ! -f "$f" ]]; then
    echo "Skip (not a file): $f" >&2
    continue
  fi

  bucket="$(bucket_for "$f")"
  if [[ -z "$bucket" ]]; then
    echo "Skip (unsupported type): $f" >&2
    continue
  fi

  key="$(guess_key "$f")"

  run_cmd "${aws_base[@]}" s3api put-object \
    --endpoint-url "$ENDPOINT" \
    --bucket "$bucket" \
    --key "$key" \
    --body "$f"

  echo "Uploaded: $f  -> s3://$bucket/$key"
done
