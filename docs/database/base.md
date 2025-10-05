# ER диаграмма базы данных PopFilms

```mermaid
erDiagram
    USER {
        bigint user_id PK
        text username
        text email
        text password_hash
        timestamptz created_at
        timestamptz updated_at
    }

    GENRE {
        smallint genre_id PK
        text genre
        timestamptz created_at
        timestamptz updated_at
    }

%% Keep MEDIA_ASSETS to change quality of content
    VIDEO {
        bigint video_id PK
        bigint series_id FK "NULL for films"
        int season_number "NULL for films"
        int episode_number "NULL for films"
        text title
        text description
        text type          "movie / episode / series"
        date release_date
        text thumbnail_s3_key
        timestamptz created_at
        timestamptz updated_at
    }

    VIDEO_GENRE {
        bigint id PK
        bigint video_id FK
        smallint genre_id FK
        timestamptz created_at
        timestamptz updated_at
    }    

    SERIES {
        bigint series_id PK
        text title
        text description
        date release_date
        text thumbnail_s3_key
        timestamptz created_at
        timestamptz updated_at
    }

    SERIES_GENRE {
        bigint id
        bigint series_id
        smallint genre_id
    }

    MEDIA_ASSET {
        bigint asset_id PK
        bigint video_id FK
        text asset_type   "video / trailer"
        text s3_key
        text quality    "1080p / 4K / etc"
        text format "mp4 / webm / vtt"
        decimal size_mb
        timestamptz created_at
        timestamptz updated_at
    }

    ACTOR {
        int actor_id PK
        text name
        date birth_date
        text bio
        timestamptz created_at
        timestamptz updated_at
    }

    VIDEO_ACTOR {
        bigint id PK
        bigint video_id FK
        int actor_id FK
        timestamptz created_at
        timestamptz updated_at
    }

    LIKE {
        bigint like_id PK
        bigint user_id FK
        bigint video_id FK
        timestamptz created_at
        timestamptz updated_at
    }

    COMMENT {
        int comment_id PK
        int user_id FK
        int video_id FK
        text content
        timestamptz created_at
        timestamptz updated_at
    }

    RATING {
        bigint rating_id PK
        bigint user_id FK
        bigint video_id FK
        smallint rating
        timestamptz created_at
        timestamptz updated_at
    }

    SAVED_VIDEO {
        bigint saved_id PK
        bigint user_id FK
        bigint video_id FK
        timestamptz created_at
        timestamptz updated_at
    }

     USER ||--o{ LIKE : "likes"
    VIDEO ||--o{ LIKE : "liked by"

    USER ||--o{ COMMENT : "writes"
    VIDEO ||--o{ COMMENT : "has"

    USER ||--o{ RATING : "rates"
    VIDEO ||--o{ RATING : "has"

    USER ||--o{ SAVED_VIDEO : "saves"
    VIDEO ||--|{ SAVED_VIDEO : "saved by"

    GENRE ||--|{ VIDEO_GENRE : "classifies"
    VIDEO ||--|{ VIDEO_GENRE : "has genre"

    SERIES ||--o{ VIDEO : "has episodes"

    VIDEO ||--|{ MEDIA_ASSET : "has"
    VIDEO ||--o{ VIDEO_ACTOR : "casts"
    ACTOR ||--o{ VIDEO_ACTOR : "acts in"

    GENRE ||--|{ SERIES_GENRE: "classifies"
    SERIES ||--|{ SERIES_GENRE: "has genre"
    
```