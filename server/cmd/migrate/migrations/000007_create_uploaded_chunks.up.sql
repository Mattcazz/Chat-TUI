CREATE TABLE IF NOT EXISTS uploaded_chunks (
    id BIGSERIAL PRIMARY KEY DEFAULT,
    upload_session_id BIGINT REFERENCES upload_sessions(id) ON DELETE CASCADE,
    chunk_index INT NOT NULL,
    size INT NOT NULL,
    checksum VARCHAR(64),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    -- Prevent duplicate chunks
    UNIQUE (upload_session_id, chunk_index)
);
