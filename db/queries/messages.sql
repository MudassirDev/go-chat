-- name: CreateMessage :one
INSERT INTO messages (
  sender_id, recipient_id, time, content, message_type, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetChatMessages :many
SELECT * FROM messages WHERE (recipient_id = $1 AND sender_id = $2) OR (sender_id = $1 AND recipient_id = $2) ORDER BY time ASC;

-- name: GetMessageWithFileName :one
SELECT * FROM messages WHERE content = $1 AND (recipient_id = $2 OR sender_id = $2);
