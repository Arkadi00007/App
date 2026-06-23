CREATE TABLE users (
                       id            BIGSERIAL PRIMARY KEY,
                       email         VARCHAR(255) UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       name          VARCHAR(255),
                       created_at    TIMESTAMP DEFAULT NOW()
);

CREATE TABLE subjects (
                          id          BIGSERIAL PRIMARY KEY,
                          name        VARCHAR(255) NOT NULL,
                          description TEXT
);

CREATE TABLE sections (
                          id         BIGSERIAL PRIMARY KEY,
                          subject_id BIGINT NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
                          name       VARCHAR(255) NOT NULL,
                          description TEXT
);

CREATE TABLE tests (
                       id          BIGSERIAL PRIMARY KEY,
                       section_id  BIGINT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
                       title       VARCHAR(255) NOT NULL,
                       description TEXT,
                       created_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE questions (
                           id            BIGSERIAL PRIMARY KEY,
                           subject_id    BIGINT NOT NULL REFERENCES subjects(id),
    -- убрали test_id отсюда, вместо этого таблица test_questions ниже
                           question_text TEXT NOT NULL,
                           question_type VARCHAR(50) NOT NULL, --Тип вопроса: 'single_choice' (один ответ), 'multiple_choice' (несколько), 'short_answer' (студент пишет сам). От этого зависит как отображать вопрос
                           image_url     TEXT,
                           explanation   TEXT,
                           difficulty    SMALLINT NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
                           points        INTEGER DEFAULT 1,
                           created_at    TIMESTAMP DEFAULT NOW()
);

-- связка вопросов с тестами (вместо test_id в questions)
CREATE TABLE test_questions (
                                id             BIGSERIAL PRIMARY KEY,
                                test_id        BIGINT NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
                                question_id    BIGINT NOT NULL REFERENCES questions(id),
                                question_order INTEGER NOT NULL,
                                UNIQUE (test_id, question_id)
);

CREATE TABLE answers (
                         id          BIGSERIAL PRIMARY KEY,
                         question_id BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
                         answer_text TEXT NOT NULL,
                         is_correct  BOOLEAN DEFAULT FALSE
);

CREATE TABLE exams (
                       id               BIGSERIAL PRIMARY KEY,
                       subject_id       BIGINT NOT NULL REFERENCES subjects(id),
                       title            VARCHAR(255) NOT NULL,
                       description      TEXT,
                       duration_minutes INTEGER,
                       year             INTEGER,
                       created_at       TIMESTAMP DEFAULT NOW()
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

CREATE TABLE test_attempts (
                               id                 BIGSERIAL PRIMARY KEY,
                               user_id            BIGINT NOT NULL REFERENCES users(id),
                               subject_id         BIGINT NOT NULL REFERENCES subjects(id),
                               test_id            BIGINT REFERENCES tests(id),   -- добавили
                               exam_id            BIGINT REFERENCES exams(id),
                               mode               VARCHAR(20) NOT NULL CHECK (mode IN ('test', 'exam', 'practice')),
                               status             VARCHAR(20) DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed', 'abandoned')),
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
                              is_correct  BOOLEAN,
                              answered_at TIMESTAMP DEFAULT NOW()
);


-- тарифные планы
CREATE TABLE subscription_plans (
                                    id             BIGSERIAL PRIMARY KEY,
                                    name           VARCHAR(100) NOT NULL,  -- "Месяц", "Год"
                                    price          NUMERIC(10,2) NOT NULL, -- 990.00
                                    currency       VARCHAR(3) DEFAULT 'AMD',
                                    duration_days  INTEGER NOT NULL,       -- 30, 365
                                    is_active      BOOLEAN DEFAULT TRUE
);

-- подписки пользователей
CREATE TABLE user_subscriptions (
                                    id          BIGSERIAL PRIMARY KEY,
                                    user_id     BIGINT NOT NULL REFERENCES users(id),
                                    plan_id     BIGINT NOT NULL REFERENCES subscription_plans(id),
                                    status      VARCHAR(20) NOT NULL CHECK (status IN ('active', 'expired', 'cancelled')),
                                    started_at  TIMESTAMP NOT NULL DEFAULT NOW(),
                                    expires_at  TIMESTAMP NOT NULL,
                                    created_at  TIMESTAMP DEFAULT NOW()
);


CREATE TABLE payments (
                          id              BIGSERIAL PRIMARY KEY,
                          user_id         BIGINT NOT NULL REFERENCES users(id),
                          plan_id         BIGINT NOT NULL REFERENCES subscription_plans(id),
                          amount          NUMERIC(10,2) NOT NULL,
                          currency        VARCHAR(3) DEFAULT 'AMD',
                          status          VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'success', 'failed')),
                          provider        VARCHAR(30) NOT NULL,  -- 'idram', 'telcell', 'stripe'
                          provider_payment_id VARCHAR(255),      -- id платежа в системе провайдера
                          created_at      TIMESTAMP DEFAULT NOW(),
                          updated_at      TIMESTAMP DEFAULT NOW()
);