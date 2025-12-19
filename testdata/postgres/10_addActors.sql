-- Al Pacino
BEGIN;
INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/alpc.jpeg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/alpc.jpeg' LIMIT 1), 266, 400);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Al Pacino', '1940-04-25', 'American actor and filmmaker, best known for The Godfather series.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Al Pacino' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/alpc.jpeg' LIMIT 1) LIMIT 1),
     'profile');


INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Al Pacino' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Heat' LIMIT 1),
        'Michael Corleone'
    );
COMMIT;
-- Robin Williams
BEGIN;
INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/GoOnlineTools-image-downloader.jpeg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/GoOnlineTools-image-downloader.jpeg' LIMIT 1), 1370, 2048);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Robin Williams', '1951-07-21', 'American actor and comedian, known for films like Jumanji, Dead Poets Society, and Good Will Hunting.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Robin Williams' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/GoOnlineTools-image-downloader.jpeg' LIMIT 1) LIMIT 1),
     'profile');
COMMIT;

BEGIN;
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Robin Williams' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Jumanji' LIMIT 1),
        'Alan Parrish'
    );
COMMIT;
-- Kirsten Dunst
BEGIN;
INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Kirsten_Dunst.jpeg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Kirsten_Dunst.jpeg' LIMIT 1), 1370, 2048);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Kirsten Dunst', '1982-04-30', 'American actress who rose to fame as a child star and is known for films such as Jumanji and Interview with the Vampire.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Kirsten Dunst' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Kirsten_Dunst.jpeg' LIMIT 1) LIMIT 1),
     'profile');
COMMIT;

BEGIN;
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Kirsten Dunst' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Jumanji' LIMIT 1),
        'Judy Shepherd'
    );
COMMIT;

-- Bonnie Hunt
BEGIN;
INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Bonnie_Hunt.jpeg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Bonnie_Hunt.jpeg' LIMIT 1), 3219, 4829);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Bonnie Hunt', '1961-09-22', 'American actress, comedian, writer, and director known for her roles in Jumanji, Jerry Maguire, and Cheaper by the Dozen.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Bonnie Hunt' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Bonnie_Hunt.jpeg' LIMIT 1) LIMIT 1),
     'profile');
COMMIT;

BEGIN;
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Bonnie Hunt' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Jumanji' LIMIT 1),
        'Sarah Whittle'
    );
COMMIT;

