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
                                                                                                                                  ('movie', 'Toy Story', 'Led by Woody, Andy''s toys live happily in his room until Andy''s birthday brings Buzz Lightyear onto the scene. Afraid of losing his place in Andy''s heart, Woody plots against Buzz. But when circumstances separate Buzz and Woody from their owner, the duo eventually learns to put aside their differences.', '1995-11-22', 8.0, 81, 13, 'United States of America', 'The adventure takes off when toys come to life!'),
                                                                                                                                  ('movie', 'Jumanji', 'When siblings Judy and Peter discover an enchanted board game that opens the door to a magical world, they unwittingly invite Alan -- an adult who''s been trapped inside the game for 26 years -- into their living room. Alan''s only hope for freedom is to finish the game, which proves risky as all three find themselves running from giant rhinoceroses, evil monkeys and other terrifying creatures.', '1995-12-15', 7.241, 104, 13, 'United States of America', 'It''s a jungle in here.'),
                                                                                                                                  ('movie', 'Grumpier Old Men', 'John and Max are still living next door to each other and still fighting, until a beautiful woman moves across the street.', '1995-12-22', 6.6, 101, 13, 'United States of America', 'The feud continues... with romance.'),
                                                                                                                                  ('movie', 'Waiting to Exhale', 'The story of four African-American women who lean on each other while navigating relationships, careers, and dreams.', '1995-12-22', 6.6, 127, 13, 'United States of America', 'Four friends. Four lives. One unforgettable journey.'),
                                                                                                                                  ('movie', 'Father of the Bride Part II', 'George Banks must deal not only with his daughter''s pregnancy but also with his wife''s.', '1995-02-10', 6.1, 106, 13, 'United States of America', 'Double the trouble, double the fun.'),
                                                                                                                                  ('movie', 'Heat', 'A group of professional bank robbers start to feel the heat from police when they unknowingly leave a clue at their latest heist.', '1995-12-15', 8.3, 170, 18, 'United States of America', 'A Los Angeles crime saga.'),
                                                                                                                                  ('movie', 'Sabrina', 'A young woman falls for a wealthy businessman while working as a chauffeur''s daughter at his estate.', '1995-12-15', 6.3, 127, 13, 'United States of America', 'Love is in the most unexpected places.'),
                                                                                                                                  ('movie', 'Tom and Huck', 'The classic adventure of Tom Sawyer and Huckleberry Finn on the Mississippi River.', '1995-12-22', 5.7, 97, 13, 'United States of America', 'The adventure of a lifetime.'),
                                                                                                                                  ('movie', 'Sudden Death', 'A former firefighter must save his daughter and the Vice President from terrorists during a Stanley Cup finals game.', '1995-10-27', 5.6, 111, 18, 'United States of America', 'Terror has no time limit.'),
                                                                                                                                  ('movie', 'GoldenEye', 'James Bond sets out to stop a Russian crime syndicate from using a satellite weapon against London.', '1995-11-17', 7.2, 130, 13, 'United Kingdom', 'Bond is back with a vengeance.');

-- Link media to genres
INSERT INTO media_genre (media_id, genre_id) VALUES
                                                 (1, 5), (1, 6), (1, 4),
                                                 (2, 5), (2, 6), (2, 1),
                                                 (3, 5), (3, 6), (3, 1),
                                                 (4, 5), (4, 6), (4, 1),
                                                 (5, 5), (5, 6), (5, 1),
                                                 (6, 5), (6, 6), (6, 1),
                                                 (7, 5), (7, 6), (7, 1),
                                                 (8, 5), (8, 6), (8, 1),
                                                 (9, 5), (9, 6), (9, 1),
                                                 (10, 5), (10, 6), (10, 1);

-- Insert ONLY the video assets first so we can track their IDs
INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
                                                        ('/trailers/1_ToyStoryTrailer.webm', 'video/webm', 1500.0),
                                                        ('/trailers/2_JumanjiTrailer.webm', 'video/webm', 1400.0),
                                                        ('/trailers/3_GrumpierOldMenTrailer.mp4', 'video/mp4', 1600.0),
                                                        ('/trailers/4_WaitingToExhaleTrailer.mp4', 'video/mp4', 1700.0),
                                                        ('/trailers/5_FatherOfTheBrideIITrailer.mp4', 'video/mp4', 1550.0),
                                                        ('/trailers/6_HeatTrailer.mp4', 'video/mp4', 2200.0),
                                                        ('/trailers/7_SabrinaTrailer.mp4', 'video/mp4', 1650.0),
                                                        ('/trailers/8_TomAndHuckTrailer.mp4', 'video/mp4', 1450.0),
                                                        ('/trailers/9_SuddenDeathTrailer.mp4', 'video/mp4', 1800.0),
                                                        ('/trailers/10_GoldenEyeTrailer.mp4', 'video/mp4', 2000.0);


