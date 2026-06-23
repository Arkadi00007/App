CREATE TABLE tests (
                       id                BIGSERIAL PRIMARY KEY,
                       section_id        BIGINT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
                       title             VARCHAR(255) NOT NULL,
                       description       TEXT,
                       created_at        TIMESTAMP DEFAULT NOW(),
                       is_premium        BOOLEAN DEFAULT FALSE,
                       show_answer_mode  VARCHAR(20) NOT NULL DEFAULT 'immediate'
                           CHECK (show_answer_mode IN ('immediate', 'end_only'))
);

-- WHERE section_id = $1
CREATE INDEX idx_tests_section_id ON tests(section_id);