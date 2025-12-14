BEGIN;

-- Insert new movies with metadata, main video assets, posters
WITH
movie_media AS (
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
    ) VALUES
        ('movie', 'Finding Nemo', 'A clownfish embarks on a journey to find his missing son, accompanied by a forgetful fish.', DATE '2003-05-30', 8.2, 100, 0, 'United States', 'After his son is captured in the Great Barrier Reef and taken to Sydney, a timid clownfish sets out on a journey to bring him home.'),
        ('movie', 'Harry Potter and the Chamber of Secrets', 'Harry returns to Hogwarts and discovers a hidden chamber threatening the school.', DATE '2002-11-15', 7.4, 161, 13, 'United Kingdom', 'An ancient prophecy seems to be coming true when a mysterious presence begins stalking the corridors of Hogwarts.'),
        ('movie', 'Harry Potter and the Prisoner of Azkaban', 'Harry learns secrets about his past as an escaped prisoner threatens him.', DATE '2004-06-04', 7.9, 142, 13, 'United Kingdom', 'Harry Potter, Ron and Hermione return to Hogwarts where Sirius Black escapes Azkaban to find Harry.'),
        ('movie', 'Harry Potter and the Goblet of Fire', 'Harry is entered into a dangerous magical tournament.', DATE '2005-11-18', 7.7, 157, 13, 'United Kingdom', 'Harry finds himself competing in the Triwizard Tournament as dark forces return.'),
        ('movie', 'Harry Potter and the Order of the Phoenix', 'Harry forms a secret group as the Ministry denies Voldemort''s return.', DATE '2007-07-11', 7.5, 138, 13, 'United Kingdom', 'Harry returns to Hogwarts to train students against dark forces.'),
        ('movie', 'Inception', 'A thief enters dreams to plant an idea.', DATE '2010-07-16', 8.8, 148, 13, 'United States', 'A skilled thief is offered a chance at redemption if he can plant an idea in a target''s subconscious.'),
        ('movie', 'The Matrix', 'A hacker discovers reality is a simulation.', DATE '1999-03-31', 8.7, 136, 18, 'United States', 'A computer hacker learns his world is a simulated reality and joins rebels against machines.'),
        ('movie', 'The Matrix Reloaded', 'Neo continues the fight against the machines.', DATE '2003-05-15', 7.2, 138, 18, 'United States', 'Freedom fighters defend Zion as Neo discovers more about his powers.')
    RETURNING media_id, title
),
-- Main video assets
movie_assets AS (
    INSERT INTO asset (s3_key, mime_type, file_size_mb)
    SELECT
        CASE title
            WHEN 'Finding Nemo' THEN 'medias/Finding_Nemo.mp4'
            WHEN 'Harry Potter and the Chamber of Secrets' THEN 'medias/Harry_Potter_and_the_Chamber_of_Secrets.mp4'
            WHEN 'Harry Potter and the Prisoner of Azkaban' THEN 'medias/Harry_Potter_and_the_Prisoner_of_Azkaban.mp4'
            WHEN 'Harry Potter and the Goblet of Fire' THEN 'medias/Harry_Potter_and_the_Goblet_of_Fire.mp4'
            WHEN 'Harry Potter and the Order of the Phoenix' THEN 'medias/Harry_Potter_and_the_Order_of_the_Phoenix.mp4'
            WHEN 'Inception' THEN 'medias/InceptionMovie.mp4'
            WHEN 'The Matrix' THEN 'medias/MatrixMovie.mp4'
            WHEN 'The Matrix Reloaded' THEN 'medias/Matrix. Reloaded(2003).mp4'
        END,
        'video/mp4',
        0
    FROM movie_media
    RETURNING asset_id, s3_key
),
movie_videos AS (
    INSERT INTO asset_video (asset_id, quality, resolution_width, resolution_height)
    SELECT asset_id, '720p', 1280, 720
    FROM movie_assets
    RETURNING asset_video_id, asset_id
),
movie_main_video AS (
    INSERT INTO media_video (media_id, asset_video_id, video_type)
    SELECT mm.media_id, mv.asset_video_id, 'main_video'
    FROM movie_media mm
    JOIN movie_assets ma ON ma.s3_key = 
        CASE mm.title
            WHEN 'Finding Nemo' THEN 'medias/Finding_Nemo.mp4'
            WHEN 'Harry Potter and the Chamber of Secrets' THEN 'medias/Harry_Potter_and_the_Chamber_of_Secrets.mp4'
            WHEN 'Harry Potter and the Prisoner of Azkaban' THEN 'medias/Harry_Potter_and_the_Prisoner_of_Azkaban.mp4'
            WHEN 'Harry Potter and the Goblet of Fire' THEN 'medias/Harry_Potter_and_the_Goblet_of_Fire.mp4'
            WHEN 'Harry Potter and the Order of the Phoenix' THEN 'medias/Harry_Potter_and_the_Order_of_the_Phoenix.mp4'
            WHEN 'Inception' THEN 'medias/InceptionMovie.mp4'
            WHEN 'The Matrix' THEN 'medias/MatrixMovie.mp4'
            WHEN 'The Matrix Reloaded' THEN 'medias/Matrix. Reloaded(2003).mp4'
        END
    JOIN movie_videos mv ON mv.asset_id = ma.asset_id
),
-- Trailers
trailer_assets AS (
    INSERT INTO asset (s3_key, mime_type, file_size_mb)
    SELECT
        CASE title
            WHEN 'Finding Nemo' THEN 'trailers/Finding_Nemo.mp4'
            WHEN 'Harry Potter and the Chamber of Secrets' THEN 'trailers/Harry_Potter_and_the_Chamber_of_Secrets.mp4'
            WHEN 'Harry Potter and the Prisoner of Azkaban' THEN 'trailers/Harry_Potter_and_the_Prisoner_of_Azkaban.mp4'
            WHEN 'Harry Potter and the Goblet of Fire' THEN 'trailers/Harry_Potter_and_the_Goblet_of_Fire.mp4'
            WHEN 'Harry Potter and the Order of the Phoenix' THEN 'trailers/Harry_Potter_and_the_Order_of_the_Phoenix.mp4'
            WHEN 'Inception' THEN 'trailers/Inception.mp4'
            WHEN 'The Matrix' THEN 'trailers/The_Matrix.mp4'
            WHEN 'The Matrix Reloaded' THEN 'trailers/The_Matrix_Reloaded.mp4'
        END,
        'video/mp4',
        0
    FROM movie_media
    RETURNING asset_id, s3_key
),
trailer_videos AS (
    INSERT INTO asset_video (asset_id, quality, resolution_width, resolution_height)
    SELECT asset_id, '720p', 1280, 720
    FROM trailer_assets
    RETURNING asset_video_id, asset_id
),
trailer_media_video AS (
    INSERT INTO media_video (media_id, asset_video_id, video_type)
    SELECT mm.media_id, tv.asset_video_id, 'trailer'
    FROM movie_media mm
    JOIN trailer_assets ta ON ta.s3_key = 
        CASE mm.title
            WHEN 'Finding Nemo' THEN 'trailers/Finding_Nemo.mp4'
            WHEN 'Harry Potter and the Chamber of Secrets' THEN 'trailers/Harry_Potter_and_the_Chamber_of_Secrets.mp4'
            WHEN 'Harry Potter and the Prisoner of Azkaban' THEN 'trailers/Harry_Potter_and_the_Prisoner_of_Azkaban.mp4'
            WHEN 'Harry Potter and the Goblet of Fire' THEN 'trailers/Harry_Potter_and_the_Goblet_of_Fire.mp4'
            WHEN 'Harry Potter and the Order of the Phoenix' THEN 'trailers/Harry_Potter_and_the_Order_of_the_Phoenix.mp4'
            WHEN 'Inception' THEN 'trailers/Inception.mp4'
            WHEN 'The Matrix' THEN 'trailers/The_Matrix.mp4'
            WHEN 'The Matrix Reloaded' THEN 'trailers/The_Matrix_Reloaded.mp4'
        END
    JOIN trailer_videos tv ON tv.asset_id = ta.asset_id
),
-- Posters
poster_assets AS (
    INSERT INTO asset (s3_key, mime_type, file_size_mb)
    SELECT 'posters/' || REPLACE(LOWER(title), ' ', '_') || '.jpg', 'image/jpeg', 0
    FROM movie_media
    RETURNING asset_id, s3_key
),
poster_images AS (
    INSERT INTO asset_image (asset_id)
    SELECT asset_id FROM poster_assets
    RETURNING asset_image_id, asset_id
),
movie_posters AS (
    INSERT INTO media_image (media_id, asset_image_id, image_type)
    SELECT mm.media_id, pi.asset_image_id, 'poster'
    FROM movie_media mm
    JOIN poster_assets pa ON pa.s3_key = 'posters/' || REPLACE(LOWER(mm.title), ' ', '_') || '.jpg'
    JOIN poster_images pi ON pi.asset_id = pa.asset_id
),
-- Genres (using existing: Action=1, Sci-Fi=2, Thriller=3, Drama=4, Adventure=5, Fantasy=6)
movie_genres AS (
    INSERT INTO media_genre (media_id, genre_id)
    SELECT mm.media_id, UNNEST(ARRAY[
        CASE WHEN mm.title LIKE 'Finding Nemo%' THEN ARRAY[5,6]::int[]
             WHEN mm.title LIKE 'Harry Potter%' THEN ARRAY[5,6]::int[]
             WHEN mm.title = 'Inception' THEN ARRAY[1,2,3]::int[]
             WHEN mm.title LIKE 'The Matrix%' THEN ARRAY[1,2]::int[]
             ELSE ARRAY[]::int[]
        END
    ])
    FROM movie_media mm
)

SELECT 'Movies added successfully with trailers';

COMMIT;