BEGIN;

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
END $$;