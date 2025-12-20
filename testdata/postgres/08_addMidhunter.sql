BEGIN;

WITH
-- 1) Series: Mindhunter
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
        'Mindhunter',
        'A psychological crime drama exploring the early days of criminal profiling at the FBI.',
        DATE '2017-10-13',
        8.6,
        NULL,
        18,
        'United States',
        'Set in the late 1970s, Mindhunter follows FBI agents Holden Ford and Bill Tench as they pioneer the science of criminal profiling. By interviewing imprisoned serial killers, they attempt to understand how murderers think, operate, and evolve — confronting disturbing truths that challenge their personal lives and the limits of traditional law enforcement.'
    )
    RETURNING media_id
),

-- 2) Poster
poster_asset AS (
    INSERT INTO asset (s3_key, mime_type, file_size_mb)
    VALUES ('posters/Mindhunter_poster.png', 'image/png', 0)
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

-- 3) Episode sources + rich titles
episode_src AS (
    SELECT *
    FROM (VALUES
        ('medias/Mindhunter.S01E01.1080p.rus.sub.FOCS-PB.mp4', 1, 1, 'Mindhunter S01E01', 'FBI negotiator Holden Ford begins questioning traditional crime-solving methods after a hostage situation exposes the limits of conventional profiling.'),
        ('medias/Mindhunter.S01E02.1080p.rus.sub.FOCS-PB.mp4', 1, 2, 'Mindhunter S01E02', 'Holden and Bill Tench conduct their first prison interview with a convicted killer, discovering how unsettling it is to step inside a murderer’s mind.'),
        ('medias/Mindhunter.S01E03.1080p.rus.sub.FOCS-PB.mp4', 1, 3, 'Mindhunter S01E03', 'The Behavioral Science Unit faces internal resistance as Holden pushes controversial theories that blur moral and professional boundaries.'),
        ('medias/Mindhunter.S01E04.1080p.rus.sub.FOCS-PB.mp4', 1, 4, 'Mindhunter S01E04', 'A new case forces the team to apply their emerging profiling techniques in the real world, with disturbing consequences.'),
        ('medias/Mindhunter.S01E05.1080p.rus.sub.FOCS-PB.mp4', 1, 5, 'Mindhunter S01E05', 'Holden’s growing obsession with serial killers strains his relationship and raises concerns about his mental stability.'),
        ('medias/Mindhunter.S01E06.1080p.rus.sub.FOCS-PB.mp4', 1, 6, 'Mindhunter S01E06', 'As the interviews intensify, Bill Tench struggles to balance family life with the emotional toll of studying violent criminals.'),
        ('medias/Mindhunter.S01E07.1080p.rus.sub.FOCS-PB.mp4', 1, 7, 'Mindhunter S01E07', 'The team refines their profiling language while encountering chilling patterns that connect different murder cases.'),
        ('medias/Mindhunter.S01E08.1080p.rus.sub.FOCS-PB.mp4', 1, 8, 'Mindhunter S01E08', 'Holden becomes dangerously close to his subjects as the line between empathy and identification begins to fade.'),
        ('medias/Mindhunter.S01E09.1080p.rus.sub.FOCS-PB.mp4', 1, 9, 'Mindhunter S01E09', 'A high-stakes investigation tests whether behavioral profiling can actually save lives—or only explain death.'),
        ('medias/Mindhunter.S01E10.1080p.rus.sub.FOCS-PB.mp4', 1, 10,'Mindhunter S01E10','Season finale: Holden faces the personal and professional cost of diving too deep into the darkest corners of the human psyche.')
    ) AS v(s3_key, season_number, episode_number, title, episode_plot)
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
        CONCAT('Episode ', season_number, 'x', LPAD(episode_number::text, 2, '0'), ' of Mindhunter.'),
        DATE '2017-10-13',
        NULL,
        55,
        18,
        'United States',
        episode_plot
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

-- 6) Series link
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
        s3_key,
        'video/mp4',
        0
    FROM episodes
    RETURNING asset_id, s3_key
),

-- 8) Asset video (1080p)
episode_videos AS (
    INSERT INTO asset_video (
        asset_id,
        quality,
        resolution_width,
        resolution_height
    )
    SELECT
        asset_id,
        '1080p',
        1920,
        1080
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
    VALUES ('trailers/Mindhunter_trailer.mp4', 'video/mp4', 0)
    RETURNING asset_id
),
trailer_video AS (
    INSERT INTO asset_video (
        asset_id,
        quality,
        resolution_width,
        resolution_height
    )
    SELECT asset_id, '1080p', 1920, 1080
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
