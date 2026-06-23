CREATE TABLE sections (
                          id          BIGSERIAL PRIMARY KEY,
                          subject_id  BIGINT NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
                          name        VARCHAR(255) NOT NULL,
                          description TEXT
);

CREATE INDEX idx_sections_subject_id ON sections(subject_id)