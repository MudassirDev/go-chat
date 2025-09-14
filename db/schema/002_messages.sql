-- +goose Up
CREATE TABLE messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  recipient_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  time TIMESTAMP NOT NULL,
  content TEXT NOT NULL,
  message_type TEXT NOT NULL CHECK(message_type = 'TEXT' OR message_type = 'AUDIO'),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE messages;
