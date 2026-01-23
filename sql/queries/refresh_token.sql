-- name: SaveRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at)
  VALUES ($1, $2, $3)
RETURNING
  *;

-- name: GetUserIdFromRefreshToken :one
SELECT
  user_id
FROM
  refresh_tokens
WHERE
  token = $1
  AND expires_at > now()
  AND revoked_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE
  refresh_tokens
SET
  revoked_at = now(),
  updated_at = now()
WHERE
  token = $1
  AND revoked_at IS NULL;

-- name: GetUserFromRefreshToken :one
SELECT
  u.id,
  u.email,
  u.created_at,
  u.updated_at,
  u.hashed_password
FROM
  users u
  INNER JOIN refresh_tokens r ON u.id = r.user_id
WHERE
  r.token = $1
  AND expires_at > now()
  AND revoked_at IS NULL;

