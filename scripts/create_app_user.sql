\echo Creating app role :app_user in database :dbname

-- 1. На dev-окружении удобно просто пересоздавать роль
DROP ROLE IF EXISTS :"app_user";

-- 2. Создаём роль приложения
CREATE ROLE :"app_user" WITH
    LOGIN
    PASSWORD :'app_password'
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    NOINHERIT
    NOREPLICATION;

-- 3. search_path только на public
ALTER ROLE :"app_user" SET search_path TO public;

-- Таймауты на уровне роли приложения
ALTER ROLE :"app_user" SET statement_timeout = '2s';
ALTER ROLE :"app_user" SET lock_timeout = '200ms';

-- 4. Доступ к базе и схеме
GRANT CONNECT ON DATABASE :dbname TO :"app_user";
GRANT USAGE ON SCHEMA public TO :"app_user";

-- =====================================================================
-- ПРАВА НА ТАБЛИЦЫ
-- =====================================================================

-- ---------------------------------------------------------------------
-- Контентные таблицы — только SELECT
-- Используются в:
--   GET /media/*
--   GET /genre/*
--   GET /actor/*
--   GET /object
--   GET /media/watch
--   GET /search
-- ---------------------------------------------------------------------
GRANT SELECT ON
    genre,
    media,
    media_genre,
    media_episode,
    media_image,
    media_video,
    media_audio,
    media_subtitle,
    asset,
    asset_video,
    asset_audio,
    asset_subtitle,
    asset_image,
    actor,
    actor_image,
    actor_role
TO :"app_user";

-- ---------------------------------------------------------------------
-- Пользователь, авторизация, сессии
--   POST /auth/signin
--   POST /auth/signup
--   GET  /auth/refresh
--   GET  /auth/signout
--   GET  /user/me
--   POST /user/me/update
--   POST /user/me/update/avatar
--   POST /user/me/update/password
-- ---------------------------------------------------------------------

-- Пользователь: читать себя, регистрировать, менять профиль/пароль
GRANT SELECT, INSERT, UPDATE ON
    "user"
TO :"app_user";

-- Сессии: создать, читать, продлевать, удалять
GRANT SELECT, INSERT, UPDATE, DELETE ON
    user_session
TO :"app_user";

-- Аватарки: создавать новые asset/asset_image, читать существующие
GRANT SELECT, INSERT ON
    asset,
    asset_image
TO :"app_user";

-- ---------------------------------------------------------------------
-- Лайки и история просмотра
--   GET /media/{id}/like
--   PUT /media/{id}/like
--   DELETE /media/{id}/like
--   GET /media/my
--   GET /media/recommendations
--   GET /media/watch
-- ---------------------------------------------------------------------

-- Лайки на media
GRANT SELECT, INSERT, UPDATE, DELETE ON
    user_like_media
TO :"app_user";

-- История просмотра
GRANT SELECT, INSERT ON
    user_watch_history
TO :"app_user";

-- Для рекомендаций могут использоваться и эти таблицы (только SELECT)
GRANT SELECT ON
    user_rating_media,
    user_like_actor,
    user_like_playlist
TO :"app_user";

-- ---------------------------------------------------------------------
-- Обращения в поддержку (appeal)
--   GET  /appeal/all
--   GET  /appeal/my
--   POST /appeal/new
--   GET  /appeal/{id}
--   PUT  /appeal/{id}/resolve
--   POST /appeal/{id}/message
--   GET  /appeal/{id}/message
-- ---------------------------------------------------------------------
GRANT SELECT, INSERT, UPDATE ON
    user_appeal
TO :"app_user";

GRANT SELECT, INSERT ON
    user_appeal_message
TO :"app_user";

-- ---------------------------------------------------------------------
-- Плейлисты и комментарии — сейчас эндпоинты их явно не меняют,
-- но могут использоваться в выборках (/object, /search и т.п.)
-- ---------------------------------------------------------------------
GRANT SELECT ON
    playlist,
    playlist_media,
    user_playlist,
    user_comment_media,
    user_comment_actor
TO :"app_user";

-- =====================================================================
-- ПРАВА НА SEQUENCE
-- =====================================================================

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO :"app_user";

\echo Done creating app role :app_user
