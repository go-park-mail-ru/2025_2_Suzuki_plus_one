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

-- Correct genre associations for each movie
INSERT INTO media_genre (media_id, genre_id) VALUES
-- Toy Story (Adventure, Fantasy, Comedy* - но Comedy нет в списке)
(1, 5), (1, 6),  -- Adventure, Fantasy

-- Jumanji (Adventure, Fantasy, Action)
(2, 5), (2, 6), (2, 1), -- Adventure, Fantasy, Action

-- Grumpier Old Men (Comedy* - но Comedy нет, ближайший Drama)
(3, 4), -- Drama

-- Waiting to Exhale (Drama)
(4, 4), -- Drama

-- Father of the Bride Part II (Comedy* - но Comedy нет, ближайший Drama)
(5, 4), -- Drama

-- Heat (Action, Thriller, Drama)
(6, 1), (6, 3), (6, 4), -- Action, Thriller, Drama

-- Sabrina (Romance* - но Romance нет, ближайший Drama)
(7, 4), -- Drama

-- Tom and Huck (Adventure)
(8, 5), -- Adventure

-- Sudden Death (Action, Thriller)
(9, 1), (9, 3), -- Action, Thriller

-- GoldenEye (Action, Adventure, Thriller)
(10, 1), (10, 5), (10, 3); -- Action, Adventure, Thriller

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
-- Add actors
('/actors/albert_brooks.jpg', 'image/jpeg', 0.64),
('/actors/alexander_gould.jpg', 'image/jpeg', 4.67),
('/actors/al_pacino.jpg', 'image/jpeg', 0.0065),
('/actors/anna_torv.jpg', 'image/jpeg', 0.0144),
('/actors/bella_ramsey.jpg', 'image/jpeg', 0.81),
('/actors/brad_renfro.jpg', 'image/jpeg', 0.03),
('/actors/carrie-anne_moss.jpg', 'image/jpeg', 0.55),
('/actors/daniel_radcliffe.jpg', 'image/jpeg', 0.19),
('/actors/david_harbour.jpg', 'image/jpeg', 0.85),
('/actors/ellen_degeneres.jpg', 'image/jpeg', 0.072),
('/actors/elliot_page.jpg', 'image/jpeg', 0.076),
('/actors/emma_watson.jpg', 'image/jpeg', 0.12),
('/actors/eric_schweig.jpg', 'image/jpeg', 0.027),
('/actors/greg_kinnear.jpg', 'image/jpeg', 0.0088),
('/actors/harrison_ford.jpg', 'image/jpeg', 0.0062),
('/actors/izabella_scorupco.jpg', 'image/jpeg', 0.0063),
('/actors/jean-claude_van_damme.jpg', 'image/jpeg', 0.022),
('/actors/john_wood.jpg', 'image/jpeg', 0.021),
('/actors/jonathan_taylor_thomas.jpg', 'image/jpeg', 0.0097),
('/actors/jon_voight.jpg', 'image/jpeg', 0.0057),
('/actors/joseph_gordon-levitt.jpg', 'image/jpeg', 0.75),
('/actors/julia_ormond.jpg', 'image/jpeg', 0.0074),
('/actors/keanu_reeves.jpg', 'image/jpeg', 0.84),
('/actors/laurence_fishburne.jpg', 'image/jpeg', 0.031),
('/actors/leonardo_dicaprio.jpg', 'image/jpeg', 0.034),
('/actors/merle_dandridge.jpg', 'image/jpeg', 0.30),
('/actors/millie_bobby_brown.jpg', 'image/jpeg', 0.48),
('/actors/nancy_marchand.jpg', 'image/jpeg', 0.0087),
('/actors/pedro_pascal.jpg', 'image/jpeg', 0.21),
('/actors/pierce_brosnan.jpg', 'image/jpeg', 0.0083),
('/actors/powers_boothe.jpg', 'image/jpeg', 0.020),
('/actors/raymond_j._barry.jpg', 'image/jpeg', 0.020),
('/actors/robert_de_niro.jpg', 'image/jpeg', 0.0066),
('/actors/rupert_grint.jpg', 'image/jpeg', 1.37),
('/actors/sean_bean.jpg', 'image/jpeg', 1.96),
('/actors/tom_sizemore.jpg', 'image/jpeg', 0.0081),
('/actors/val_kilmer.jpg', 'image/jpeg', 0.0075),
('/actors/winona_ryder.jpg', 'image/jpeg', 0.47),

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
    ('Harrison Ford', '1942-07-13', 'Harrison Ford is an American actor.'),
