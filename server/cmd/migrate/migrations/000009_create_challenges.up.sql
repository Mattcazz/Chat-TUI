CREATE TABLE auth_challenges (
    user_id     BIGINT NOT NULL PRIMARY KEY, -- One active challenge per user
    nonce       VARCHAR(32) NOT NULL,            -- The random string
    expires_at  TIMESTAMP NOT NULL,       -- Security timeout (e.g., 30 seconds)
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);