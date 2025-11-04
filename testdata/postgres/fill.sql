-- Test data for local development (minimal, matching migrations)
-- 1) Avatar asset and image
INSERT INTO asset (s3_key, size_mb, mime_type)
VALUES (
    '/avatars/Chris.png',
    204800::numeric / 1024 / 1024,
    'image/png'
  );
INSERT INTO asset_image (asset_id, resolution_width, resolution_height)
VALUES (
    (
      SELECT asset_id
      FROM asset
      WHERE s3_key = '/avatars/Chris.png'
    ),
    256,
    256
  );
-- 2) Test user referencing the avatar
INSERT INTO "user" (username, asset_image_id, password_hash, email)
VALUES (
    'testuser',
    (
      SELECT asset_image_id
      FROM asset_image
      WHERE asset_id = (
          SELECT asset_id
          FROM asset
          WHERE s3_key = '/avatars/Chris.png'
        )
    ),
    '$2y$10$U5D2NWz2Q9TDsl5YKfHQ5O5qlgCH4SAAva7406ZyDQ/sj53Aoif.G', -- plaintext: 'Password123!'
    'test@example.com'
  );
-- 3) Media: Inception
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
    'Inception',
    'A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a C.E.O., but his tragic past may doom the project and his team to disaster.',
    '2010-07-16',
    8.8,
    148,
    13,
    'USA',
    'Dom Cobb is a skilled thief, the absolute best in the dangerous art of extraction: stealing valuable secrets from deep within the subconscious during the dream state when the mind is at its most vulnerable. Cobb''s rare ability has made him a coveted player in this treacherous new world of corporate espionage, but it has also made him an international fugitive and cost him everything he has ever loved. Now Cobb is being offered a chance at redemption. One last job could give him his life back but only if he can accomplish the impossible - inception. Instead of the perfect heist, Cobb and his team of specialists have to pull off the reverse: their task is not to steal an idea but to plant one. If they succeed, it could be the perfect crime. But no amount of careful planning or expertise can prepare the team for the dangerous enemy that seems to predict their every move. An enemy that only Cobb could have seen coming.'
  );
-- 4) Poster asset and image
INSERT INTO asset (s3_key, size_mb, mime_type)
VALUES (
    '/posters/InceptionPoster.png',
    512000::numeric / 1024 / 1024,
    'image/png'
  );
INSERT INTO asset_image (asset_id, resolution_width, resolution_height)
VALUES (
    (
      SELECT asset_id
      FROM asset
      WHERE s3_key = '/posters/InceptionPoster.png'
    ),
    1200,
    1800
  );
INSERT INTO media_image (media_id, asset_image_id, image_type)
VALUES (
    (
      SELECT media_id
      FROM media
      WHERE title = 'Inception'
    ),
    (
      SELECT asset_image_id
      FROM asset_image
      WHERE asset_id = (
          SELECT asset_id
          FROM asset
          WHERE s3_key = '/posters/InceptionPoster.png'
        )
    ),
    'poster'
  );
-- 5) Main video asset, asset_video and media_video
INSERT INTO asset (s3_key, size_mb, mime_type)
VALUES (
    '/medias/InceptionMovie.webm',
    104857600::numeric / 1024 / 1024,
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
      WHERE s3_key = '/medias/InceptionMovie.webm'
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
      WHERE title = 'Inception'
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
          WHERE title = 'Inception'
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
          WHERE s3_key = '/medias/InceptionMovie.webm'
        )
      ORDER BY asset_video_id DESC
      LIMIT 1
    )
  );
-- 6) Trailer asset, asset_video and media_video
INSERT INTO asset (s3_key, size_mb, mime_type)
VALUES (
    '/trailers/InceptionTrailer.webm',
    20971520::numeric / 1024 / 1024,
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
      WHERE s3_key = '/trailers/InceptionTrailer.webm'
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
      WHERE title = 'Inception'
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
          WHERE title = 'Inception'
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
          WHERE s3_key = '/trailers/InceptionTrailer.webm'
        )
      ORDER BY asset_video_id DESC
      LIMIT 1
    )
  );
-- end of minimal test data