-- add actors
  ('Albert Brooks', '1947-07-22', 'Albert Brooks (born Albert Lawrence Einstein; July 22, 1947) is an American actor, comedian, director and screenwriter.'),
    ('Alexander Gould', '1994-05-04', 'Alexander Jerome Gould (born May 4, 1994) is an American actor and United States Marine.'),
    ('Anna Torv', '1979-06-07', 'Anna Torv (born 7 June 1979) is an Australian actress.'),
    ('Bella Ramsey', '2003-09-25', 'Bella Ramsey (born 25 September 2003) is an English actor.'),
    ('Brad Renfro', '1982-07-25', 'Brad Barron Renfro (July 25, 1982 – January 15, 2008) was an American actor.'),
    ('Carrie-Anne Moss', '1967-08-21', 'Carrie-Anne Moss (born August 21, 1967) is a Canadian actress.'),
    ('Daniel Radcliffe', '1989-07-23', 'Daniel Jacob Radcliffe (born 23 July 1989) is an English actor.'),
    ('David Harbour', '1975-04-10', 'David Harbour (born April 10, 1975) is an American actor.'),
    ('Ellen DeGeneres', '1958-01-26', 'Ellen Lee DeGeneres (born January 26, 1958) is an American comedian, television host, actress, writer, and producer.'),
    ('Elliot Page', '1987-02-21', 'Elliot Page (born February 21, 1987) is a Canadian actor and producer.'),
    ('Emma Watson', '1990-04-15', 'Emma Charlotte Duerre Watson (born 15 April 1990) is an English actress and activist.'),
    ('Eric Schweig', '1967-06-19', 'Eric Schweig (born June 19, 1967) is a Canadian film and television actor.'),
    ('Greg Kinnear', '1963-06-17', 'Gregory Buck Kinnear (born June 17, 1963) is an American actor and television personality.'),
    ('Izabella Scorupco', '1970-06-04', 'Izabella Dorota Scorupco (born 4 June 1970) is a Polish actress, model and singer.'),
    ('Jean-Claude Van Damme', '1960-10-18', 'Jean-Claude Van Damme (born Jean-Claude Camille François Van Varenberg; 18 October 1960) is a Belgian martial artist and actor.'),
    ('John Wood', '1930-07-05', 'John Wood (5 July 1930 – 6 August 2011) was an English actor.'),
    ('Jonathan Taylor Thomas', '1981-09-08', 'Jonathan Taylor Thomas (born September 8, 1981) is an American actor and director.'),
    ('Jon Voight', '1938-12-29', 'Jonathan Vincent Voight (born December 29, 1938) is an American actor.'),
    ('Joseph Gordon-Levitt', '1981-02-17', 'Joseph Leonard Gordon-Levitt (born February 17, 1981) is an American actor.'),
    ('Julia Ormond', '1965-01-04', 'Julia Karin Ormond (born 4 January 1965) is an English actress.'),
    ('Keanu Reeves', '1964-09-02', 'Keanu Charles Reeves (born September 2, 1964) is a Canadian actor.'),
    ('Laurence Fishburne', '1961-07-30', 'Laurence John Fishburne III (born July 30, 1961) is an American actor.'),
    ('Leonardo DiCaprio', '1974-11-11', 'Leonardo Wilhelm DiCaprio (born November 11, 1974) is an American actor and film producer.'),
    ('Merle Dandridge', '1975-05-31', 'Merle Dandridge (born May 31, 1975) is an American actress and singer.'),
    ('Millie Bobby Brown', '2004-02-19', 'Millie Bobby Brown (born 19 February 2004) is a British actress and model.'),
    ('Nancy Marchand', '1928-06-19', 'Nancy Marchand (June 19, 1928 – June 18, 2000) was an American actress.'),
    ('Pedro Pascal', '1975-04-02', 'Pedro Pascal (born José Pedro Balmaceda Pascal; April 2, 1975) is a Chilean-American actor.'),
    ('Pierce Brosnan', '1953-05-16', 'Pierce Brendan Brosnan (born 16 May 1953) is an Irish actor and film producer.'),
    ('Powers Boothe', '1948-06-01', 'Powers Allen Boothe (June 1, 1948 – May 14, 2017) was an American actor.'),
    ('Raymond J. Barry', '1939-03-14', 'Raymond J. Barry (born March 14, 1939) is an American character actor.'),
    ('Rupert Grint', '1988-08-24', 'Rupert Alexander Lloyd Grint (born 24 August 1988) is an English actor.'),
    ('Sean Bean', '1959-04-17', 'Shaun Mark Bean (born 17 April 1959), known professionally as Sean Bean, is an English actor.'),
    ('Tom Sizemore', '1961-11-29', 'Thomas Edward Sizemore Jr. (November 29, 1961 – March 3, 2023) was an American actor.'),
    ('Val Kilmer', '1959-12-31', 'Val Edward Kilmer (born December 31, 1959) is an American actor.'),
    ('Winona Ryder', '1971-10-29', 'Winona Laura Horowitz (born October 29, 1971), known professionally as Winona Ryder, is an American actress.');


