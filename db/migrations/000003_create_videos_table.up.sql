BEGIN;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'video_format') THEN
        CREATE TYPE video_format AS ENUM ('dvd', 'bluray', 'digital', 'vhs');
    END IF;
END $$;


CREATE TABLE IF NOT EXISTS videos(
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR NOT NULL,
  overview VARCHAR,
  genre_ids INT[],
  format video_format,
  production_company VARCHAR,
  rent_price DECIMAL(12, 2),
  cover_path VARCHAR,
  total_stock INTEGER NOT NULL DEFAULT 0,
  available_stock INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);

COMMIT;