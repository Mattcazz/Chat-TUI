CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(25) UNIQUE NOT NULL,
    public_key TEXT NOT NULL, -- SSH public key
    fingerprint TEXT,         
    last_seen TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);