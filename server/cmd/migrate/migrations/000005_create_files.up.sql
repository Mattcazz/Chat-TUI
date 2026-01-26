CREATE TYPE file_status AS ENUM ('uploading', 'ready', 'expired');

CREATE TABLE IF NOT EXISTS files (
    id BIGSERIAL PRIMARY KEY DEFAULT,
    conversation_id BIGINT REFERENCES conversations(id) ON DELETE CASCADE,
    uploader_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL, 
    checksum VARCHAR(64), -- SHA256
    storage_path TEXT NOT NULL, 
    status file_status DEFAULT 'uploading',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ
);