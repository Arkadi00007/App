-- 000008_create_exams_collections.up.sql
CREATE TABLE exams (
                       id               BIGSERIAL PRIMARY KEY,
                       subject_id       BIGINT NOT NULL REFERENCES subjects(id),
                       title            VARCHAR(255) NOT NULL,
                       description      TEXT,
                       duration_minutes INTEGER,
                       year             INTEGER,
                       created_at       TIMESTAMP DEFAULT NOW(),
                       is_premium BOOLEAN DEFAULT FALSE
);

CREATE TABLE exam_questions (
                                id             BIGSERIAL PRIMARY KEY,
                                exam_id        BIGINT NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
                                question_id    BIGINT NOT NULL REFERENCES questions(id),
                                question_order INTEGER NOT NULL,
                                UNIQUE (exam_id, question_id)
);

CREATE TABLE collections (
                             id          BIGSERIAL PRIMARY KEY,
                             subject_id  BIGINT NOT NULL REFERENCES subjects(id),
                             title       VARCHAR(255) NOT NULL,
                             description TEXT
);

CREATE TABLE collection_questions (
                                      id            BIGSERIAL PRIMARY KEY,
                                      collection_id BIGINT NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
                                      question_id   BIGINT NOT NULL REFERENCES questions(id)
);

CREATE INDEX idx_exams_subject_id ON exams(subject_id);
CREATE INDEX idx_exam_questions_exam_id_order ON exam_questions(exam_id, question_order);
CREATE INDEX idx_collections_subject_id ON collections(subject_id);
CREATE INDEX idx_collection_questions_collection_id ON collection_questions(collection_id);


-- стоит добавить порядок вопросов как в test_questions
ALTER TABLE collection_questions ADD COLUMN question_order INTEGER NOT NULL DEFAULT 0;

-- и индекс
CREATE INDEX idx_collection_questions_collection_id_order
    ON collection_questions(collection_id, question_order);