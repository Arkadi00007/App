-- 000006_create_test_questions.up.sql
CREATE TABLE test_questions (
                                id             BIGSERIAL PRIMARY KEY,
                                test_id        BIGINT NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
                                question_id    BIGINT NOT NULL REFERENCES questions(id),
                                question_order INTEGER NOT NULL,
                                UNIQUE (test_id, question_id)
);

-- WHERE test_id = $1 ORDER BY question_order
-- UNIQUE(test_id, question_id) уже создаёт составной индекс,
-- но он работает как (test_id, question_id) — для ORDER BY добавим отдельный
CREATE INDEX idx_test_questions_test_id_order ON test_questions(test_id, question_order);