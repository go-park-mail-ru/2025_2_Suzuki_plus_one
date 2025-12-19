-- Al Pacino

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

-- Robin Williams

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



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Robin Williams' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Jumanji' LIMIT 1),
        'Alan Parrish'
    );

-- Kirsten Dunst

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



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Kirsten Dunst' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Jumanji' LIMIT 1),
        'Judy Shepherd'
    );


-- Bonnie Hunt

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



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Bonnie Hunt' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Jumanji' LIMIT 1),
        'Sarah Whittle'
    );


-- Jack Lemmon

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Jack_Lemmon.jpg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Jack_Lemmon.jpg' LIMIT 1), 400, 394);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Jack Lemmon', '1925-02-08', 'American actor known for classic films such as Some Like It Hot and The Apartment, and for his comedic partnership with Walter Matthau.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Jack Lemmon' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Jack_Lemmon.jpg' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Jack Lemmon' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Grumpier Old Men' LIMIT 1),
        'John Gustafson'
    );


-- Walter Matthau

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Walter_Matthau.jpg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Walter_Matthau.jpg' LIMIT 1), 675, 675);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Walter Matthau', '1920-10-01', 'American actor celebrated for his sharp comedic style and frequent collaborations with Jack Lemmon, including The Odd Couple.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Walter Matthau' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Walter_Matthau.jpg' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Walter Matthau' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Grumpier Old Men' LIMIT 1),
        'Max Goldman'
    );


-- Ann-Margret

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Ann-Margret.jpg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Ann-Margret.jpg' LIMIT 1), 235, 243);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Ann-Margret', '1941-04-28', 'Swedish-American actress, singer, and dancer known for roles in Bye Bye Birdie and Viva Las Vegas, and for her performance in Grumpier Old Men.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Ann-Margret' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Ann-Margret.jpg' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Ann-Margret' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Grumpier Old Men' LIMIT 1),
        'Ariel Truax'
    );


-- Whitney Houston

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Whitney_Houston.jpg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Whitney_Houston.jpg' LIMIT 1), 924, 1200);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Whitney Houston', '1963-08-09', 'American singer and actress who starred in The Bodyguard and Waiting to Exhale, and is one of the best-selling music artists of all time.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Whitney Houston' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Whitney_Houston.jpg' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Whitney Houston' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Waiting to Exhale' LIMIT 1),
        'Savannah Jackson'
    );


-- Angela Bassett

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Angela_Bassett.webp', 'image/webp', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Angela_Bassett.webp' LIMIT 1), 960, 640);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Angela Bassett', '1958-08-16', 'American actress known for powerful performances in films such as What’s Love Got to Do with It, Waiting to Exhale, and Black Panther.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Angela Bassett' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Angela_Bassett.webp' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Angela Bassett' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Waiting to Exhale' LIMIT 1),
        'Bernadine Harris'
    );


-- Loretta Devine

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Loretta_Devine.jpg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Loretta_Devine.jpg' LIMIT 1), 758, 758);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Loretta Devine', '1949-08-21', 'American actress and singer known for roles in Waiting to Exhale, Boston Public, and Grey’s Anatomy.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Loretta Devine' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Loretta_Devine.jpg' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Loretta Devine' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Waiting to Exhale' LIMIT 1),
        'Gloria Matthews'
    );


-- Steve Martin

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Steve_Martin.jpg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Steve_Martin.jpg' LIMIT 1), 500, 608);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Steve Martin', '1945-08-14', 'American actor, comedian, writer, and musician known for films such as Father of the Bride, Planes, Trains and Automobiles, and Roxanne.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Steve Martin' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Steve_Martin.jpg' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Steve Martin' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Father of the Bride Part II' LIMIT 1),
        'George Banks'
    );


-- Diane Keaton

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Diane_Keaton.webp', 'image/webp', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Diane_Keaton.webp' LIMIT 1), 3000, 2105);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Diane Keaton', '1946-01-05', 'American actress and filmmaker known for her collaborations with Woody Allen and films such as Annie Hall, The Godfather trilogy, and Father of the Bride.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Diane Keaton' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Diane_Keaton.webp' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Diane Keaton' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Father of the Bride Part II' LIMIT 1),
        'Nina Banks'
    );


-- Martin Short

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Martin_Short.webp', 'image/webp', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Martin_Short.webp' LIMIT 1), 1200, 630);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Martin Short', '1950-03-26', 'Canadian-American actor and comedian known for SCTV, Three Amigos, and his role as wedding planner Franck in the Father of the Bride films.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Martin Short' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Martin_Short.webp' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Martin Short' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Father of the Bride Part II' LIMIT 1),
        'Franck Eggelhoffer'
    );


