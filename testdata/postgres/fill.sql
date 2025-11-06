-- Begin transaction
BEGIN;

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
-- Movies
('movie', 'Inception', 'A thief who steals corporate secrets through dream-sharing technology is given the inverse task of planting an idea into the mind of a C.E.O.', '2010-07-16', 8.8, 148, 12, 'USA', 'Dom Cobb is a skilled thief, the absolute best in the dangerous art of extraction, stealing valuable secrets from deep within the subconscious during the dream state.'),
('movie', 'The Matrix', 'A computer hacker learns from mysterious rebels about the true nature of his reality and his role in the war against its controllers.', '1999-03-31', 8.7, 136, 15, 'USA', 'Thomas Anderson, a computer programmer, is led to fight an underground war against powerful computers who have constructed his entire reality with a system called the Matrix.'),
-- Series
('series', 'Breaking Bad', 'A high school chemistry teacher diagnosed with inoperable lung cancer turns to manufacturing and selling methamphetamine in order to secure his family''s future.', '2008-01-20', 9.5, 50, 18, 'USA', 'Walter White, a chemistry teacher, discovers he has cancer and decides to get into the meth-making business to repay his medical debts.'),
-- Episodes for Breaking Bad
('episode', 'Pilot', 'Diagnosed with terminal lung cancer, chemistry teacher Walter White teams up with former student Jesse Pinkman to cook and sell crystal meth.', '2008-01-20', 8.9, 58, 18, 'USA', 'Walter White learns he has lung cancer and begins his journey into the criminal underworld.'),
('episode', 'Cat''s in the Bag...', 'After their first drug deal goes terribly wrong, Walt and Jesse struggle to cover their tracks.', '2008-01-27', 8.7, 48, 18, 'USA', 'Walt and Jesse deal with the aftermath of their first encounter with the criminal world.');

-- Insert episode relationships
INSERT INTO media_episode (episode_id, series_id, season_number, episode_number) VALUES
(4, 3, 1, 1),  -- Pilot episode
(5, 3, 1, 2);  -- Cat's in the Bag episode

-- Link media to genres
INSERT INTO media_genre (media_id, genre_id) VALUES
-- Inception genres
(1, 1), (1, 2), (1, 3),  -- Action, Sci-Fi, Thriller
-- Matrix genres
(2, 1), (2, 2), (2, 3),  -- Action, Sci-Fi, Thriller
-- Breaking Bad genres
(3, 4), (3, 3);  -- Drama, Thriller

-- Insert assets (files stored in S3)
INSERT INTO asset (s3_key, mime_type, size_mb) VALUES
-- Actor images
('/actors/leo.png', 'image/png', 2.5),
('/actors/morgan.png', 'image/png', 2.3),
-- Avatars
('/avatars/Chris.png', 'image/png', 1.8),
-- Posters
('/posters/InceptionPoster.png', 'image/png', 3.2),
('/posters/MatrixPoster.png', 'image/png', 3.1),
-- Media files
('/medias/InceptionMovie.webm', 'video/webm', 1500.0),
('/medias/MatrixMovie.webm', 'video/webm', 1400.0),
-- Trailers
('/trailers/InceptionTrailer.webm', 'video/webm', 120.5),
('/trailers/MatrixTrailer.webm', 'video/webm', 110.3);

-- Insert asset images
INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
(2, 300, 450),  -- leo.png
(3, 300, 450),  -- morgan.png
(4, 200, 200),  -- Chris.png (avatar)
(5, 400, 600),  -- InceptionPoster.png
(6, 400, 600);  -- MatrixPoster.png

-- Insert asset videos
INSERT INTO asset_video (asset_id, quality, resolution_width, resolution_height) VALUES
(7, '1080p', 1920, 1080),  -- InceptionMovie.webm
(8, '1080p', 1920, 1080),  -- MatrixMovie.webm
(9, '720p', 1280, 720),    -- InceptionTrailer.webm
(10, '720p', 1280, 720);    -- MatrixTrailer.webm

-- Link media to images (posters)
INSERT INTO media_image (media_id, asset_image_id, image_type) VALUES
(1, 5, 'poster'),  -- Inception poster
(2, 6, 'poster'),  -- Matrix poster
(3, 5, 'poster');  -- Breaking Bad uses Inception poster for example

-- Link media to videos
INSERT INTO media_video (media_id, asset_video_id, video_type) VALUES
(1, 7, 'main_video'),  -- Inception main video
(1, 9, 'trailer'),     -- Inception trailer
(2, 8, 'main_video'),  -- Matrix main video
(2, 10, 'trailer');     -- Matrix trailer

