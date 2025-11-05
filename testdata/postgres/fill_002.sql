-- Additional test media for local development (matched to migrations)

-- Media: The Matrix
INSERT INTO media (
    media_type,
    title,
    description,
    release_date,
    rating,
    duration_minutes,
    age_rating,
    country,
    plot_summary
  )
VALUES (
    'movie',
    'The Matrix',
    'A computer hacker learns about the true nature of his reality and his role in the war against its controllers.',
    '1999-03-31',
    8.7,
    136,
    16,
    'USA',
    'Thomas Anderson, a computer programmer by day and hacker by night, is drawn into a rebellion against the machines which have trapped humanity in a simulated reality.'
  );

-- Poster asset and image for The Matrix
INSERT INTO asset (s3_key, size_mb, mime_type)
VALUES (
    '/posters/MatrixPoster.png',
    420000::numeric / 1024 / 1024,
    'image/png'
  );

INSERT INTO asset_image (asset_id, resolution_width, resolution_height)
VALUES (
    (
      SELECT asset_id
      FROM asset
      WHERE s3_key = '/posters/MatrixPoster.png'
    ),
    1000,
    1500
  );

INSERT INTO media_image (media_id, asset_image_id, image_type)
VALUES (
    (
      SELECT media_id
      FROM media
      WHERE title = 'The Matrix'
    ),
    (
      SELECT asset_image_id
      FROM asset_image
      WHERE asset_id = (
          SELECT asset_id
          FROM asset
          WHERE s3_key = '/posters/MatrixPoster.png'
        )
    ),
    'poster'
  );

-- Main video asset, asset_video and media_video for The Matrix
INSERT INTO asset (s3_key, size_mb, mime_type)
VALUES (
    '/medias/MatrixMovie.webm',
    90000000::numeric / 1024 / 1024,
    'video/webm'
  );

INSERT INTO asset_video (
    asset_id,
    quality,
    resolution_width,
    resolution_height
  )
VALUES (
    (
      SELECT asset_id
      FROM asset
      WHERE s3_key = '/medias/MatrixMovie.webm'
    ),
    '1080p',
    1920,
    1080
  );

INSERT INTO media_video (media_id, video_type)
VALUES (
    (
      SELECT media_id
      FROM media
      WHERE title = 'The Matrix'
    ),
    'main_video'
  );

INSERT INTO media_video_asset (media_video_id, asset_video_id)
VALUES (
    (
      SELECT media_video_id
      FROM media_video
      WHERE media_id = (
          SELECT media_id
          FROM media
          WHERE title = 'The Matrix'
        )
        AND video_type = 'main_video'
      ORDER BY media_video_id DESC
      LIMIT 1
    ), (
      SELECT asset_video_id
      FROM asset_video
      WHERE asset_id = (
          SELECT asset_id
          FROM asset
          WHERE s3_key = '/medias/MatrixMovie.webm'
        )
      ORDER BY asset_video_id DESC
      LIMIT 1
    )
  );

-- Trailer asset, asset_video and media_video for The Matrix
INSERT INTO asset (s3_key, size_mb, mime_type)
VALUES (
    '/trailers/MatrixTrailer.webm',
    25000000::numeric / 1024 / 1024,
    'video/webm'
  );

INSERT INTO asset_video (
    asset_id,
    quality,
    resolution_width,
    resolution_height
  )
VALUES (
    (
      SELECT asset_id
      FROM asset
      WHERE s3_key = '/trailers/MatrixTrailer.webm'
    ),
    '720p',
    1280,
    720
  );

INSERT INTO media_video (media_id, video_type)
VALUES (
    (
      SELECT media_id
      FROM media
      WHERE title = 'The Matrix'
    ),
    'trailer'
  );

INSERT INTO media_video_asset (media_video_id, asset_video_id)
VALUES (
    (
      SELECT media_video_id
      FROM media_video
      WHERE media_id = (
          SELECT media_id
          FROM media
          WHERE title = 'The Matrix'
        )
        AND video_type = 'trailer'
      ORDER BY media_video_id DESC
      LIMIT 1
    ), (
      SELECT asset_video_id
      FROM asset_video
      WHERE asset_id = (
          SELECT asset_id
          FROM asset
          WHERE s3_key = '/trailers/MatrixTrailer.webm'
        )
      ORDER BY asset_video_id DESC
      LIMIT 1
    )
  );

-- end of additional test data
