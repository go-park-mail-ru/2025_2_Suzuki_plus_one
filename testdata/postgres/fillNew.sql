-- Begin transaction
BEGIN;

-- First, let's check what assets already exist
SELECT asset_id, s3_key FROM asset ORDER BY asset_id;

-- Insert genres
INSERT INTO genre (name, description) VALUES
                                          ('Animation', 'Animated films and series'),
                                          ('Adventure', 'Exciting journeys and exploration'),
                                          ('Family', 'Family-friendly content'),
                                          ('Comedy', 'Funny and entertaining films');

-- Insert media (only Toy Story)
INSERT INTO media (media_type, title, description, release_date, rating, duration_minutes, age_rating, country, plot_summary) VALUES
    ('movie', 'Toy Story', 'Led by Woody, Andy''s toys live happily in his room until Andy''s birthday brings Buzz Lightyear onto the scene. Afraid of losing his place in Andy''s heart, Woody plots against Buzz. But when circumstances separate Buzz and Woody from their owner, the duo eventually learns to put aside their differences.', '1995-11-22', 8.0, 81, 0, 'United States of America', 'The adventure takes off when toys come to life!');

-- Link media to genres
INSERT INTO media_genre (media_id, genre_id) VALUES
                                                 (1, 1), (1, 2), (1, 3), (1, 4);

-- Insert ONLY the video assets first so we can track their IDs
INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
                                                        ('/medias/InceptionMovie.webm', 'video/webm', 1500.0),
                                                        ('/trailers/InceptionTrailer.webm', 'video/webm', 120.5);

-- Get the asset IDs for the videos we just inserted
DO $$
DECLARE
video1_id BIGINT;
    trailer1_id BIGINT;
BEGIN
SELECT asset_id INTO video1_id FROM asset WHERE s3_key = '/medias/InceptionMovie.webm';
SELECT asset_id INTO trailer1_id FROM asset WHERE s3_key = '/trailers/InceptionTrailer.webm';

RAISE NOTICE 'Video assets inserted with IDs: %, %', video1_id, trailer1_id;

    -- Now insert asset_video records with the correct asset IDs
INSERT INTO asset_video (asset_id, quality, resolution_width, resolution_height) VALUES
                                                                                     (video1_id, '1080p', 1920, 1080),
                                                                                     (trailer1_id, '720p', 1280, 720);
END $$;

-- Now insert the rest of the assets (actor images and posters)
INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
-- Actor images
('/actors/Tom_Hanks.png', 'image/png', 0.07),
('/actors/Tim_Allen.png', 'image/png', 0.07),
-- Posters
('/posters/1_Toy_Story.png', 'image/png', 0.1);

-- Insert asset_images for actors and posters
INSERT INTO asset_image (asset_id, resolution_width, resolution_height)
SELECT asset_id, 500, 750
FROM asset
WHERE s3_key LIKE '/actors/%' OR s3_key LIKE '/posters/%';

-- Link media to videos using the asset_video records
INSERT INTO media_video (media_id, asset_video_id, video_type) VALUES
                                                                   (1, 1, 'main_video'),  -- Toy Story main video
                                                                   (1, 2, 'trailer');     -- Toy Story trailer

-- Link media to posters
INSERT INTO media_image (media_id, asset_image_id, image_type)
SELECT 1, asset_image_id, 'poster'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/posters/1_Toy_Story.png';

-- Insert actors
INSERT INTO actor (name, birth_date, bio) VALUES
                                              ('Tom Hanks', '1956-07-09', 'Thomas Jeffrey Hanks (born July 9, 1956) is an American actor and filmmaker.'),
                                              ('Tim Allen', '1953-06-13', 'Tim Allen (born Timothy Allen Dick; June 13, 1953) is an American comedian, actor, voice-over artist, and entertainer.');

-- Link actor images
INSERT INTO actor_image (actor_id, asset_image_id, image_type)
SELECT 1, asset_image_id, 'profile'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/actors/Tom_Hanks.png';

INSERT INTO actor_image (actor_id, asset_image_id, image_type)
SELECT 2, asset_image_id, 'profile'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/actors/Tim_Allen.png';

-- Insert actor roles
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
                                                           (1, 1, 'Woody (voice)'),
                                                           (2, 1, 'Buzz Lightyear (voice)');

-- Insert one test user
INSERT INTO "user" (username, password_hash, date_of_birth, phone_number, email) VALUES
    ('testuser', '$2b$10$examplehashedpassword123456789012', '1990-01-01', '+1234567890', 'testuser@example.com');

-- Insert user session
INSERT INTO user_session (user_id, session_token, expires_at) VALUES
    (1, 'testuser_session_token_123', CURRENT_TIMESTAMP + INTERVAL '30 days');

-- Insert playlist
INSERT INTO playlist (user_id, name, description, visibility) VALUES
    (1, 'My Favorite Movies', 'A collection of my all-time favorite films', 'public');

-- Link playlist media
INSERT INTO playlist_media (playlist_id, media_id) VALUES
    (1, 1);  -- Toy Story in favorite movies

-- Insert user playlist role
INSERT INTO user_playlist (user_id, playlist_id, role) VALUES
    (1, 1, 'owner');

-- Insert watch history
INSERT INTO user_watch_history (user_id, media_id, progress_seconds) VALUES
    (1, 1, 2400);  -- testuser watched 40 min of Toy Story

-- Insert likes
INSERT INTO user_like_media (user_id, media_id) VALUES
    (1, 1);  -- testuser likes Toy Story

INSERT INTO user_like_actor (user_id, actor_id) VALUES
                                                    (1, 1),  -- testuser likes Tom Hanks
                                                    (1, 2);  -- testuser likes Tim Allen

-- Insert rating
INSERT INTO user_rating_media (user_id, media_id, rating) VALUES
    (1, 1, 5);  -- testuser rates Toy Story 5 stars

-- Commit transaction
COMMIT;

-- Verify data counts
DO $$
DECLARE
media_count INTEGER;
    user_count INTEGER;
    actor_count INTEGER;
    asset_count INTEGER;
    asset_video_count INTEGER;
BEGIN
SELECT COUNT(*) INTO media_count FROM media;
SELECT COUNT(*) INTO user_count FROM "user";
SELECT COUNT(*) INTO actor_count FROM actor;
SELECT COUNT(*) INTO asset_count FROM asset;
SELECT COUNT(*) INTO asset_video_count FROM asset_video;

RAISE NOTICE 'Database populated successfully:';
    RAISE NOTICE '- % media entries', media_count;
    RAISE NOTICE '- % users', user_count;
    RAISE NOTICE '- % actors', actor_count;
    RAISE NOTICE '- % assets', asset_count;
    RAISE NOTICE '- % asset videos', asset_video_count;
    RAISE NOTICE '- Genres, playlist, likes, and rating added';
END $$;