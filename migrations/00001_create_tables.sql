-- Initial schema for PopFilms (aligned with docs/database/base.md ER diagram)

CREATE TABLE "user" (
    user_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Playlists
CREATE TABLE playlist (
    playlist_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    visibility TEXT NOT NULL CHECK (visibility IN ('public','private','unlisted')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_playlist_user FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON DELETE CASCADE
);

-- Core media entity (movie / series / episode)
CREATE TABLE media (
    media_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    media_type TEXT NOT NULL CHECK (media_type IN ('movie','series','episode')),
    title TEXT NOT NULL,
    description TEXT,
    release_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Genres
CREATE TABLE genre (
    genre_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Media genres (many-to-many)
CREATE TABLE media_genre (
    media_id BIGINT NOT NULL,
    genre_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (media_id, genre_id),
    CONSTRAINT fk_mg_media FOREIGN KEY (media_id) REFERENCES media (media_id) ON DELETE CASCADE,
    CONSTRAINT fk_mg_genre FOREIGN KEY (genre_id) REFERENCES genre (genre_id) ON DELETE CASCADE
);

CREATE TABLE playlist_media (
    playlist_id BIGINT NOT NULL,
    media_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (playlist_id, media_id),
    CONSTRAINT fk_pm_playlist FOREIGN KEY (playlist_id) REFERENCES playlist (playlist_id) ON DELETE CASCADE,
    CONSTRAINT fk_pm_media FOREIGN KEY (media_id) REFERENCES media (media_id) ON DELETE CASCADE
);

-- Episode-specific details (episode.media_id references media.media_id)
CREATE TABLE media_episode (
    episode_id BIGINT PRIMARY KEY,
    series_id BIGINT NOT NULL,
    season_number INTEGER NOT NULL,
    episode_number INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_episode_media FOREIGN KEY (episode_id) REFERENCES media (media_id) ON DELETE CASCADE,
    CONSTRAINT fk_episode_series FOREIGN KEY (series_id) REFERENCES media (media_id) ON DELETE CASCADE,
    -- Ensure episode_id refers to a media row with media_type='episode'
    CONSTRAINT chk_episode_id_is_episode CHECK (
        EXISTS (
            SELECT 1 FROM media m WHERE m.media_id = episode_id AND m.media_type = 'episode'
        )
    ),
    -- Ensure series_id refers to a media row with media_type='series'
    CONSTRAINT chk_series_id_is_series CHECK (
        EXISTS (
            SELECT 1 FROM media m WHERE m.media_id = series_id AND m.media_type = 'series'
        )
    )
);

-- Media <-> Images
CREATE TABLE asset (
    asset_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    s3_key TEXT NOT NULL UNIQUE,
    mime_type TEXT NOT NULL,
    size_mb NUMERIC(12,3) NOT NULL CHECK (size_mb >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE asset_image (
    asset_image_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    asset_id BIGINT NOT NULL UNIQUE,
    resolution_width INTEGER,
    resolution_height INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_asset_image_asset FOREIGN KEY (asset_id) REFERENCES asset (asset_id) ON DELETE CASCADE
);

CREATE TABLE media_image (
    media_image_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    media_id BIGINT NOT NULL,
    asset_image_id BIGINT NOT NULL,
    image_type TEXT NOT NULL CHECK (image_type IN ('poster','preview','thumbnail')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_media_image_media FOREIGN KEY (media_id) REFERENCES media (media_id) ON DELETE CASCADE,
    CONSTRAINT fk_media_image_asset FOREIGN KEY (asset_image_id) REFERENCES asset_image (asset_image_id) ON DELETE CASCADE
);

-- Video assets
CREATE TABLE asset_video (
    asset_video_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    asset_id BIGINT NOT NULL UNIQUE,
    quality TEXT NOT NULL,
    resolution_width INTEGER,
    resolution_height INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_asset_video_asset FOREIGN KEY (asset_id) REFERENCES asset (asset_id) ON DELETE CASCADE
);

CREATE TABLE media_video (
    media_video_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    media_id BIGINT NOT NULL,
    video_type TEXT NOT NULL CHECK (video_type IN ('main_video','trailer','teaser')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_media_video_media FOREIGN KEY (media_id) REFERENCES media (media_id) ON DELETE CASCADE
);

CREATE TABLE media_video_asset (
    media_video_id BIGINT NOT NULL,
    asset_video_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (media_video_id, asset_video_id),
    CONSTRAINT fk_mva_media_video FOREIGN KEY (media_video_id) REFERENCES media_video (media_video_id) ON DELETE CASCADE,
    CONSTRAINT fk_mva_asset_video FOREIGN KEY (asset_video_id) REFERENCES asset_video (asset_video_id) ON DELETE CASCADE
);

-- Audio assets
CREATE TABLE asset_audio (
    asset_audio_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    asset_id BIGINT NOT NULL UNIQUE,
    language TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_asset_audio_asset FOREIGN KEY (asset_id) REFERENCES asset (asset_id) ON DELETE CASCADE
);

CREATE TABLE media_audio (
    media_video_id BIGINT NOT NULL,
    asset_audio_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (media_video_id, asset_audio_id),
    CONSTRAINT fk_media_audio_media_video FOREIGN KEY (media_video_id) REFERENCES media_video (media_video_id) ON DELETE CASCADE,
    CONSTRAINT fk_media_audio_asset FOREIGN KEY (asset_audio_id) REFERENCES asset_audio (asset_audio_id) ON DELETE CASCADE
);

-- Subtitle assets
CREATE TABLE asset_subtitle (
    asset_subtitle_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    asset_id BIGINT NOT NULL UNIQUE,
    language TEXT NOT NULL,
    format TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_asset_subtitle_asset FOREIGN KEY (asset_id) REFERENCES asset (asset_id) ON DELETE CASCADE
);

CREATE TABLE media_subtitle (
    media_video_id BIGINT NOT NULL,
    asset_subtitle_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (media_video_id, asset_subtitle_id),
    CONSTRAINT fk_media_subtitle_media_video FOREIGN KEY (media_video_id) REFERENCES media_video (media_video_id) ON DELETE CASCADE,
    CONSTRAINT fk_media_subtitle_asset FOREIGN KEY (asset_subtitle_id) REFERENCES asset_subtitle (asset_subtitle_id) ON DELETE CASCADE
);

-- Actors and roles
CREATE TABLE actor (
    actor_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT NOT NULL,
    birth_date DATE,
    bio TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE actor_image (
    actor_image_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    actor_id BIGINT NOT NULL,
    asset_image_id BIGINT NOT NULL,
    image_type TEXT NOT NULL CHECK (image_type IN ('profile','other')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_actor_image_actor FOREIGN KEY (actor_id) REFERENCES actor (actor_id) ON DELETE CASCADE,
    CONSTRAINT fk_actor_image_asset FOREIGN KEY (asset_image_id) REFERENCES asset_image (asset_image_id) ON DELETE CASCADE
);

CREATE TABLE actor_role (
    actor_role_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    actor_id BIGINT NOT NULL,
    media_id BIGINT NOT NULL,
    role_name TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_actor_role_actor FOREIGN KEY (actor_id) REFERENCES actor (actor_id) ON DELETE CASCADE,
    CONSTRAINT fk_actor_role_media FOREIGN KEY (media_id) REFERENCES media (media_id) ON DELETE CASCADE
);

-- User-playlist link (collaboration/role)
CREATE TABLE user_playlist (
    user_id BIGINT NOT NULL,
    playlist_id BIGINT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('owner','collaborator','viewer')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, playlist_id),
    CONSTRAINT fk_up_user FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON DELETE CASCADE,
    CONSTRAINT fk_up_playlist FOREIGN KEY (playlist_id) REFERENCES playlist (playlist_id) ON DELETE CASCADE
);

-- Watch history
CREATE TABLE user_watch_history (
    watch_history_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT NOT NULL,
    media_id BIGINT NOT NULL,
    watched_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    progress_seconds INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_wh_user FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON DELETE CASCADE,
    CONSTRAINT fk_wh_media FOREIGN KEY (media_id) REFERENCES media (media_id) ON DELETE CASCADE
);

-- Likes
CREATE TABLE user_like_media (
    user_id BIGINT NOT NULL,
    media_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, media_id),
    CONSTRAINT fk_ulm_user FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON DELETE CASCADE,
    CONSTRAINT fk_ulm_media FOREIGN KEY (media_id) REFERENCES media (media_id) ON DELETE CASCADE
);

CREATE TABLE user_like_actor (
    user_id BIGINT NOT NULL,
    actor_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, actor_id),
    CONSTRAINT fk_ula_user FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON DELETE CASCADE,
    CONSTRAINT fk_ula_actor FOREIGN KEY (actor_id) REFERENCES actor (actor_id) ON DELETE CASCADE
);

CREATE TABLE user_like_playlist (
    user_id BIGINT NOT NULL,
    playlist_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, playlist_id),
    CONSTRAINT fk_ulp_user FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON DELETE CASCADE,
    CONSTRAINT fk_ulp_playlist FOREIGN KEY (playlist_id) REFERENCES playlist (playlist_id) ON DELETE CASCADE
);

-- Comments
CREATE TABLE user_comment_media (
    user_comment_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT NOT NULL,
    media_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ucm_user FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON DELETE CASCADE,
    CONSTRAINT fk_ucm_media FOREIGN KEY (media_id) REFERENCES media (media_id) ON DELETE CASCADE
);

CREATE TABLE user_comment_actor (
    user_comment_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT NOT NULL,
    actor_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_uca_user FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON DELETE CASCADE,
    CONSTRAINT fk_uca_actor FOREIGN KEY (actor_id) REFERENCES actor (actor_id) ON DELETE CASCADE
);

-- Ratings
CREATE TABLE user_rating_media (
    user_id BIGINT NOT NULL,
    media_id BIGINT NOT NULL,
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, media_id),
    CONSTRAINT fk_urm_user FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON DELETE CASCADE,
    CONSTRAINT fk_urm_media FOREIGN KEY (media_id) REFERENCES media (media_id) ON DELETE CASCADE
);

-- Saved (bookmarks)
CREATE TABLE saved_media (
    saved_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT NOT NULL,
    media_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_saved_user FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON DELETE CASCADE,
    CONSTRAINT fk_saved_media FOREIGN KEY (media_id) REFERENCES media (media_id) ON DELETE CASCADE,
    CONSTRAINT unique_saved UNIQUE (user_id, media_id)
);
