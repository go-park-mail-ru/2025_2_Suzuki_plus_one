-- Test data for local development (minimal, matching migrations)
-- 1) Avatar asset and image
INSERT INTO asset (s3_key, file_size_mb, mime_type)
VALUES (
    '/avatars/dima.jpeg',
    204800::numeric / 1024 / 1024,
    'image/jpeg'
  );
INSERT INTO asset_image (asset_id, resolution_width, resolution_height)
VALUES (
    (
      SELECT asset_id
      FROM asset
      WHERE s3_key = '/avatars/dima.jpeg'
    ),
    256,
    256
  );
-- 2) Test user referencing the avatar
INSERT INTO "user" (user_id, username, asset_image_id, password_hash, email, date_of_birth, phone_number)
VALUES (
        11,
    'testuser',
    (
      SELECT asset_image_id
      FROM asset_image
      WHERE asset_id = (
          SELECT asset_id
          FROM asset
          WHERE s3_key = '/avatars/dima.jpeg'
        )
    ),
    '$2y$10$U5D2NWz2Q9TDsl5YKfHQ5O5qlgCH4SAAva7406ZyDQ/sj53Aoif.G', -- plaintext: 'Password123!'
    'test@example.com',
    '1990-01-01',
    '+1234567890'
  );