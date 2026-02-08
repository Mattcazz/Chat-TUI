CREATE TYPE contact_status AS ENUM ('pending', 'accepted', 'blocked');

CREATE TABLE contacts (
    id BIGSERIAL PRIMARY KEY,
    from_user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    to_user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    nickname VARCHAR(100),
    status contact_status NOT NULL DEFAULT 'pending',
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (from_user_id, to_user_id)
);

CREATE INDEX idx_contacts_user ON contacts(from_user_id);    