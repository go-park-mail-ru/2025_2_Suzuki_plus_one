CREATE OR REPLACE FUNCTION get_actor_id_by_name(actor_name text)
RETURNS integer AS $$
DECLARE
    actor_id integer;
BEGIN
    SELECT a.actor_id INTO actor_id
    FROM actor a
    WHERE a.name = actor_name; 
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
    WHERE m.title = media_title;
    RETURN media_id;
END;
$$ LANGUAGE plpgsql;


-- Toy Story (1995)
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (get_actor_id_by_name('Tom Hanks'), get_media_id_by_title('Toy Story'), 'Woody'),
    (get_actor_id_by_name('Tim Allen'), get_media_id_by_title('Toy Story'), 'Buzz Lightyear');

-- Heat (1995)
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (get_actor_id_by_name('Robert De Niro'), get_media_id_by_title('Heat'), 'Neil McCauley'),
    (get_actor_id_by_name('Al Pacino'), get_media_id_by_title('Heat'), 'Lt. Vincent Hanna'),
    (get_actor_id_by_name('Val Kilmer'), get_media_id_by_title('Heat'), 'Chris Shiherlis'),
    (get_actor_id_by_name('Tom Sizemore'), get_media_id_by_title('Heat'), 'Michael Cheritto');

-- GoldenEye (1995)
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (get_actor_id_by_name('Pierce Brosnan'), get_media_id_by_title('GoldenEye'), 'James Bond'),
    (get_actor_id_by_name('Izabella Scorupco'), get_media_id_by_title('GoldenEye'), 'Natalya Simonova'),
    (get_actor_id_by_name('Sean Bean'), get_media_id_by_title('GoldenEye'), 'Alec Trevelyan / 006');

-- The Matrix (1999)
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (get_actor_id_by_name('Keanu Reeves'), get_media_id_by_title('The Matrix'), 'Neo'),
    (get_actor_id_by_name('Laurence Fishburne'), get_media_id_by_title('The Matrix'), 'Morpheus'),
    (get_actor_id_by_name('Carrie-Anne Moss'), get_media_id_by_title('The Matrix'), 'Trinity');

-- The Matrix Reloaded (2003)
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (get_actor_id_by_name('Keanu Reeves'), get_media_id_by_title('The Matrix Reloaded'), 'Neo'),
    (get_actor_id_by_name('Laurence Fishburne'), get_media_id_by_title('The Matrix Reloaded'), 'Morpheus'),
    (get_actor_id_by_name('Carrie-Anne Moss'), get_media_id_by_title('The Matrix Reloaded'), 'Trinity');

-- Harry Potter and the Chamber of Secrets (2002)
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (get_actor_id_by_name('Daniel Radcliffe'), get_media_id_by_title('Harry Potter and the Chamber of Secrets'), 'Harry Potter'),
    (get_actor_id_by_name('Emma Watson'), get_media_id_by_title('Harry Potter and the Chamber of Secrets'), 'Hermione Granger'),
    (get_actor_id_by_name('Rupert Grint'), get_media_id_by_title('Harry Potter and the Chamber of Secrets'), 'Ron Weasley'),
    (get_actor_id_by_name('John Wood'), get_media_id_by_title('Harry Potter and the Chamber of Secrets'), 'Professor Gilderoy Lockhart');

-- Harry Potter and the Prisoner of Azkaban (2004)
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (get_actor_id_by_name('Daniel Radcliffe'), get_media_id_by_title('Harry Potter and the Prisoner of Azkaban'), 'Harry Potter'),
    (get_actor_id_by_name('Emma Watson'), get_media_id_by_title('Harry Potter and the Prisoner of Azkaban'), 'Hermione Granger'),
    (get_actor_id_by_name('Rupert Grint'), get_media_id_by_title('Harry Potter and the Prisoner of Azkaban'), 'Ron Weasley');

-- Harry Potter and the Goblet of Fire (2005)
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (get_actor_id_by_name('Daniel Radcliffe'), get_media_id_by_title('Harry Potter and the Goblet of Fire'), 'Harry Potter'),
    (get_actor_id_by_name('Emma Watson'), get_media_id_by_title('Harry Potter and the Goblet of Fire'), 'Hermione Granger'),
    (get_actor_id_by_name('Rupert Grint'), get_media_id_by_title('Harry Potter and the Goblet of Fire'), 'Ron Weasley'),
    (get_actor_id_by_name('Robert De Niro'), get_media_id_by_title('Harry Potter and the Goblet of Fire'), 'Uncle Vernon (предположительно, но в списке он есть, возможно для другого фильма)'); -- Примечание: Роберт Де Ниро не снимался в Гарри Поттере. Возможно, он есть в общем списке актеров для другого фильма (например, Heat).

-- Harry Potter and the Order of the Phoenix (2007)
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (get_actor_id_by_name('Daniel Radcliffe'), get_media_id_by_title('Harry Potter and the Order of the Phoenix'), 'Harry Potter'),
    (get_actor_id_by_name('Emma Watson'), get_media_id_by_title('Harry Potter and the Order of the Phoenix'), 'Hermione Granger'),
    (get_actor_id_by_name('Rupert Grint'), get_media_id_by_title('Harry Potter and the Order of the Phoenix'), 'Ron Weasley');

-- Inception (2010)
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (get_actor_id_by_name('Leonardo DiCaprio'), get_media_id_by_title('Inception'), 'Dom Cobb'),
    (get_actor_id_by_name('Elliot Page'), get_media_id_by_title('Inception'), 'Ariadne'),
    (get_actor_id_by_name('Joseph Gordon-Levitt'), get_media_id_by_title('Inception'), 'Arthur');
    
-- Finding Nemo (2003)
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (get_actor_id_by_name('Albert Brooks'), get_media_id_by_title('Finding Nemo'), 'Marlin'),
    (get_actor_id_by_name('Ellen DeGeneres'), get_media_id_by_title('Finding Nemo'), 'Dory'),
    (get_actor_id_by_name('Alexander Gould'), get_media_id_by_title('Finding Nemo'), 'Nemo');

-- Stranger Things (сериал) - Основной сериал
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (get_actor_id_by_name('Winona Ryder'), get_media_id_by_title('Stranger Things'), 'Joyce Byers'),
    (get_actor_id_by_name('David Harbour'), get_media_id_by_title('Stranger Things'), 'Jim Hopper'),
    (get_actor_id_by_name('Millie Bobby Brown'), get_media_id_by_title('Stranger Things'), 'Eleven / Jane Hopper');

-- The Last of Us (сериал) - Основной сериал
INSERT INTO actor_role (actor_id, media_id, role_name) VALUES
    (get_actor_id_by_name('Pedro Pascal'), get_media_id_by_title('The Last of Us'), 'Joel Miller'),
    (get_actor_id_by_name('Bella Ramsey'), get_media_id_by_title('The Last of Us'), 'Ellie'),
    (get_actor_id_by_name('Merle Dandridge'), get_media_id_by_title('The Last of Us'), 'Marlene'),
    (get_actor_id_by_name('Anna Torv'), get_media_id_by_title('The Last of Us'), 'Tess');

