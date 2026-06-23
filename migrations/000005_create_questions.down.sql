-- 000005_create_questions.down.sql
DROP INDEX IF EXISTS idx_questions_difficulty;
DROP INDEX IF EXISTS idx_questions_subject_id;
DROP TABLE IF EXISTS questions;