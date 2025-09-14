-- name: CreateUser :one
INSERT INTO users (
  username, password, created_at, updated_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, username, created_at, updated_at;

-- name: GetUserWithUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserWithID :one
SELECT id, username, created_at, updated_at FROM users WHERE id = $1;

-- name: GetAllUsersExceptCurrent :many
SELECT id, username, created_at, updated_at FROM users WHERE id != $1;
