CREATE OR REPLACE FUNCTION get_actor_id_by_name(actor_name text)
RETURNS integer AS $$
DECLARE
    actor_id integer;
BEGIN
    SELECT a.actor_id INTO actor_id
    FROM actor a
    WHERE a.name = actor_name
    LIMIT 1; 
    RETURN actor_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_media_id_by_title(media_title text)
RETURNS integer AS $$
DECLARE
    media_id integer;
BEGIN
    SELECT m.media_id INTO media_id
    FROM media m
    WHERE m.title = media_title
    LIMIT 1;
    RETURN media_id;
END;
$$ LANGUAGE plpgsql;

-- Insert helper that only inserts when both actor and media exist
CREATE OR REPLACE FUNCTION insert_actor_role(actor_name text, media_title text, role_name text)
RETURNS void AS $$
DECLARE
    a_id integer;
    m_id integer;
BEGIN
    SELECT actor_id INTO a_id FROM actor WHERE name = actor_name LIMIT 1;
    SELECT media_id INTO m_id FROM media WHERE title = media_title LIMIT 1;

    IF a_id IS NULL THEN
        RAISE NOTICE 'Skipping insert: actor "%" not found', actor_name;
        RETURN;
    END IF;
    IF m_id IS NULL THEN
        RAISE NOTICE 'Skipping insert: media "%" not found', media_title;
        RETURN;
    END IF;

    INSERT INTO actor_role (actor_id, media_id, role_name) VALUES (a_id, m_id, role_name);
END;
$$ LANGUAGE plpgsql;

-- Toy Story (1995)
SELECT insert_actor_role('Tom Hanks', 'Toy Story', 'Woody');
SELECT insert_actor_role('Tim Allen', 'Toy Story', 'Buzz Lightyear');

-- Heat (1995)
SELECT insert_actor_role('Robert De Niro', 'Heat', 'Neil McCauley');
SELECT insert_actor_role('Al Pacino', 'Heat', 'Lt. Vincent Hanna');
SELECT insert_actor_role('Val Kilmer', 'Heat', 'Chris Shiherlis');
SELECT insert_actor_role('Tom Sizemore', 'Heat', 'Michael Cheritto');

-- GoldenEye (1995)
SELECT insert_actor_role('Pierce Brosnan', 'GoldenEye', 'James Bond');
SELECT insert_actor_role('Izabella Scorupco', 'GoldenEye', 'Natalya Simonova');
SELECT insert_actor_role('Sean Bean', 'GoldenEye', 'Alec Trevelyan / 006');

-- The Matrix (1999)
SELECT insert_actor_role('Keanu Reeves', 'The Matrix', 'Neo');
SELECT insert_actor_role('Laurence Fishburne', 'The Matrix', 'Morpheus');
SELECT insert_actor_role('Carrie-Anne Moss', 'The Matrix', 'Trinity');

-- The Matrix Reloaded (2003)
SELECT insert_actor_role('Keanu Reeves', 'The Matrix Reloaded', 'Neo');
SELECT insert_actor_role('Laurence Fishburne', 'The Matrix Reloaded', 'Morpheus');
SELECT insert_actor_role('Carrie-Anne Moss', 'The Matrix Reloaded', 'Trinity');

-- Harry Potter and the Chamber of Secrets (2002)
SELECT insert_actor_role('Daniel Radcliffe', 'Harry Potter and the Chamber of Secrets', 'Harry Potter');
SELECT insert_actor_role('Emma Watson', 'Harry Potter and the Chamber of Secrets', 'Hermione Granger');
SELECT insert_actor_role('Rupert Grint', 'Harry Potter and the Chamber of Secrets', 'Ron Weasley');
SELECT insert_actor_role('John Wood', 'Harry Potter and the Chamber of Secrets', 'Professor Gilderoy Lockhart');

-- Harry Potter and the Prisoner of Azkaban (2004)
SELECT insert_actor_role('Daniel Radcliffe', 'Harry Potter and the Prisoner of Azkaban', 'Harry Potter');
SELECT insert_actor_role('Emma Watson', 'Harry Potter and the Prisoner of Azkaban', 'Hermione Granger');
SELECT insert_actor_role('Rupert Grint', 'Harry Potter and the Prisoner of Azkaban', 'Ron Weasley');

-- Harry Potter and the Goblet of Fire (2005)
SELECT insert_actor_role('Daniel Radcliffe', 'Harry Potter and the Goblet of Fire', 'Harry Potter');
SELECT insert_actor_role('Emma Watson', 'Harry Potter and the Goblet of Fire', 'Hermione Granger');
SELECT insert_actor_role('Rupert Grint', 'Harry Potter and the Goblet of Fire', 'Ron Weasley');
SELECT insert_actor_role('Robert De Niro', 'Harry Potter and the Goblet of Fire', 'Uncle Vernon (предположительно, но в списке он есть, возможно для другого фильма)');

-- Harry Potter and the Order of the Phoenix (2007)
SELECT insert_actor_role('Daniel Radcliffe', 'Harry Potter and the Order of the Phoenix', 'Harry Potter');
SELECT insert_actor_role('Emma Watson', 'Harry Potter and the Order of the Phoenix', 'Hermione Granger');
SELECT insert_actor_role('Rupert Grint', 'Harry Potter and the Order of the Phoenix', 'Ron Weasley');

-- Inception (2010)
SELECT insert_actor_role('Leonardo DiCaprio', 'Inception', 'Dom Cobb');
SELECT insert_actor_role('Elliot Page', 'Inception', 'Ariadne');
SELECT insert_actor_role('Joseph Gordon-Levitt', 'Inception', 'Arthur');

-- Finding Nemo (2003)
SELECT insert_actor_role('Albert Brooks', 'Finding Nemo', 'Marlin');
SELECT insert_actor_role('Ellen DeGeneres', 'Finding Nemo', 'Dory');
SELECT insert_actor_role('Alexander Gould', 'Finding Nemo', 'Nemo');

-- Stranger Things (сериал) - Основной сериал
SELECT insert_actor_role('Winona Ryder', 'Stranger Things', 'Joyce Byers');
SELECT insert_actor_role('David Harbour', 'Stranger Things', 'Jim Hopper');
SELECT insert_actor_role('Millie Bobby Brown', 'Stranger Things', 'Eleven / Jane Hopper');

-- The Last of Us (сериал) - Основной сериал
SELECT insert_actor_role('Pedro Pascal', 'The Last of Us', 'Joel Miller');
SELECT insert_actor_role('Bella Ramsey', 'The Last of Us', 'Ellie');
SELECT insert_actor_role('Merle Dandridge', 'The Last of Us', 'Marlene');
SELECT insert_actor_role('Anna Torv', 'The Last of Us', 'Tess');

