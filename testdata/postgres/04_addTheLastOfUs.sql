BEGIN;

WITH
-- Series metadata
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
    ) VALUES (
        'series',
        'The Last of Us',
        'A hardened survivor escorts a young girl across a post-apocalyptic United States.',
        DATE '2023-01-15',
        9.2,
        NULL,
        18,
        'United States',
        'Twenty years after a fungal pandemic, Joel smuggles Ellie, who may be key to a cure, across a dangerous world.'
    )
    RETURNING media_id
),
-- Poster for series
poster_asset AS (
    INSERT INTO asset (s3_key, mime_type, file_size_mb)
    VALUES ('posters/The_Last_of_Us.jpg', 'image/jpeg', 0)
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
    FROM series_media sm, poster_image pi
),
-- Episode sources
episode_src AS (
    SELECT * FROM (VALUES
        ('medias/The.Last.of.Us.S01E01.When.Youre.Lost.in.the.Darkness.mp4', 1, 1, 'The Last of Us S01E01 - When You''re Lost in the Darkness'),
        ('medias/The.Last.of.Us.S01E02.Infected.mp4', 1, 2, 'The Last of Us S01E02 - Infected'),
        ('medias/The.Last.of.Us.S01E03.Long.Long.Time.mp4', 1, 3, 'The Last of Us S01E03 - Long, Long Time'),
        ('medias/The.Last.of.Us.S01E04.Please.Hold.to.My.Hand.mp4', 1, 4, 'The Last of Us S01E04 - Please Hold to My Hand'),
        ('medias/The.Last.of.Us.S01E05.Endure.and.Survive.mp4', 1, 5, 'The Last of Us S01E05 - Endure and Survive'),
        ('medias/The.Last.of.Us.S01E06.Kin.mp4', 1, 6, 'The Last of Us S01E06 - Kin'),
        ('medias/The.Last.of.Us.S01E07.Left.Behind.mp4', 1, 7, 'The Last of Us S01E07 - Left Behind'),
        ('medias/The.Last.of.Us.S01E08.When.We.Are.in.Need.mp4', 1, 8, 'The Last of Us S01E08 - When We Are in Need'),
        ('medias/The.Last.of.Us.S01E09.Look.for.the.Light.mp4', 1, 9, 'The Last of Us S01E09 - Look for the Light')
    ) AS v(s3_key, season_number, episode_number, title)
),
-- Episode media
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
        es.title,
        CONCAT('Episode ', es.season_number, 'x', LPAD(es.episode_number::text, 2, '0'), ' of The Last of Us.'),
        DATE '2023-01-15' + INTERVAL '7 days' * (es.episode_number - 1),
        NULL,
        60,
        18,
        'United States',
        NULL
    FROM episode_src es
    RETURNING media_id, title
),
-- Link episodes to source
episodes AS (
    SELECT em.media_id AS episode_id, es.s3_key, es.season_number, es.episode_number
    FROM episode_media em
    JOIN episode_src es ON es.title = em.title
),
-- Link poster to episodes
episode_poster AS (
    INSERT INTO media_image (media_id, asset_image_id, image_type)
    SELECT e.episode_id, pi.asset_image_id, 'poster'
    FROM episodes e, poster_image pi
),
-- Link to series
media_episode_insert AS (
    INSERT INTO media_episode (episode_id, series_id, season_number, episode_number)
    SELECT e.episode_id, sm.media_id, e.season_number, e.episode_number
    FROM episodes e, series_media sm
),
-- Episode video assets
episode_assets AS (
    INSERT INTO asset (s3_key, mime_type, file_size_mb)
    SELECT s3_key, 'video/mp4', 0
    FROM episodes
    RETURNING asset_id, s3_key
),
episode_videos AS (
    INSERT INTO asset_video (asset_id, quality, resolution_width, resolution_height)
    SELECT asset_id, '720p', 1280, 720
    FROM episode_assets
    RETURNING asset_video_id, asset_id
),
episode_main_video AS (
    INSERT INTO media_video (media_id, asset_video_id, video_type)
    SELECT e.episode_id, av.asset_video_id, 'main_video'
    FROM episodes e
    JOIN episode_assets ea ON ea.s3_key = e.s3_key
    JOIN episode_videos av ON av.asset_id = ea.asset_id
)

SELECT 'The Last of Us series added successfully';

COMMIT;