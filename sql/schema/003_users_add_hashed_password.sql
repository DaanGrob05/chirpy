-- +goose Up
ALTER TABLE users
  ADD hashed_password text NOT NULL DEFAULT 'UNSET';

ALTER TABLE users
  ALTER COLUMN hashed_password DROP DEFAULT;

-- +goose Down
ALTER TABLE users
  DROP COLUMN hashed_password;

