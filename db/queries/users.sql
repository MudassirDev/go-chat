-- name: CreateUser :one
INSERT INTO users (
  id, username, created_at, updated_at
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;
