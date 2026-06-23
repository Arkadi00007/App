-- 000009_create_attempts.up.sql
CREATE TABLE test_attempts (
                               id                 BIGSERIAL PRIMARY KEY,
                               user_id            BIGINT NOT NULL REFERENCES users(id),
                               subject_id         BIGINT NOT NULL REFERENCES subjects(id),
                               test_id            BIGINT REFERENCES tests(id),
                               exam_id            BIGINT REFERENCES exams(id),
                               mode               VARCHAR(20) NOT NULL CHECK (mode IN ('test', 'exam', 'practice')),
                               status             VARCHAR(20) DEFAULT 'in_progress'
                                   CHECK (status IN ('in_progress', 'completed', 'abandoned')),
                               started_at         TIMESTAMP DEFAULT NOW(),
                               finished_at        TIMESTAMP,
                               time_limit_minutes INTEGER,
                               score              INTEGER DEFAULT 0,
                               max_score          INTEGER DEFAULT 0,
                               percentage         NUMERIC(5,2)
);

CREATE TABLE user_answers (
                              id          BIGSERIAL PRIMARY KEY,
                              attempt_id  BIGINT NOT NULL REFERENCES test_attempts(id) ON DELETE CASCADE,
                              question_id BIGINT NOT NULL REFERENCES questions(id),
                              answer_id   BIGINT REFERENCES answers(id),
                              text_answer TEXT,
                              answer_ids BIGINT[],
                              is_correct  BOOLEAN,
                              answered_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_attempts_user_id ON test_attempts(user_id);
CREATE INDEX idx_attempts_user_status ON test_attempts(user_id, status);
CREATE INDEX idx_attempts_subject_id ON test_attempts(subject_id);
CREATE INDEX idx_user_answers_attempt_id ON user_answers(attempt_id);
CREATE INDEX idx_user_answers_question_id ON user_answers(question_id);