-- Harrison Ford

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Harrison_Ford.webp', 'image/webp', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Harrison_Ford.webp' LIMIT 1), 216, 320);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Harrison Ford', '1942-07-13', 'American actor best known for iconic roles in Star Wars, Indiana Jones, and romantic lead performances such as in Sabrina.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Harrison Ford' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Harrison_Ford.webp' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Harrison Ford' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Sabrina' LIMIT 1),
        'Linus Larrabee'
    );


-- Julia Ormond

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Julia_Ormond.jpg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Julia_Ormond.jpg' LIMIT 1), 1050, 549);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Julia Ormond', '1965-01-04', 'English actress known for films such as Sabrina, Legends of the Fall, and First Knight.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Julia Ormond' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Julia_Ormond.jpg' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Julia Ormond' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Sabrina' LIMIT 1),
        'Sabrina Fairchild'
    );


-- Greg Kinnear

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Greg_Kinnear.jpeg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Greg_Kinnear.jpeg' LIMIT 1), 300, 171);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Greg Kinnear', '1963-06-17', 'American actor and former television host known for roles in Sabrina, As Good as It Gets, and Little Miss Sunshine.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Greg Kinnear' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Greg_Kinnear.jpeg' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Greg Kinnear' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Sabrina' LIMIT 1),
        'David Larrabee'
    );


-- Jonathan Taylor Thomas

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Jonathan_Taylor_Thomas.webp', 'image/webp', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Jonathan_Taylor_Thomas.webp' LIMIT 1), 780, 438);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Jonathan Taylor Thomas', '1981-09-08', 'American actor and voice actor who rose to fame in the 1990s, known for Home Improvement, The Lion King, and Tom and Huck.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Jonathan Taylor Thomas' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Jonathan_Taylor_Thomas.webp' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Jonathan Taylor Thomas' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Tom and Huck' LIMIT 1),
        'Tom Sawyer'
    );


-- Brad Renfro

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Brad_Renfro.webp', 'image/webp', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Brad_Renfro.webp' LIMIT 1), 640, 1017);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Brad Renfro', '1982-07-25', 'American actor known for powerful performances in films such as The Client, Sleepers, and Tom and Huck.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Brad Renfro' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Brad_Renfro.webp' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Brad Renfro' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Tom and Huck' LIMIT 1),
        'Huckleberry Finn'
    );


-- Eric Schweig

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Eric_Schweig.jpg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Eric_Schweig.jpg' LIMIT 1), 319, 449);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Eric Schweig', '1967-06-19', 'Canadian actor known for roles in Tom and Huck, The Last of the Mohicans, and Big Eden.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Eric Schweig' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Eric_Schweig.jpg' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Eric Schweig' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Tom and Huck' LIMIT 1),
        'Injun Joe'
    );


-- Jean-Claude Van Damme

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Jean-Claude_Van_Damme.jpg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Jean-Claude_Van_Damme.jpg' LIMIT 1), 1000, 1500);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Jean-Claude Van Damme', '1960-10-18', 'Belgian actor and martial artist known for action films such as Bloodsport, Kickboxer, and Sudden Death.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Jean-Claude Van Damme' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Jean-Claude_Van_Damme.jpg' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Jean-Claude Van Damme' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Sudden Death' LIMIT 1),
        'Darren McCord'
    );


-- Powers Boothe

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Powers_Boothe.jpg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Powers_Boothe.jpg' LIMIT 1), 605, 469);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Powers Boothe', '1948-06-01', 'American actor known for intense roles in films and television, including Sudden Death, Tombstone, and Deadwood.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Powers Boothe' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Powers_Boothe.jpg' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Powers Boothe' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Sudden Death' LIMIT 1),
        'Joshua Foss'
    );


-- Raymond J. Barry

INSERT INTO asset (s3_key, mime_type, file_size_mb) VALUES
    ('actors/Raymond_J._Barry.jpg', 'image/jpeg', 2);

INSERT INTO asset_image (asset_id, resolution_width, resolution_height) VALUES
    ((SELECT asset_id FROM asset WHERE s3_key = 'actors/Raymond_J._Barry.jpg' LIMIT 1), 800, 576);

INSERT INTO actor (name, birth_date, bio) VALUES
    ('Raymond J. Barry', '1939-08-31', 'American actor known for roles in film and television, including Sudden Death and later appearances in projects such as Justified and Better Call Saul.');

INSERT INTO actor_image (actor_id, asset_image_id, image_type) VALUES
    ((SELECT actor_id FROM actor WHERE name = 'Raymond J. Barry' LIMIT 1),
     (SELECT asset_image_id FROM asset_image WHERE asset_id = (SELECT asset_id FROM asset WHERE s3_key = 'actors/Raymond_J._Barry.jpg' LIMIT 1) LIMIT 1),
     'profile');



INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (
        (SELECT actor_id FROM actor WHERE name = 'Raymond J. Barry' LIMIT 1),
        (SELECT media_id FROM media WHERE title = 'Sudden Death' LIMIT 1),
        'Vice President'
    );