-- Insert actors
INSERT INTO actor (name, birth_date, bio) VALUES
('Leonardo DiCaprio', '1974-11-11', 'Academy Award-winning actor known for his roles in Titanic, Inception, and The Revenant.'),
('Morgan Freeman', '1937-06-01', 'Renowned actor known for his distinctive voice and roles in The Shawshank Redemption and Driving Miss Daisy.'),
('Keanu Reeves', '1964-09-02', 'Canadian actor known for The Matrix series and John Wick franchise.');

-- Link actor images
INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
(1, 2, 'profile'),  -- Leo profile image
(2, 3, 'profile');  -- Morgan profile image

-- Insert actor roles
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
(1, 1, 'Dom Cobb'),      -- Leo in Inception
(3, 2, 'Neo');           -- Keanu in Matrix

-- Insert users
INSERT INTO "user" (username, asset_image_id, password_hash, date_of_birth, phone_number, email) VALUES
('Chris', 3, '$2b$10$examplehashedpassword123456789012', '1990-05-15', '+1234567890', 'chris@example.com'),
('Alex', NULL, '$2b$10$examplehashedpassword123456789013', '1985-08-20', '+0987654321', 'alex@example.com');

-- Insert user sessions
INSERT INTO user_session (user_id, session_token, expires_at) VALUES
(1, 'chris_session_token_123', CURRENT_TIMESTAMP + INTERVAL '30 days'),
(2, 'alex_session_token_456', CURRENT_TIMESTAMP + INTERVAL '30 days');

-- Insert playlists
INSERT INTO playlist (user_id, name, description, visibility) VALUES
(1, 'My Favorite Movies', 'A collection of my all-time favorite films', 'public'),
(1, 'Sci-Fi Collection', 'The best science fiction movies and shows', 'unlisted'),
(2, 'Private Watchlist', 'Movies I plan to watch', 'private');

-- Link playlist media
INSERT INTO playlist_media (playlist_id, media_id) VALUES
(1, 1), (1, 2),  -- Chris favorites: Inception and Matrix
(2, 1), (2, 2), (2, 3),  -- Sci-Fi collection
(3, 1);  -- Alex private: Inception

-- Insert user playlist roles
INSERT INTO user_playlist (user_id, playlist_id, role) VALUES
(1, 1, 'owner'),
(1, 2, 'owner'),
(2, 3, 'owner'),
(2, 1, 'viewer');  -- Alex can view Chris's favorite movies

-- Insert watch history
INSERT INTO user_watch_history (user_id, media_id, progress_seconds) VALUES
(1, 1, 8880),  -- Chris watched 2h28m of Inception (148min movie)
(1, 2, 8160),  -- Chris watched 2h16m of Matrix
(2, 1, 3600);  -- Alex watched 1h of Inception

-- Insert likes
INSERT INTO user_like_media (user_id, media_id) VALUES
(1, 1), (1, 2), (1, 3),  -- Chris likes all media
(2, 1);  -- Alex likes Inception

INSERT INTO user_like_actor (user_id, actor_id) VALUES
(1, 1), (1, 3),  -- Chris likes Leo and Keanu
(2, 2);  -- Alex likes Morgan

INSERT INTO user_like_playlist (user_id, playlist_id) VALUES
(2, 1);  -- Alex likes Chris's favorite movies playlist

-- Insert comments
INSERT INTO user_comment_media (user_id, media_id, content) VALUES
(1, 1, 'Mind-blowing concept and incredible visuals! One of Nolan''s best works.'),
(1, 2, 'Revolutionary film that changed action movies forever.'),
(2, 1, 'The dream within a dream concept still confuses me but I love it!');

INSERT INTO user_comment_actor (user_id, actor_id, content) VALUES
(1, 1, 'Leo always delivers outstanding performances!'),
(2, 2, 'Morgan Freeman has the most iconic voice in Hollywood.');

-- Insert ratings
INSERT INTO user_rating_media (user_id, media_id, rating) VALUES
(1, 1, 5),  -- Chris rates Inception 5 stars
(1, 2, 5),  -- Chris rates Matrix 5 stars
(1, 3, 4),  -- Chris rates Breaking Bad 4 stars
(2, 1, 5),  -- Alex rates Inception 5 stars
(2, 2, 4);  -- Alex rates Matrix 4 stars

-- Commit transaction
COMMIT;

-- Verify data counts
DO $$
DECLARE
    media_count INTEGER;
    user_count INTEGER;
    actor_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO media_count FROM media;
    SELECT COUNT(*) INTO user_count FROM "user";
    SELECT COUNT(*) INTO actor_count FROM actor;

    RAISE NOTICE 'Database populated successfully:';
    RAISE NOTICE '- % media entries', media_count;
    RAISE NOTICE '- % users', user_count;
    RAISE NOTICE '- % actors', actor_count;
    RAISE NOTICE '- Genres, playlists, comments, and ratings added';
END $$;