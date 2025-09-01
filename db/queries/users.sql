-- name: CreateUser :one
INSERT INTO users (
  username, password, created_at, updated_at
) VALUES (
  ?, ?, ?, ?
)
RETURNING id, username, created_at, updated_at;

-- name: GetUserWithUsername :one
SELECT * FROM users WHERE username = ?;

-- name: GetUserWithID :one
SELECT id, username, created_at, updated_at FROM users WHERE id = ?;

-- name: GetAllUsersExceptCurrent :many
SELECT id, username, created_at, updated_at FROM users WHERE id != ?;
