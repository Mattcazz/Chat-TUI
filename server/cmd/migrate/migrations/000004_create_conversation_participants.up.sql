CREATE TYPE participant_role AS ENUM ('admin', 'member');

CREATE TABLE conversation_participants (
    conversation_id BIGINT REFERENCES conversations(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    role participant_role DEFAULT 'member',
    unread_count INT DEFAULT 0,
    PRIMARY KEY (conversation_id, user_id)
);

CREATE INDEX idx_participants_user ON conversation_participants(user_id);