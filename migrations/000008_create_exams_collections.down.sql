-- 000008_create_exams_collections.down.sql
DROP INDEX IF EXISTS idx_collection_questions_collection_id;
DROP INDEX IF EXISTS idx_collections_subject_id;
DROP INDEX IF EXISTS idx_exam_questions_exam_id_order;
DROP INDEX IF EXISTS idx_exams_subject_id;
DROP TABLE IF EXISTS collection_questions;
DROP TABLE IF EXISTS collections;
DROP TABLE IF EXISTS exam_questions;
DROP TABLE IF EXISTS exams;