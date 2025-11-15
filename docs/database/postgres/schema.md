# ER диаграмма базы данных PopFilms

```mermaid
erDiagram

    %% # Genres
    GENRE {
        smallint genre_id PK
        %% ---
        text name
        text description
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    %% # Playlists

    PLAYLIST {
        bigint playlist_id PK
        bigint user_id FK "creator of the playlist"
        %% ---
        text name
        text description
        text visibility "public / private / unlisted"
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    PLAYLIST_MEDIA {
        bigint playlist_media_id PK
        bigint playlist_id FK
        bigint media_id FK
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    %% # Representations of movies and series

    %% Core MEDIA entity which can be a movie, series, or episode
    MEDIA {
        bigint media_id PK
        %% ---
        text media_type   "movie / series / episode"
        text title
        text description
        date release_date
        float rating
        integer duration_minutes
        integer age_rating
        text country
        text plot_summary
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    MEDIA_GENRE {
        bigint media_genre_id PK
        smallint genre_id FK
        bigint media_id FK
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    %% Specific details for MEDIA of type 'episode'
    MEDIA_EPISODE {
        bigint episode_id PK "refers to MEDIA of type 'episode'"
        bigint series_id FK "refers to MEDIA of type 'series'"
        %% ---
        integer season_number
        integer episode_number
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    %% Connects MEDIA to its various IMAGE assets (posters, previews, etc)
    MEDIA_IMAGE {
        bigint media_image_id PK
        bigint media_id FK
        bigint asset_image_id FK
        %% ---
        text image_type   "poster / preview"
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    %% Connects MEDIA to its various VIDEO assets (main video, trailers, etc)
    MEDIA_VIDEO {
        bigint media_video_id PK
        bigint media_id FK
        bigint asset_video_id FK
        %% ---
        text video_type   "main_video / trailer"
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    %% Connects MEDIA_VIDEO_ASSET to its AUDIOs
    MEDIA_AUDIO {
        bigint media_audio_id PK
        bigint media_video_id FK
        bigint asset_audio_id FK
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    MEDIA_SUBTITLE {
        bigint media_subtitle_id PK
        bigint media_video_id FK
        bigint asset_subtitle_id FK
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }


    %% # S3 representations of different asset types

    %% An asset can be a video file, audio track, subtitle file, or image
    ASSET {
        bigint asset_id PK
        %% ---
        text s3_key
        text mime_type
        real file_size_mb
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    ASSET_VIDEO {
        bigint asset_video_id PK
        bigint asset_id FK
        %% ---
        text quality
        integer resolution_width
        integer resolution_height
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    ASSET_SUBTITLE {
        bigint asset_subtitle_id PK
        bigint asset_id FK
        %% ---
        text language
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    ASSET_AUDIO {
        bigint asset_audio_id PK
        bigint asset_id FK
        %% ---
        text language
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    ASSET_IMAGE {
        bigint asset_image_id PK
        bigint asset_id FK
        %% ---
        integer resolution_width
        integer resolution_height
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    %% # Actors and their association with videos
    ACTOR {
        bigint actor_id PK
        %% ---
        text name
        date birth_date
        text bio
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    ACTOR_IMAGE {
        bigint actor_image_id PK
        bigint actor_id FK
        bigint asset_image_id FK
        %% ---
        text image_type   "profile / other"
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    ACTOR_ROLE {
        bigint actor_role_id PK
        bigint actor_id FK
        bigint media_id FK
        %% ---
        text role_name
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    %% # Relationships involving users
    USER {
        bigint user_id PK
        %% ---
        text username
        bigint asset_image_id FK
        text email
        text password_hash
        date date_of_birth
        text phone_number
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    USER_APPEAL {
        bigint user_appeal_id PK
        bigint user_id FK
        %% ---
        text tag      "bug / feature / complaint / other"
        %% name is first line of a message
        text name
        text status   "open / in_progress / resolved"
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    USER_APPEAL_MESSAGE {
        bigint user_appeal_message_id PK
        bigint user_appeal_id FK
        %% ---
        boolean is_response
        text message_content
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    USER_SESSION {
        bigint user_session_id PK
        bigint user_id FK
        %% ---
        text session_token
        timestamptz expires_at
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    USER_PLAYLIST {
        bigint user_playlist_id PK
        bigint user_id FK
        bigint playlist_id FK
        %% ---
        text role "collaborator / viewer"
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    USER_WATCH_HISTORY {
        bigint watch_history_id PK
        bigint user_id FK
        bigint media_id FK
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    %% ## Likes

    USER_LIKE_MEDIA {
        bigint user_like_media_id PK
        bigint user_id FK
        bigint media_id FK
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    USER_LIKE_ACTOR {
        bigint user_like_actor_id PK
        bigint user_id FK
        bigint actor_id FK
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    USER_LIKE_PLAYLIST {
        bigint user_like_playlist_id PK
        bigint user_id FK
        bigint playlist_id FK
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    %% ## Comments

    USER_COMMENT_MEDIA {
        bigint user_comment_id PK
        bigint user_id FK
        bigint media_id FK
        %% ---
        text content
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    USER_COMMENT_ACTOR {
        bigint user_comment_id PK
        bigint user_id FK
        bigint actor_id FK
        %% ---
        text content
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    %% ## Rating

    USER_RATING_MEDIA {
        bigint user_rating_media_id PK
        bigint user_id FK
        bigint media_id FK
        %% ---
        smallint rating "1 to 5"
        %% ---
        timestamptz created_at
        timestamptz updated_at
    }

    %% # General relationships
    GENRE ||--o{ MEDIA_GENRE : "media_genre has many genres"
        MEDIA_GENRE ||--|| MEDIA : "media_genre links media and genre"

    %% # Playlist relationships
    PLAYLIST ||--o{ PLAYLIST_MEDIA : "playlist has media"
        PLAYLIST_MEDIA ||--|| MEDIA : "playlist_media links playlist and media"

    %% # Core media relationships
    MEDIA ||--o{ MEDIA_EPISODE : "media(type=series) has episodes"
    MEDIA ||--o{ MEDIA_IMAGE : "media has images"
    MEDIA ||--|{ MEDIA_VIDEO : "media has videos"

        MEDIA_VIDEO ||--|{ ASSET_VIDEO : "media_video has video assets"
        MEDIA_VIDEO ||--|{ MEDIA_AUDIO : "media_video has audio"
        MEDIA_VIDEO ||--|{ MEDIA_SUBTITLE : "media_video has subtitles"

        MEDIA_AUDIO ||--o{ ASSET_AUDIO : "media_audio has audio assets"
        MEDIA_SUBTITLE ||--o{ ASSET_SUBTITLE : "media_subtitle has subtitle assets"
        MEDIA_IMAGE ||--o{ ASSET_IMAGE : "media_image has image assets"

    %% ## Asset relationships
    ASSET_VIDEO ||--|| ASSET : "asset_video is an asset"
    ASSET_AUDIO ||--|{ ASSET : "asset_audio is an asset"
    ASSET_SUBTITLE ||--|{ ASSET : "asset_subtitle is an asset"
    ASSET_IMAGE ||--|{ ASSET : "asset_image is an asset"


    %% # Actor relationships
    ACTOR ||--o{ ACTOR_IMAGE : "actor has images"
        ACTOR_IMAGE ||--o{ ASSET_IMAGE : "actor_image links actor and asset_image"
    ACTOR ||--o{ ACTOR_ROLE : "actor has roles in media"
        ACTOR_ROLE ||--|| MEDIA : "actor_role links actor and media"


    %% # User relationships
    USER ||--o{ USER_WATCH_HISTORY : "user has watch history"
        USER_WATCH_HISTORY ||--|| MEDIA : "user_watch_history links user and media"

    USER ||--o{ USER_PLAYLIST : "user has playlists"
        USER_PLAYLIST ||--|| PLAYLIST : "user_playlist links user and playlist"

    USER }o--|| ASSET_IMAGE : "user has profile image"

    USER ||--o{ USER_SESSION : "user has sessions"

    %% ## Likes
    USER ||--o{ USER_LIKE_MEDIA : "user likes media"
        USER_LIKE_MEDIA ||--|| MEDIA : "user_like_media links user and media"
    USER ||--o{ USER_LIKE_ACTOR : "user likes actor"
        USER_LIKE_ACTOR ||--|| ACTOR : "user_like_actor links user and actor"
    USER ||--o{ USER_LIKE_PLAYLIST : "user likes playlist"
        USER_LIKE_PLAYLIST ||--|| PLAYLIST : "user_like_playlist links user and playlist"

    %% ## Comments
    USER ||--o{ USER_COMMENT_MEDIA : "user comments on media"
        USER_COMMENT_MEDIA ||--|| MEDIA : "user_comment_media links user and media"
    USER ||--o{ USER_COMMENT_ACTOR : "user comments on actor"
        USER_COMMENT_ACTOR ||--|| ACTOR : "user_comment_actor links user and actor"

    %% ## Rating
    USER ||--o{ USER_RATING_MEDIA : "user rates media"
        USER_RATING_MEDIA ||--|| MEDIA : "user_rating_media links user and media"
```
