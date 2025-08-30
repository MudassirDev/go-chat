-- name: CreateUser :one
INSERT INTO users (
  username, created_at, updated_at
) VALUES (
  ?, ?, ?
)
RETURNING *;
