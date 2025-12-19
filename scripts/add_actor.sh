#!/usr/bin/env bash
# /home/chinalap/Documents/VKEdu/Semester2/backend/scripts/add_actor.sh
# set -euo pipefail

if [[ "$#" -eq 0 ]]; then
    usage
    exit 1
fi

if [[ -z "${AWS_ACCESS_KEY_ID:-}" || -z "${AWS_SECRET_ACCESS_KEY:-}" || -z "${AWS_DEFAULT_REGION:-}" ]]; then
    echo "AWS credentials not set in environment variables."
    exit 1
fi

BUCKET="actors"

usage() {
    cat <<EOF
Usage: $0 --image-path PATH --actor-name NAME --birth-date YYYY-MM-DD --bio BIO --media-title TITLE --role-name ROLE
Options:
  -i, --image-path   Path to actor image file (required)
  -n, --actor-name   Actor name (required)
  -b, --birth-date   Birth date YYYY-MM-DD (required)
  -B, --bio          Bio text (required)
  -m, --media-title  Media title to link role to (required)
  -r, --role-name    Role name (required)
  -h, --help         Show this help
EOF
}

# Parse long options
PARSED=$(getopt -o i:n:b:B:m:r:h -l image-path:,actor-name:,birth-date:,bio:,media-title:,role-name:,help -- "$@") || { usage; exit 2; }
eval set -- "$PARSED"

IMAGE_PATH=""
ACTOR_NAME=""
BIRTH_DATE=""
BIO=""
MEDIA_TITLE=""
ROLE_NAME=""

while true; do
    case "$1" in
        -i|--image-path) IMAGE_PATH="$2"; shift 2 ;;
        -n|--actor-name) ACTOR_NAME="$2"; shift 2 ;;
        -b|--birth-date) BIRTH_DATE="$2"; shift 2 ;;
        -B|--bio) BIO="$2"; shift 2 ;;
        -m|--media-title) MEDIA_TITLE="$2"; shift 2 ;;
        -r|--role-name) ROLE_NAME="$2"; shift 2 ;;
        -h|--help) usage; exit 0 ;;
        --) shift; break ;;
        *) echo "Unexpected option: $1"; usage; exit 3 ;;
    esac
done

# Validate required args
if [[ -z "$IMAGE_PATH" || -z "$ACTOR_NAME" || -z "$BIRTH_DATE" || -z "$BIO" || -z "$MEDIA_TITLE" || -z "$ROLE_NAME" ]]; then
    echo "Missing required option(s)."
    usage
    exit 1
fi

if [[ ! -f "$IMAGE_PATH" ]]; then
    echo "Image file not found: $IMAGE_PATH"
    exit 1
fi

# Upload image to s3
if ! command -v aws >/dev/null 2>&1; then
    echo "AWS CLI not found. Please install it to proceed."
    exit 1
fi

KEY="$(basename "$IMAGE_PATH")"
echo "Uploading $IMAGE_PATH to s3://$BUCKET/$KEY ..."
aws s3api put-object --endpoint-url https://s3.cloud.ru --bucket "$BUCKET" --key "$KEY" --body "$IMAGE_PATH" --no-cli-pager
echo "Image uploaded to S3 with key: $KEY"

# Derive image info
IMAGE_NAME="$KEY"

if command -v file >/dev/null 2>&1; then
    MIME_TYPE=$(file --mime-type -b "$IMAGE_PATH")
else
    MIME_TYPE="application/octet-stream"
fi

if command -v identify >/dev/null 2>&1; then
    read -r WIDTH HEIGHT < <(identify -format "%w %h" "$IMAGE_PATH")
else
    WIDTH=0
    HEIGHT=0
fi

echo "Image info - MIME type: $MIME_TYPE, Width: $WIDTH, Height: $HEIGHT"

# Escape single quotes for safe SQL insertion
escape_sql() { printf "%s" "$1" | sed "s/'/''/g"; }
IMAGE_NAME=$(escape_sql "$IMAGE_NAME")
MIME_TYPE=$(escape_sql "$MIME_TYPE")
ACTOR_NAME=$(escape_sql "$ACTOR_NAME")
BIRTH_DATE=$(escape_sql "$BIRTH_DATE")
BIO=$(escape_sql "$BIO")
MEDIA_TITLE=$(escape_sql "$MEDIA_TITLE")
ROLE_NAME=$(escape_sql "$ROLE_NAME")

export IMAGE_NAME MIME_TYPE WIDTH HEIGHT ACTOR_NAME BIRTH_DATE BIO MEDIA_TITLE ROLE_NAME

# Append envsubst-processed SQL to target file
cat <<'SQL' | envsubst >> $HOME/Documents/VKEdu/Semester2/backend/testdata/postgres/10_addActors.sql
-- ${ACTOR_NAME}
BEGIN;
INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/${IMAGE_NAME}', '${MIME_TYPE}', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/${IMAGE_NAME}' LIMIT 1), ${WIDTH}, ${HEIGHT});

INSERT INTO actor (name, birth_date, bio) VALUES
    ('${ACTOR_NAME}', '${BIRTH_DATE}', '${BIO}');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = '${ACTOR_NAME}' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/${IMAGE_NAME}' LIMIT 1) LIMIT 1),
     'profile');
COMMIT;

BEGIN;
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = '${ACTOR_NAME}' LIMIT 1),
        (SELECT media_id FROM media WHERE title = '${MEDIA_TITLE}' LIMIT 1),
        '${ROLE_NAME}'
    );
COMMIT;

SQL

echo "SQL appended to ../testdata/postgres/10_addActors.sql"