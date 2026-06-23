package tests

import (
	"Test_App/internal/domain"
	"context"
	"database/sql"
)

type questionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) domain.QuestionRepository {
	return &questionRepository{db: db}
}

// GetByTestID — все вопросы конкретного теста, через связку test_questions
func (r *questionRepository) GetQuestionsByTestID(ctx context.Context, testID int64) ([]*domain.Question, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT q.id, q.subject_id, q.question_text, q.question_type,
		       q.image_url, q.explanation, q.difficulty, q.points, q.created_at
		FROM questions q
		JOIN test_questions tq ON tq.question_id = q.id
		WHERE tq.test_id = $1
		ORDER BY tq.question_order
	`, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*domain.Question
	for rows.Next() {
		q := &domain.Question{}
		if err := rows.Scan(
			&q.ID, &q.SubjectID, &q.QuestionText, &q.QuestionType,
			&q.ImageURL, &q.Explanation, &q.Difficulty, &q.Points, &q.CreatedAt,
		); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

func (r *questionRepository) GetQuestionByID(ctx context.Context, id int64) (*domain.Question, error) {
	q := &domain.Question{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, subject_id, question_text, question_type,
		       image_url, explanation, difficulty, points, created_at
		FROM questions
		WHERE id = $1
	`, id).Scan(
		&q.ID, &q.SubjectID, &q.QuestionText, &q.QuestionType,
		&q.ImageURL, &q.Explanation, &q.Difficulty, &q.Points, &q.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return q, nil
}
