-- +goose Up
CREATE TABLE messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  sender_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  recipient_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  time TIMESTAMP NOT NULL,
  content TEXT NOT NULL,
  message_type TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE messages;
