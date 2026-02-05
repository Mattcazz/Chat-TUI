CREATE TYPE upload_status AS ENUM ('pending', 'uploading', 'assembling', 'completed');

CREATE TABLE upload_sessions (
    id BIGSERIAL PRIMARY KEY,
    file_id BIGINT REFERENCES files(id) ON DELETE CASCADE,
    total_size BIGINT NOT NULL,
    chunk_size INT NOT NULL,
    total_chunks INT NOT NULL,
    status upload_status DEFAULT 'pending',
    expires_at TIMESTAMP NOT NULL
);

