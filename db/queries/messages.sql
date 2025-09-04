-- name: CreateMessage :one
INSERT INTO messages (
  sender_id, recipient_id, time, message, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetChatMessages :many
SELECT * FROM messages WHERE (recipient_id = ? AND sender_id = ?) OR (sender_id = ? AND recipient_id = ?) ORDER BY time ASC;
