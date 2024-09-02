CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL,
    access_token TEXT NOT NULL,
    id_token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
