-- 000007_create_answers.up.sql
CREATE TABLE answers (
                         id          BIGSERIAL PRIMARY KEY,
                         question_id BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
                         answer_text TEXT NOT NULL,
                         is_correct  BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_answers_question_id ON answers(question_id);