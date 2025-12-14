BEGIN;

WITH
-- 0) One mapping table: movie title -> s3_key
movie_asset_map(title, s3_key) AS (
  VALUES
    ('Toy Story',                 'medias/1_ToyStoryMovie.mp4'),
    ('Jumanji',                   'medias/Jumanji.mp4'),
    ('Grumpier Old Men',          'medias/3_GrumpierOldManMovie.mp4'),
    ('Waiting to Exhale',         'medias/Waiting_to_exhale.mp4'),
    ('Father of the Bride Part II','medias/Father_of_the_Bride 2.mp4'),
    ('Heat',                      'medias/6_HeatMovie.mp4'),
    ('Sabrina',                   'medias/Sabrina_1995.mp4'),
    ('Tom and Huck',              'medias/Tom_and_Huck.mp4'),
    ('Sudden Death',              'medias/Sudden_Death.mp4'),
    ('GoldenEye',                 'medias/10_GoldenEyeMovie.mp4')
),

-- 1) Find the media rows we want to attach videos to
movie_media AS (
  SELECT m.media_id, m.title
  FROM media m
  JOIN movie_asset_map map ON map.title = m.title
  WHERE m.media_type = 'movie'
),

-- 2) Ensure ASSET rows exist for these s3_keys
ins_asset AS (
  INSERT INTO asset (s3_key, mime_type, file_size_mb)
  SELECT map.s3_key, 'video/mp4', 0::real
  FROM movie_asset_map map
  WHERE NOT EXISTS (
    SELECT 1 FROM asset a WHERE a.s3_key = map.s3_key
  )
  RETURNING asset_id, s3_key
),

-- 3) Collect all assets (inserted or already existing)
all_assets AS (
  SELECT asset_id, s3_key
  FROM ins_asset
  UNION
  SELECT a.asset_id, a.s3_key
  FROM asset a
  JOIN movie_asset_map map ON map.s3_key = a.s3_key
),

-- 4) Ensure ASSET_VIDEO rows exist (one 720p record per asset)
ins_asset_video AS (
  INSERT INTO asset_video (asset_id, quality, resolution_width, resolution_height)
  SELECT aa.asset_id, '720p', 1280, 720
  FROM all_assets aa
  WHERE NOT EXISTS (
    SELECT 1
    FROM asset_video av
    WHERE av.asset_id = aa.asset_id
      AND av.quality = '720p'
  )
  RETURNING asset_video_id, asset_id
),

-- 6) Link MEDIA -> ASSET_VIDEO via MEDIA_VIDEO (main_video)
ins_media_video AS (
  INSERT INTO media_video (media_id, asset_video_id, video_type)
  SELECT
    mm.media_id,
    aav.asset_video_id,
    'main_video'
  FROM movie_media mm
  JOIN movie_asset_map map ON map.title = mm.title
  JOIN all_assets aa ON aa.s3_key = map.s3_key
  JOIN ins_asset_video aav ON aav.asset_id = aa.asset_id
  WHERE NOT EXISTS (
    SELECT 1
    FROM media_video mv
    WHERE mv.media_id = mm.media_id
      AND mv.video_type = 'main_video'
  )
  RETURNING media_video_id, media_id, asset_video_id
)

-- Small summary so you see what happened
SELECT
  (SELECT COUNT(*) FROM movie_media)      AS movies_found,
  (SELECT COUNT(*) FROM ins_asset)        AS assets_inserted,
  (SELECT COUNT(*) FROM ins_asset_video)  AS asset_videos_inserted,
  (SELECT COUNT(*) FROM ins_media_video)  AS media_videos_inserted;

COMMIT;


