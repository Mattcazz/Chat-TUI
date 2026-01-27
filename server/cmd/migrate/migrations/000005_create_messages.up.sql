CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    content TEXT NOT NULL, -- The actual message
    message_type VARCHAR(20) DEFAULT 'text', 
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_messages_conversation_time ON messages(conversation_id, created_at DESC);