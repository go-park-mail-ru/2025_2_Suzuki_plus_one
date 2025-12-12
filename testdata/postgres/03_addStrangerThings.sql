BEGIN;

WITH
-- 1) Create the series row with richer metadata
series_media AS (
    INSERT INTO media (
        media_type,
        title,
        description,
        release_date,
        rating,
        duration_minutes,
        age_rating,
        country,
        plot_summary
    )
    VALUES (
        'series',
        'Stranger Things',
        'In 1980s Hawkins, a group of kids and adults face secret experiments, strange powers and a dangerous parallel dimension.',
        DATE '2016-07-15',
        8.7,           -- overall series rating
        NULL,          -- episodes vary in length, leave series-level duration null
        16,            -- Dutch-like 16+ age advice
        'United States',
        'Set in 1980s Hawkins, Indiana, the series follows kids, families and local authorities as a telekinetic girl escapes a government lab and a portal to a hostile dimension unleashes monsters and other threats on the town.'
    )
    RETURNING media_id
),

-- 2) Poster for the series
poster_asset AS (
    INSERT INTO asset (s3_key, mime_type, file_size_mb)
    VALUES ('posters/Stranger_Things.jpg', 'image/jpeg', 0)
    RETURNING asset_id
),
poster_image AS (
    INSERT INTO asset_image (asset_id)
    SELECT asset_id
    FROM poster_asset
    RETURNING asset_image_id
),
series_image AS (
    INSERT INTO media_image (media_id, asset_image_id, image_type)
    SELECT sm.media_id, pi.asset_image_id, 'poster'
    FROM series_media sm
    CROSS JOIN poster_image pi
),

-- 3) All episode sources: file path + season/episode + title
episode_src AS (
    SELECT *
    FROM (VALUES
        ('medias/Stranger.Things.S01E01.Chapter.One.The.Vanishing.Of.Will.Byers.720p.mp4', 1, 1, 'Stranger Things S01E01 - Chapter One: The Vanishing of Will Byers'),
        ('medias/Stranger.Things.S01E02.Chapter.Two.The.Weirdo.On.Maple.Street.720p.mp4',   1, 2, 'Stranger Things S01E02 - Chapter Two: The Weirdo on Maple Street'),
        ('medias/Stranger.Things.S01E03.Chapter.Three.Holly.Jolly.720p.mp4',               1, 3, 'Stranger Things S01E03 - Chapter Three: Holly Jolly'),
        ('medias/Stranger.Things.S01E04.Chapter.Four.The.Body.720p.mp4',                   1, 4, 'Stranger Things S01E04 - Chapter Four: The Body'),
        ('medias/Stranger.Things.S01E05.Chapter.Five.The.Flea.And.The.Acrobat.720p.mp4',   1, 5, 'Stranger Things S01E05 - Chapter Five: The Flea and the Acrobat'),
        ('medias/Stranger.Things.S01E06.Chapter.Six.The.Monster.720p.mp4',                 1, 6, 'Stranger Things S01E06 - Chapter Six: The Monster'),
        ('medias/Stranger.Things.S01E07.Chapter.Seven.The.Bathtub.720p.mp4',               1, 7, 'Stranger Things S01E07 - Chapter Seven: The Bathtub'),
        ('medias/Stranger.Things.S01E08.Chapter.Eight.The.Upside.Down.720p.mp4',           1, 8, 'Stranger Things S01E08 - Chapter Eight: The Upside Down'),
        ('medias/Stranger.Things.S02E01.720p.mp4',                                         2, 1, 'Stranger Things S02E01 - Chapter One: MADMAX'),
        ('medias/Stranger.Things.S02E02.720p.mp4',                                         2, 2, 'Stranger Things S02E02 - Chapter Two: Trick or Treat, Freak'),
        ('medias/Stranger.Things.S02E03.720p.mp4',                                         2, 3, 'Stranger Things S02E03 - Chapter Three: The Pollywog'),
        ('medias/Stranger.Things.S02E04.720p.mp4',                                         2, 4, 'Stranger Things S02E04 - Chapter Four: Will the Wise'),
        ('medias/Stranger.Things.S02E05.720p.mkv',                                         2, 5, 'Stranger Things S02E05 - Chapter Five: Dig Dug'),
        ('medias/Stranger.Things.S02E06.720p.mkv',                                         2, 6, 'Stranger Things S02E06 - Chapter Six: The Spy'),
        ('medias/Stranger.Things.S02E07.720p.mkv',                                         2, 7, 'Stranger Things S02E07 - Chapter Seven: The Lost Sister'),
        ('medias/Stranger.Things.S02E08.720p.mkv',                                         2, 8, 'Stranger Things S02E08 - Chapter Eight: The Mind Flayer'),
        ('medias/Stranger.Things.S02E09.720p.mkv',                                         2, 9, 'Stranger Things S02E09 - Chapter Nine: The Gate')
    ) AS v(s3_key, season_number, episode_number, title)
),

