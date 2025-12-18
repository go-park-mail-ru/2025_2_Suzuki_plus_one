BEGIN;

WITH
-- 1) Series: The OA
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
        'The OA',
        'A mysterious woman resurfaces after being missing for seven years, now calling herself The OA.',
        DATE '2016-12-16',
        7.8,
        NULL,
        16,
        'United States',
        'After disappearing for seven years, Prairie Johnson returns with her sight restored and strange scars, revealing a story involving parallel dimensions, near-death experiences, and a secret mission.'
    )
    RETURNING media_id
),

-- 2) Poster
poster_asset AS (
    INSERT INTO asset (s3_key, mime_type, file_size_mb)
    VALUES ('posters/OA_poster.png', 'image/png', 0)
    RETURNING asset_id
),
poster_image AS (
    INSERT INTO asset_image (asset_id)
    SELECT asset_id FROM poster_asset
    RETURNING asset_image_id
),
series_image AS (
    INSERT INTO media_image (media_id, asset_image_id, image_type)
    SELECT sm.media_id, pi.asset_image_id, 'poster'
    FROM series_media sm
    CROSS JOIN poster_image pi
),

-- 3) Episode sources (S3 keys from your data)
episode_src AS (
    SELECT *
    FROM (VALUES
        -- Season 1
        ('medias/The.OA.S01E01.WEBRip.Rus.Eng.DV.LostFilm.mp4', 1, 1, 'The OA S01E01'),
        ('medias/The.OA.S01E02.WEBRip.Rus.Eng.DV.LostFilm.mp4', 1, 2, 'The OA S01E02'),
        ('medias/The.OA.S01E03.WEBRip.Rus.Eng.DV.LostFilm.mp4', 1, 3, 'The OA S01E03'),
        ('medias/The.OA.S01E04.WEBRip.Rus.Eng.DV.LostFilm.mp4', 1, 4, 'The OA S01E04'),
        ('medias/The.OA.S01E05.WEBRip.Rus.Eng.DV.LostFilm.mp4', 1, 5, 'The OA S01E05'),
        ('medias/The.OA.S01E06.WEBRip.Rus.Eng.DV.LostFilm.mp4', 1, 6, 'The OA S01E06'),
        ('medias/The.OA.S01E07.WEBRip.Rus.Eng.DV.LostFilm.mp4', 1, 7, 'The OA S01E07'),
        ('medias/The.OA.S01E08.WEBRip.Rus.Eng.DV.LostFilm.mp4', 1, 8, 'The OA S01E08'),

        -- Season 2
        ('medias/The.OA.S02E01.WEBRip.Rus.Eng.LostFilm.mp4', 2, 1, 'The OA S02E01'),
        ('medias/The.OA.S02E02.WEBRip.Rus.Eng.LostFilm.mp4', 2, 2, 'The OA S02E02'),
        ('medias/The.OA.S02E03.WEBRip.Rus.Eng.LostFilm.mp4', 2, 3, 'The OA S02E03'),
        ('medias/The.OA.S02E04.WEBRip.Rus.Eng.LostFilm.mp4', 2, 4, 'The OA S02E04'),
        ('medias/The.OA.S02E05.WEBRip.Rus.Eng.LostFilm.mp4', 2, 5, 'The OA S02E05'),
        ('medias/The.OA.S02E06.WEBRip.Rus.Eng.LostFilm.mp4', 2, 6, 'The OA S02E06'),
        ('medias/The.OA.S02E07.WEBRip.Rus.Eng.LostFilm.mp4', 2, 7, 'The OA S02E07'),
        ('medias/The.OA.S02E08.WEBRip.Rus.Eng.LostFilm.mp4', 2, 8, 'The OA S02E08')
    ) AS v(s3_key, season_number, episode_number, title)
),

-- 4) Episode media
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
        CONCAT('Episode ', season_number, 'x', LPAD(episode_number::text, 2, '0'), ' of The OA.'),
        CASE
            WHEN season_number = 1 THEN DATE '2016-12-16'
            WHEN season_number = 2 THEN DATE '2019-03-22'
        END,
        NULL,
        55,
        16,
        'United States',
        NULL
    FROM episode_src
    RETURNING media_id, title
),

-- 5) Episodes join
episodes AS (
    SELECT
        em.media_id AS episode_id,
        es.s3_key,
        es.season_number,
        es.episode_number
    FROM episode_media em
    JOIN episode_src es USING (title)
),

-- 5.5) Poster for episodes
episode_media_image AS (
    INSERT INTO media_image (media_id, asset_image_id, image_type)
    SELECT e.episode_id, pi.asset_image_id, 'poster'
    FROM episodes e
    CROSS JOIN poster_image pi
),

-- 6) Link episodes to series
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

-- 7) Episode assets
episode_assets AS (
    INSERT INTO asset (s3_key, mime_type, file_size_mb)
    SELECT
        e.s3_key,
        'video/mp4',
        0
    FROM episodes e
    RETURNING asset_id, s3_key
),

-- 8) Video assets
episode_videos AS (
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
    FROM episode_assets
    RETURNING asset_video_id, asset_id
),

-- 9) Media video
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

-- 10) Bridge
episode_mva AS (
    INSERT INTO media_video_asset (media_video_id, asset_video_id)
    SELECT media_video_id, asset_video_id
    FROM episode_media_video
),

-- 11) Trailer
trailer_asset AS (
    INSERT INTO asset (s3_key, mime_type, file_size_mb)
    VALUES ('trailers/OA_trailer.mp4', 'video/mp4', 0)
    RETURNING asset_id
),
trailer_video AS (
    INSERT INTO asset_video (
        asset_id,
        quality,
        resolution_width,
        resolution_height
    )
    SELECT asset_id, '720p', 1280, 720
    FROM trailer_asset
    RETURNING asset_video_id
),
trailer_media_video AS (
    INSERT INTO media_video (
        media_id,
        asset_video_id,
        video_type
    )
    SELECT sm.media_id, tv.asset_video_id, 'trailer'
    FROM series_media sm
    CROSS JOIN trailer_video tv
    RETURNING media_video_id, asset_video_id
)

-- 12) Trailer bridge
INSERT INTO media_video_asset (media_video_id, asset_video_id)
SELECT media_video_id, asset_video_id
FROM trailer_media_video;

COMMIT;