-- Get the asset IDs for the videos we just inserted
DO $$
DECLARE
trailer1_id BIGINT;
trailer2_id BIGINT;
trailer3_id BIGINT;
trailer4_id BIGINT;
trailer5_id BIGINT;
trailer6_id BIGINT;
trailer7_id BIGINT;
trailer8_id BIGINT;
trailer9_id BIGINT;
trailer10_id BIGINT;
BEGIN
SELECT asset_id INTO trailer1_id FROM asset WHERE s3_key = '/trailers/1_ToyStoryTrailer.webm';
SELECT asset_id INTO trailer2_id FROM asset WHERE s3_key = '/trailers/2_JumanjiTrailer.webm';
SELECT asset_id INTO trailer3_id FROM asset WHERE s3_key = '/trailers/3_GrumpierOldMenTrailer.mp4';
SELECT asset_id INTO trailer4_id FROM asset WHERE s3_key = '/trailers/4_WaitingToExhaleTrailer.mp4';
SELECT asset_id INTO trailer5_id FROM asset WHERE s3_key = '/trailers/5_FatherOfTheBrideIITrailer.mp4';
SELECT asset_id INTO trailer6_id FROM asset WHERE s3_key = '/trailers/6_HeatTrailer.mp4';
SELECT asset_id INTO trailer7_id FROM asset WHERE s3_key = '/trailers/7_SabrinaTrailer.mp4';
SELECT asset_id INTO trailer8_id FROM asset WHERE s3_key = '/trailers/8_TomAndHuckTrailer.mp4';
SELECT asset_id INTO trailer9_id FROM asset WHERE s3_key = '/trailers/9_SuddenDeathTrailer.mp4';
SELECT asset_id INTO trailer10_id FROM asset WHERE s3_key = '/trailers/10_GoldenEyeTrailer.mp4';

RAISE NOTICE 'Video assets inserted with IDs: %, %, %, %, %, %, %, %, %, %', trailer1_id, trailer2_id, trailer3_id, trailer4_id, trailer5_id, trailer6_id, trailer7_id, trailer8_id, trailer9_id, trailer10_id;

    -- Now insert asset_video records with the correct asset IDs
INSERT INTO asset_video (asset_id, quality, resolution_width, resolution_height) VALUES
(trailer1_id, '720p', 1280, 720),
(trailer2_id, '720p', 1280, 720),
(trailer3_id, '720p', 1280, 720),
(trailer4_id, '720p', 1280, 720),
(trailer5_id, '720p', 1280, 720),
(trailer6_id, '720p', 1280, 720),
(trailer7_id, '720p', 1280, 720),
(trailer8_id, '720p', 1280, 720),
(trailer9_id, '720p', 1280, 720),
(trailer10_id, '720p', 1280, 720);
END $$;

-- Now insert the rest of the assets (actor images and posters)
INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
-- Actor images
('/actors/Tom_Hanks.png', 'image/png', 0.07),
('/actors/Tim_Allen.png', 'image/png', 0.07),
('/actors/Don_Rickles.png', 'image/png', 0.07),
('/actors/Jim_Varney.png', 'image/png', 0.07),
('/actors/Wallace_Shawn.png', 'image/png', 0.07),
('/actors/Robert_De_Niro.png', 'image/png', 0.07),
('/actors/Al_Pacino.png', 'image/png', 0.07),
('/actors/Harrison_Ford.png', 'image/png', 0.07),
-- Posters
('/posters/1_Toy_Story.png', 'image/png', 0.1),
('/posters/2_Jumanji.png', 'image/png', 0.1),
('/posters/3_Grumpier_Old_Men.png', 'image/png', 0.1),
('/posters/4_Waiting_to_Exhale.png', 'image/png', 0.1),
('/posters/5_Father_of_the_Bride_Part_II.png', 'image/png', 0.1),
('/posters/6_Heat.png', 'image/png', 0.1),
('/posters/7_Sabrina.png', 'image/png', 0.1),
('/posters/8_Tom_and_Huck.png', 'image/png', 0.1),
('/posters/9_Sudden_Death.png', 'image/png', 0.1),
('/posters/10_GoldenEye.png', 'image/png', 0.1);

-- Insert asset_images for actors and posters
INSERT INTO asset_image (asset_id, resolution_width, resolution_height)
SELECT asset_id, 500, 750
FROM asset
WHERE s3_key LIKE '/actors/%' OR s3_key LIKE '/posters/%';

