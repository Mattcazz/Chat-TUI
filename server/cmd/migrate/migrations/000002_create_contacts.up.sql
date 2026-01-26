CREATE TYPE contact_status AS ENUM ('request', 'accepted', 'blocked');

CREATE TABLE IF NOT EXISTS contacts (
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    contact_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    nickname VARCHAR(100),
    status contact_status NOT NULL DEFAULT 'request',
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, contact_id) 
);

CREATE INDEX idx_contacts_user ON contacts(user_id);