-- Link actor images

CREATE OR REPLACE FUNCTION normalize_name(name text) RETURNS text AS $$
    BEGIN
		RETURN LOWER(REGEXP_REPLACE(name, '_', ' ', 'g'));
	END
$$ LANGUAGE plpgsql IMMUTABLE;

INSERT INTO actor_image (actor_id, asset_image_id, image_type)
SELECT 
	act.actor_id,
	ai.asset_image_id,
	'profile'
FROM actor act
JOIN asset  on normalize_name(asset.s3_key) LIKE '%' || normalize_name(act.name) || '%'
JOIN asset_image ai on ai.asset_id = asset.asset_id;

-- Не актуально

-- INSERT INTO actor_image (actor_id, asset_image_id, image_type)
-- SELECT 1, asset_image_id, 'profile'
-- FROM asset_image ai
--          JOIN asset a ON ai.asset_id = a.asset_id
-- WHERE a.s3_key = '/actors/Tom_Hanks.png';

-- INSERT INTO actor_image (actor_id, asset_image_id, image_type)
-- SELECT 2, asset_image_id, 'profile'
-- FROM asset_image ai
--          JOIN asset a ON ai.asset_id = a.asset_id
-- WHERE a.s3_key = '/actors/Tim_Allen.png';

-- INSERT INTO actor_image (actor_id, asset_image_id, image_type)
-- SELECT 3, asset_image_id, 'profile'
-- FROM asset_image ai
--          JOIN asset a ON ai.asset_id = a.asset_id
-- WHERE a.s3_key = '/actors/Robert_De_Niro.png';

-- INSERT INTO actor_image (actor_id, asset_image_id, image_type)
-- SELECT 4, asset_image_id, 'profile'
-- FROM asset_image ai
--          JOIN asset a ON ai.asset_id = a.asset_id
-- WHERE a.s3_key = '/actors/Al_Pacino.png';

-- INSERT INTO actor_image (actor_id, asset_image_id, image_type)
-- SELECT 5, asset_image_id, 'profile'
-- FROM asset_image ai
--          JOIN asset a ON ai.asset_id = a.asset_id
-- WHERE a.s3_key = '/actors/Harrison_Ford.png';

-- Insert actor roles


-- Перенесено в файл 07_add_actor_roles.sql
-- INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
--                                                            (1, 1, 'Woody (voice)'),
--                                                            (2, 1, 'Buzz Lightyear (voice)'),
--                                                            (3, 6, 'Neil McCauley'),
--                                                            (4, 6, 'Lt. Vincent Hanna'),
--                                                            (5, 10, 'James Bond');



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

INSERT INTO user_like_media (user_id, media_id, is_dislike) VALUES
                                                    (2, 1, false), (2, 2, false),
                                                    (3, 1, false);

INSERT INTO user_like_actor (user_id, actor_id, is_dislike) VALUES
                                                    (2, 1, false), (2, 2, false);

INSERT INTO user_like_playlist (user_id, playlist_id, is_dislike) VALUES
    (3, 1, false);

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