-- Link media to videos using the asset_video records
INSERT INTO media_video (media_id, asset_video_id, video_type) VALUES
                                                                   (1, 1, 'trailer'),  -- Toy Story trailer
                                                                   (2, 2, 'trailer'),  -- Jumanji trailer
                                                                   (3, 3, 'trailer'),  -- Grumpier Old Men
                                                                   (4, 4, 'trailer'),  -- Waiting to Exhale
                                                                   (5, 5, 'trailer'),  -- Father of the Bride Part II
                                                                   (6, 6, 'trailer'),  -- Heat
                                                                   (7, 7, 'trailer'),  -- Sabrina
                                                                   (8, 8, 'trailer'),  -- Tom and Huck
                                                                   (9, 9, 'trailer'),  -- Sudden Death
                                                                   (10, 10, 'trailer'); -- GoldenEye

-- Link media to posters
INSERT INTO media_image (media_id, asset_image_id, image_type)
SELECT 1, asset_image_id, 'poster'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/posters/1_Toy_Story.png';

INSERT INTO media_image (media_id, asset_image_id, image_type)
SELECT 2, asset_image_id, 'poster'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/posters/2_Jumanji.png';

INSERT INTO media_image (media_id, asset_image_id, image_type)
SELECT 3, asset_image_id, 'poster'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/posters/3_Grumpier_Old_Men.png';

INSERT INTO media_image (media_id, asset_image_id, image_type)
SELECT 4, asset_image_id, 'poster'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/posters/4_Waiting_to_Exhale.png';

INSERT INTO media_image (media_id, asset_image_id, image_type)
SELECT 5, asset_image_id, 'poster'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/posters/5_Father_of_the_Bride_Part_II.png';

INSERT INTO media_image (media_id, asset_image_id, image_type)
SELECT 6, asset_image_id, 'poster'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/posters/6_Heat.png';

INSERT INTO media_image (media_id, asset_image_id, image_type)
SELECT 7, asset_image_id, 'poster'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/posters/7_Sabrina.png';

INSERT INTO media_image (media_id, asset_image_id, image_type)
SELECT 8, asset_image_id, 'poster'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/posters/8_Tom_and_Huck.png';

INSERT INTO media_image (media_id, asset_image_id, image_type)
SELECT 9, asset_image_id, 'poster'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/posters/9_Sudden_Death.png';

INSERT INTO media_image (media_id, asset_image_id, image_type)
SELECT 10, asset_image_id, 'poster'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/posters/10_GoldenEye.png';

-- Insert actors
INSERT INTO actor (name, birth_date, bio) VALUES
                                              ('Tom Hanks', '1956-07-09', 'Thomas Jeffrey Hanks (born July 9, 1956) is an American actor and filmmaker.'),
                                              ('Tim Allen', '1953-06-13', 'Tim Allen (born Timothy Allen Dick; June 13, 1953) is an American comedian, actor, voice-over artist, and entertainer.'),
                                              ('Robert De Niro', '1943-08-17', 'Robert Anthony De Niro is an American actor, producer, and director.'),
                                              ('Al Pacino', '1940-04-25', 'Alfredo James Pacino is an American actor and filmmaker.'),
                                              ('Harrison Ford', '1942-07-13', 'Harrison Ford is an American actor.');

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

INSERT INTO actor_image (actor_id, asset_image_id, image_type)
SELECT 3, asset_image_id, 'profile'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/actors/Robert_De_Niro.png';

INSERT INTO actor_image (actor_id, asset_image_id, image_type)
SELECT 4, asset_image_id, 'profile'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/actors/Al_Pacino.png';

INSERT INTO actor_image (actor_id, asset_image_id, image_type)
SELECT 5, asset_image_id, 'profile'
FROM asset_image ai
         JOIN asset a ON ai.asset_id = a.asset_id
WHERE a.s3_key = '/actors/Harrison_Ford.png';

-- Insert actor roles
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
                                                           (1, 1, 'Woody (voice)'),
                                                           (2, 1, 'Buzz Lightyear (voice)'),
                                                           (3, 6, 'Neil McCauley'),
                                                           (4, 6, 'Lt. Vincent Hanna'),
                                                           (5, 10, 'James Bond');

-- Insert additional users
INSERT INTO "user" (username, asset_image_id, password_hash, date_of_birth, phone_number, email) VALUES
                                                                                                     ('Chris', 3, '$2b$10$examplehashedpassword123456789012', '1990-05-15', '+1234567890', 'chris@example.com'),
                                                                                                     ('Alex', 4, '$2b$10$examplehashedpassword123456789013', '1985-08-20', '+0987654321', 'alex@example.com');

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