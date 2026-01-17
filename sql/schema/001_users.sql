-- +goose Up
CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  email text NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;

