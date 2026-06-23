CREATE TABLE verification_codes (
                                    id         BIGSERIAL PRIMARY KEY,
                                    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                    code       VARCHAR(6) NOT NULL,
                                    type       VARCHAR(20) NOT NULL CHECK (type IN ('email_verify', 'reset_password')),
                                    expires_at TIMESTAMP NOT NULL,
                                    used_at    TIMESTAMP,         -- NULL пока не использован
                                    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_verification_codes_lookup
    ON verification_codes(user_id, type);