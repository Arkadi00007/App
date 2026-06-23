package tests

import (
	"Test_App/internal/domain"
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type answerRepository struct {
	db *sql.DB
}

func NewAnswerRepository(db *sql.DB) domain.AnswerRepository {
	return &answerRepository{db: db}
}

// GetByQuestionIDs — варианты ответов сразу для нескольких вопросов
// (используется при открытии теста — грузим ответы всех вопросов одним запросом)
func (r *answerRepository) GetByQuestionIDs(ctx context.Context, questionIDs []int64) ([]*domain.Answer, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, question_id, answer_text, is_correct
		FROM answers
		WHERE question_id = ANY($1)
	`, pq.Array(questionIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []*domain.Answer
	for rows.Next() {
		a := &domain.Answer{}
		if err := rows.Scan(&a.ID, &a.QuestionID, &a.AnswerText, &a.IsCorrect); err != nil {
			return nil, err
		}
		answers = append(answers, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return answers, nil
}

// GetByQuestionID — варианты ответов ОДНОГО вопроса
// (используется при проверке multiple_choice и short_answer)
func (r *answerRepository) GetByQuestionID(ctx context.Context, questionID int64) ([]*domain.Answer, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, question_id, answer_text, is_correct
		FROM answers
		WHERE question_id = $1
	`, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []*domain.Answer
	for rows.Next() {
		a := &domain.Answer{}
		if err := rows.Scan(&a.ID, &a.QuestionID, &a.AnswerText, &a.IsCorrect); err != nil {
			return nil, err
		}
		answers = append(answers, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return answers, nil
}

func (r *answerRepository) GetByID(ctx context.Context, id int64) (*domain.Answer, error) {
	a := &domain.Answer{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, question_id, answer_text, is_correct
		FROM answers
		WHERE id = $1
	`, id).Scan(&a.ID, &a.QuestionID, &a.AnswerText, &a.IsCorrect)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return a, nil
}
