BEGIN;

WITH
-- 1) Series: Scream Queens
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
        'Scream Queens',
        'A satirical horror-comedy series blending slasher tropes with dark humor and campus drama.',
        DATE '2015-09-22',
        7.2,
        NULL,
        18,
        'United States',
        'Set on a college campus ruled by an elite sorority, Scream Queens mixes classic slasher horror with sharp comedy. When a masked serial killer known as the Red Devil begins murdering students, secrets surface, power structures collapse, and the line between satire and horror becomes increasingly blurred.'
    )
    RETURNING media_id
),

-- 2) Poster
poster_asset AS (
    INSERT INTO asset (s3_key, mime_type, file_size_mb)
    VALUES ('posters/Scream_Queens_poster.png', 'image/png', 0)
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

-- 3) Episode sources + titles + detailed plots
episode_src AS (
    SELECT *
    FROM (VALUES
        ('medias/Scream.Queens.2015.S01E01E02.Pilot-Hell.Week.1080p.WEB-DL.DD5.1.H.264.Rus.Eng-BkD.mp4', 1, 1, 'Pilot / Hell Week',
         'The series opens with a deadly sorority secret from the past. Chanel Oberlin rules Kappa Kappa Tau with cruelty and precision until a masked killer begins targeting the campus during Hell Week.'),
        ('medias/Scream.Queens.2015.S01E03.Chainsaw.1080p.WEB-DL.DD5.1.H.264.Rus.Eng-BkD.mp4', 1, 3, 'Chainsaw',
         'A campus fundraiser turns into chaos when the Red Devil strikes again, forcing the sorority to confront the growing body count.'),
        ('medias/Scream.Queens.2015.S01E04.Haunted.House.1080p.WEB-DL.DD5.1.H.264.Rus.Eng-BkD.mp4', 1, 4, 'Haunted House',
         'A deadly haunted house attraction becomes the perfect hunting ground as paranoia spreads among students and faculty.'),
        ('medias/Scream.Queens.2015.S01E05.Pumpkin.Patch.1080p.WEB-DL.DD5.1.H.264.Rus.Eng-BkD.mp4', 1, 5, 'Pumpkin Patch',
         'The killer exploits Halloween traditions while Grace and Denise dig deeper into the sorority’s dark history.'),
        ('medias/Scream.Queens.2015.S01E06.Seven.Minutes.in.Hell.1080p.WEB-DL.DD5.1.H.264.Rus.Eng-BkD.mp4', 1, 6, 'Seven Minutes in Hell',
         'Secrets and alliances shift during a charity event as romantic tensions rise and another murder rocks the campus.'),
        ('medias/Scream.Queens.2015.S01E07.Beware.of.Young.Girls.1080p.WEB-DL.DD5.1.H.264.Rus.Eng-BkD.mp4', 1, 7, 'Beware of Young Girls',
         'Flashbacks reveal the origins of the Red Devil myth while present-day events spiral out of control.'),
        ('medias/Scream.Queens.2015.S01E08.Mommie.Dearest.1080p.WEB-DL.DD5.1.H.264.Rus.Eng-BkD.mp4', 1, 8, 'Mommie Dearest',
         'Family visits expose buried trauma and explain the twisted motivations behind several suspects.'),
        ('medias/Scream.Queens.2015.S01E09.Ghost.Stories.1080p.WEB-DL.DD5.1.H.264.Rus.Eng-BkD.mp4', 1, 9, 'Ghost Stories',
         'Supernatural rumors and campus legends fuel fear as the killer’s identity seems closer than ever.'),
        ('medias/Scream.Queens.2015.S01E10.Thanksgiving.1080p.WEB-DL.DD5.1.H.264.Rus.Eng-BkD.mp4', 1, 10, 'Thanksgiving',
         'A holiday dinner becomes a bloodbath as long-suspected characters meet shocking fates.'),
        ('medias/Scream.Queens.2015.S01E11.Black.Friday.1080p.WEB-DL.DD5.1.H.264.Rus.Eng-BkD.mp4', 1, 11, 'Black Friday',
         'A department store massacre pushes the investigation into public chaos and media frenzy.'),
        ('medias/Scream.Queens.2015.S01E12.Dorkus.1080p.WEB-DL.DD5.1.H.264.Rus.Eng-BkD.mp4', 1, 12, 'Dorkus',
         'The mystery unravels as past and present collide, revealing the truth behind the killings.'),
        ('medias/Scream.Queens.2015.S01E13.The.Final.Girl(S).1080p.WEB-DL.DD5.1.H.264.Rus.Eng-BkD.mp4', 1, 13, 'The Final Girl(s)',
         'Season finale: Survivors face the ultimate showdown as the Red Devil’s identity is finally exposed.')
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
        CONCAT('Episode ', season_number, 'x', LPAD(episode_number::text, 2, '0'), ' of Scream Queens.'),
        DATE '2015-09-22',
        NULL,
        45,
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
    VALUES ('trailers/Scream_Queens_trailer.mp4', 'video/mp4', 0)
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
