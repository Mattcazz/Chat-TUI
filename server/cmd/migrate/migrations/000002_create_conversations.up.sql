CREATE TABLE conversations (
    id BIGSERIAL PRIMARY KEY,
    last_message_at TIMESTAMPTZ,
    last_message_preview TEXT, -- For the inbox list UI
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);