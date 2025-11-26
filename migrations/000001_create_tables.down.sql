-- Drop tables in reverse order of creation to handle dependencies

-- First drop all triggers and functions
DROP TRIGGER IF EXISTS trg_check_media_episode_types ON media_episode;
DROP FUNCTION IF EXISTS check_media_episode_types();

-- Drop user-related tables (with foreign keys) - must drop user_playlist BEFORE the enum
DROP TABLE IF EXISTS user_rating_media;
DROP TABLE IF EXISTS user_comment_actor;
DROP TABLE IF EXISTS user_comment_media;
DROP TABLE IF EXISTS user_like_playlist;
DROP TABLE IF EXISTS user_like_actor;
DROP TABLE IF EXISTS user_like_media;
DROP TABLE IF EXISTS user_watch_history;
DROP TABLE IF EXISTS user_playlist;  -- This table uses the playlist_role enum
DROP TABLE IF EXISTS playlist_media;
DROP TABLE IF EXISTS playlist;
DROP TABLE IF EXISTS user_session;
DROP TABLE IF EXISTS "user";

-- Now we can safely drop the enum type
DROP TYPE IF EXISTS playlist_role;

-- Drop actor-related tables
DROP TABLE IF EXISTS actor_role;
DROP TABLE IF EXISTS actor_image;
DROP TABLE IF EXISTS actor;

-- Drop media-asset relationship tables
DROP TABLE IF EXISTS media_subtitle;
DROP TABLE IF EXISTS media_audio;
DROP TABLE IF EXISTS media_video_asset;
DROP TABLE IF EXISTS media_video;
DROP TABLE IF EXISTS media_image;

-- Drop asset subtype tables
DROP TABLE IF EXISTS asset_subtitle;
DROP TABLE IF EXISTS asset_audio;
DROP TABLE IF EXISTS asset_video;
DROP TABLE IF EXISTS asset_image;
DROP TABLE IF EXISTS asset;

-- Drop media structure tables
DROP TABLE IF EXISTS media_episode;
DROP TABLE IF EXISTS media_genre;

-- Drop core tables
DROP TABLE IF EXISTS genre;
DROP TABLE IF EXISTS media;