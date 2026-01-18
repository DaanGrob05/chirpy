-- +goose Up
CREATE TABLE chirps (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),
  body text NOT NULL,
  user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;

