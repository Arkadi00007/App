-- 000009_create_attempts.down.sql
DROP INDEX IF EXISTS idx_user_answers_question_id;
DROP INDEX IF EXISTS idx_user_answers_attempt_id;
DROP INDEX IF EXISTS idx_attempts_subject_id;
DROP INDEX IF EXISTS idx_attempts_user_status;
DROP INDEX IF EXISTS idx_attempts_user_id;
DROP TABLE IF EXISTS user_answers;
DROP TABLE IF EXISTS test_attempts;