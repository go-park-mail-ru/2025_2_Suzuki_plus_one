-- Begin transaction
BEGIN;

-- First, let's check what assets already exist
SELECT asset_id, s3_key FROM asset ORDER BY asset_id;

-- Insert genres
INSERT INTO genre (name, description) VALUES
('Action', 'High-energy films with physical stunts and chases'),
('Sci-Fi', 'Futuristic technology, space exploration, and scientific themes'),
('Thriller', 'Suspenseful stories that keep viewers on edge'),
('Drama', 'Serious, character-driven stories focusing on emotional themes'),
('Adventure', 'Exciting journeys and exploration'),
('Fantasy', 'Magical elements, mythical creatures, and imaginary worlds');

-- Insert media (movies and series)
INSERT INTO media (media_type, title, description, release_date, rating, duration_minutes, age_rating, country, plot_summary) VALUES
('movie', 'Toy Story', 'Led by Woody, Andy''s toys live happily in his room until Andy''s birthday brings Buzz Lightyear onto the scene. Afraid of losing his place in Andy''s heart, Woody plots against Buzz. But when circumstances separate Buzz and Woody from their owner, the duo eventually learns to put aside their differences.', '1995-11-22', 8.0, 81, 13, 'United States of America', 'The adventure takes off when toys come to life!');

-- Link media to genres
INSERT INTO media_genre (media_id, genre_id) VALUES
(1, 5), (1, 6), (1, 4),
(2, 5), (2, 6), (2, 1);

-- Insert ONLY the video assets first so we can track their IDs
INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
('/medias/1_ToyStoryMovie.mp4', 'video/mp4', 1500.0),
('/trailers/1_ToyStoryTrailer.mp4', 'video/mp4', 120.5);

-- Get the asset IDs for the videos we just inserted
DO $$
DECLARE
video1_id BIGINT;
trailer1_id BIGINT;
BEGIN
SELECT asset_id INTO video1_id FROM asset WHERE s3_key = '/medias/1_ToyStory.mp4';
SELECT asset_id INTO trailer1_id FROM asset WHERE s3_key = '/trailers/1_ToyStoryTrailer.mp4';

RAISE NOTICE 'Video assets inserted with IDs: %, %', video1_id, trailer1_id;

    -- Now insert asset_video records with the correct asset IDs
INSERT INTO asset_video (asset_id, quality, resolution_width, resolution_height) VALUES
(video1_id, '720p', 1280, 720),
(trailer1_id, '720p', 1280, 720);
END $$;

-- Now insert the rest of the assets (actor images and posters)
INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
-- Actor images
('/actors/Tom_Hanks.png', 'image/png', 0.07),
('/actors/Tim_Allen.png', 'image/png', 0.07),
('/actors/Don_Rickles.png', 'image/png', 0.07),
('/actors/Jim_Varney.png', 'image/png', 0.07),
('/actors/Wallace_Shawn.png', 'image/png', 0.07),
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
(1, 3, 'trailer');     -- Toy Story trailer

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

-- Insert additional users
INSERT INTO "user" (username, asset_image_id, password_hash, date_of_birth, phone_number, email) VALUES
('Chris', 3, '$2b$10$examplehashedpassword123456789012', '1990-05-15', '+1234567890', 'chris@example.com'),
('Alex', NULL, '$2b$10$examplehashedpassword123456789013', '1985-08-20', '+0987654321', 'alex@example.com');

-- Insert user sessions
INSERT INTO user_session (user_id, session_token, expires_at) VALUES
(2, 'chris_session_token_123', CURRENT_TIMESTAMP + INTERVAL '30 days'),
(3, 'alex_session_token_456', CURRENT_TIMESTAMP + INTERVAL '30 days');

-- Insert playlists
INSERT INTO playlist (user_id, name, description, visibility) VALUES
(2, 'My Favorite Movies', 'A collection of my all-time favorite films', 'public'),
(2, 'Sci-Fi Collection', 'The best science fiction movies and shows', 'unlisted'),
(3, 'Private Watchlist', 'Movies I plan to watch', 'private');

-- Link playlist media
INSERT INTO playlist_media (playlist_id, media_id) VALUES
(1, 1), (1, 2),
(2, 1), (2, 2),
(3, 1);

-- Insert user playlist roles
INSERT INTO user_playlist (user_id, playlist_id, role) VALUES
(2, 1, 'owner'),
(2, 2, 'owner'),
(3, 3, 'owner'),
(3, 1, 'viewer');

-- Insert watch history
INSERT INTO user_watch_history (user_id, media_id, progress_seconds) VALUES
(2, 1, 8880),
(2, 2, 8160),
(3, 1, 3600);

-- Insert likes
INSERT INTO user_like_media (user_id, media_id) VALUES
(2, 1), (2, 2),
(3, 1);

INSERT INTO user_like_actor (user_id, actor_id) VALUES
(2, 1), (2, 2);

INSERT INTO user_like_playlist (user_id, playlist_id) VALUES
    (3, 1);

-- Insert comments
INSERT INTO user_comment_media (user_id, media_id, content) VALUES
(2, 1, 'Mind-blowing concept and incredible visuals! One of the best animated films.'),
(2, 2, 'Revolutionary film that combined live action and animation beautifully.'),
(3, 1, 'The chemistry between Woody and Buzz is amazing!');

INSERT INTO user_comment_actor (user_id, actor_id, content) VALUES
(2, 1, 'Tom Hanks always delivers outstanding performances!'),
(3, 2, 'Tim Allen has great comedic timing.');

-- Insert ratings
INSERT INTO user_rating_media (user_id, media_id, rating) VALUES
(2, 1, 5),
(2, 2, 5),
(3, 1, 5),
(3, 2, 4);

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
RAISE NOTICE '- Genres, playlists, comments, and ratings added';
END $$;
