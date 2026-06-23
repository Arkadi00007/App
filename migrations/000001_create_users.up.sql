CREATE TABLE users (
                       id            BIGSERIAL PRIMARY KEY,
                       email         VARCHAR(255) UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       name          VARCHAR(255),
                       is_verified BOOLEAN DEFAULT FALSE,
                       created_at    TIMESTAMP DEFAULT NOW()
);