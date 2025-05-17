CREATE TABLE  IF NOT EXISTS rentals (
  id BIGSERIAL PRIMARY KEY,
  video_id BIGINT NOT NULL,
  payment_id BIGINT NOT NULL,
  due_date TIMESTAMP NOT NULL,
  return_date TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP,
  FOREIGN KEY (video_id) REFERENCES videos(id),
  FOREIGN KEY (payment_id) REFERENCES payments(id)
);