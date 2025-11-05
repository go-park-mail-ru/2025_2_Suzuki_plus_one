INSERT INTO asset (s3_key, mime_type, size_mb)
VALUES
('actors/leo.png', 'image/png', 1.2),
('actors/morgan.png', 'image/png', 1.1);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height)
VALUES
((SELECT asset_id FROM asset WHERE s3_key = 'actors/leo.png'), 800, 800),
((SELECT asset_id FROM asset WHERE s3_key = 'actors/morgan.png'), 800, 800);

INSERT INTO actor (name, birth_date, bio)
VALUES
('Leonardo DiCaprio', '1974-11-11', 'Leonardo Wilhelm DiCaprio is an American actor and film producer.'),
('Morgan Freeman', '1937-06-01', 'Morgan Freeman is an American actor, director, and narrator.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type)
VALUES
(1, (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/leo.png')), 'profile'),
(2, (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/morgan.png')), 'profile');

INSERT INTO actor_role (actor_id, media_id, role_name)
VALUES
(1, 1, 'Dom Cobb'),
(2, 2, 'Red'),
(1, 2, 'Hugh Glass');