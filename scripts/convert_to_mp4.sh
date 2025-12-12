#!/usr/bin/env bash
set -euo pipefail

########################################
# Config
########################################

# Folder with your source video files
DIR=$(realpath "../testdata/minio/medias/")

########################################
# Script
########################################

# If no matches, don't keep the literal pattern
shopt -s nullglob

cd "$DIR"

for in_file in *; do
  # Skip if not a regular file
  [[ -f "$in_file" ]] || continue

  # Skip already-mp4 files (case-insensitive)
  case "$in_file" in
    *.mp4|*.MP4)
      echo "Skipping MP4: $in_file"
      continue
      ;;
  esac

  filename="$in_file"
  name="${filename%.*}"
  out_file="${name}.mp4"

  echo "Converting: \"$in_file\" -> \"$out_file\""

  # Convert to streaming-friendly MP4 (H.264 + AAC + faststart)
  ffmpeg -y -i "$in_file" \
    -c:v libx264 -preset medium -crf 23 \
    -c:a aac -b:a 128k \
    -movflags +faststart \
    "$out_file"

  echo "Conversion done, removing original: \"$in_file\""
  rm -- "$in_file"
  echo
done

echo "All done in folder: $DIR"
