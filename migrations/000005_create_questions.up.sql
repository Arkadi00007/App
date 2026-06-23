CREATE TABLE questions (
                           id            BIGSERIAL PRIMARY KEY,
                           subject_id    BIGINT NOT NULL REFERENCES subjects(id),
                           question_text TEXT NOT NULL,
                           question_type VARCHAR(50) NOT NULL,
                           image_url     TEXT,
                           explanation   TEXT,
                           difficulty    SMALLINT NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
                           points        INTEGER DEFAULT 1,
                           created_at    TIMESTAMP DEFAULT NOW()
);

-- WHERE subject_id = $1 (фильтр вопросов по предмету)
CREATE INDEX idx_questions_subject_id ON questions(subject_id);

-- WHERE difficulty = $1 (фильтр по сложности)
CREATE INDEX idx_questions_difficulty ON questions(difficulty);