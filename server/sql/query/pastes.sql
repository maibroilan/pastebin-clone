-- name: CreatePaste :one
INSERT INTO pastes (slug, content, password_hash, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPaste :one
SELECT * FROM pastes
WHERE slug = $1;

-- name: Ping :one
SELECT 1;