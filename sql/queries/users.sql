-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
  VALUES (gen_random_uuid (), now(), now(), $1, $2)
RETURNING
  *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT
  *
FROM
  users
WHERE
  email = $1;

-- name: UpdateUserCredentials :one
UPDATE
  users
SET
  hashed_password = $2,
  email = $3,
  updated_at = now()
WHERE
  id = $1
RETURNING
  *;

