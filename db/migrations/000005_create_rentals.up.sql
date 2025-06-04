BEGIN;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'rental_status') THEN
    -- pending is status for rented videos with payment status of pending
    -- failed is status for rented videos with payment status of expired
    -- rented is status for rented videos with payment status of success
        CREATE TYPE rental_status AS ENUM ('rented', 'returned', 'pending', 'failed');
    END IF;
END $$;


CREATE TABLE  IF NOT EXISTS rentals (
  id BIGSERIAL PRIMARY KEY,
  video_id BIGINT NOT NULL,
  rental_payment_id BIGINT,
  latefee_payment_id BIGINT,
  user_id UUID NOT NULL,
  due_date TIMESTAMP NOT NULL,
  return_date TIMESTAMP,
  status rental_status NOT NULL DEFAULT 'pending',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP,
  FOREIGN KEY (video_id) REFERENCES videos(id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (rental_payment_id) REFERENCES payments(id),
  FOREIGN KEY (latefee_payment_id) REFERENCES payments(id)
);

COMMIT;