-- 4) Insert media rows for all episodes with basic metadata
episode_media AS (
    INSERT INTO media (
        media_type,
        title,
        description,
        release_date,
        rating,
        duration_minutes,
        age_rating,
        country,
        plot_summary
    )
    SELECT
        'episode',
        title,
        -- short generic description; you can refine later per episode
        CONCAT('Episode ', season_number, 'x', LPAD(episode_number::text, 2, '0'),
               ' of Stranger Things in Hawkins, Indiana.'),
        CASE
            WHEN season_number = 1 THEN DATE '2016-07-15'
            WHEN season_number = 2 THEN DATE '2017-10-27'
        END,
        NULL, -- per-episode rating can be filled later if desired
        CASE
            WHEN season_number = 1 THEN 50
            WHEN season_number = 2 THEN 50
        END,
        16,
        'United States',
        NULL
    FROM episode_src
    RETURNING media_id, title
),

-- 5) Join inserted episodes back to their src info (file path + S/E)
episodes AS (
    SELECT
        em.media_id       AS episode_id,
        es.s3_key,
        es.season_number,
        es.episode_number
    FROM episode_media em
    JOIN episode_src es USING (title)
),

-- 5.5) Link series poster image to episode media as well
episode_media_image AS (
    INSERT INTO media_image (media_id, asset_image_id, image_type)
SELECT
    e.episode_id,
    pi.asset_image_id,
    'poster'
FROM episodes e
CROSS JOIN poster_image pi),

-- 6) Link each episode to the series in media_episode
media_episode_insert AS (
    INSERT INTO media_episode (
        episode_id,
        series_id,
        season_number,
        episode_number
    )
    SELECT
        e.episode_id,
        sm.media_id,
        e.season_number,
        e.episode_number
    FROM episodes e
    CROSS JOIN series_media sm
),

-- 7) Create asset rows for each episode file
episode_assets AS (
    INSERT INTO asset (s3_key, mime_type, file_size_mb)
    SELECT DISTINCT
        e.s3_key,
        CASE
            WHEN e.s3_key LIKE '%.mkv'  THEN 'video/x-matroska'
            WHEN e.s3_key LIKE '%.mp4'  THEN 'video/mp4'
            WHEN e.s3_key LIKE '%.webm' THEN 'video/webm'
        END,
        0
    FROM episodes e
    RETURNING asset_id, s3_key
),

-- 8) Wrap those assets as asset_video (assume 720p)
episode_videos AS (
    INSERT INTO asset_video (
        asset_id,
        quality,
        resolution_width,
        resolution_height
    )
    SELECT
        ea.asset_id,
        '720p',
        1280,
        720
    FROM episode_assets ea
    RETURNING asset_video_id, asset_id
),

-- 9) Link each episode media row to its video asset as main_video
episode_media_video AS (
    INSERT INTO media_video (
        media_id,
        asset_video_id,
        video_type
    )
    SELECT
        e.episode_id,
        av.asset_video_id,
        'main_video'
    FROM episodes e
    JOIN episode_assets ea USING (s3_key)
    JOIN episode_videos av ON av.asset_id = ea.asset_id
    RETURNING media_video_id, asset_video_id
),

-- 10) Also fill the media_video_asset bridge for episodes
episode_mva AS (
    INSERT INTO media_video_asset (media_video_id, asset_video_id)
    SELECT
        media_video_id,
        asset_video_id
    FROM episode_media_video
),

-- 11) Series trailer asset & video
trailer_asset AS (
    INSERT INTO asset (s3_key, mime_type, file_size_mb)
    VALUES ('trailers/Stranger_Things_Season_1_Trailer_1.mp4', 'video/mp4', 0)
    RETURNING asset_id
),
trailer_video AS (
    INSERT INTO asset_video (
        asset_id,
        quality,
        resolution_width,
        resolution_height
    )
    SELECT
        asset_id,
        '720p',
        1280,
        720
    FROM trailer_asset
    RETURNING asset_video_id
),
trailer_media_video AS (
    INSERT INTO media_video (
        media_id,
        asset_video_id,
        video_type
    )
    SELECT
        sm.media_id,
        tv.asset_video_id,
        'trailer'
    FROM series_media sm
    CROSS JOIN trailer_video tv
    RETURNING media_video_id, asset_video_id
)

-- 12) Bridge row for the trailer (series-level trailer)
INSERT INTO media_video_asset (media_video_id, asset_video_id)
SELECT
    media_video_id,
    asset_video_id
FROM trailer_media_video;

COMMIT;
