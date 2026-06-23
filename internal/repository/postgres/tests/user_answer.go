package tests

import (
	"Test_App/internal/domain"
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type userAnswerRepository struct {
	db *sql.DB
}

func NewUserAnswerRepository(db *sql.DB) *userAnswerRepository {
	return &userAnswerRepository{db: db}
}

// CreateUserAnswer — сохраняет ответ студента на вопрос (первый ответ)
func (r *userAnswerRepository) CreateUserAnswer(ctx context.Context, ua *domain.UserAnswer) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO user_answers (attempt_id, question_id, answer_id, answer_ids, text_answer, is_correct)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, answered_at
	`,
		ua.AttemptID,
		ua.QuestionID,
		ua.AnswerID,
		pq.Array(ua.AnswerIDs),
		ua.TextAnswer,
		ua.IsCorrect,
	).Scan(&ua.ID, &ua.AnsweredAt)
}

// UpdateUserAnswer — обновляет существующий ответ (студент изменил выбор, вернувшись назад)
func (r *userAnswerRepository) UpdateUserAnswer(ctx context.Context, ua *domain.UserAnswer) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE user_answers SET
			answer_id   = $1,
			answer_ids  = $2,
			text_answer = $3,
			is_correct  = $4,
			answered_at = NOW()
		WHERE id = $5
	`,
		ua.AnswerID,
		pq.Array(ua.AnswerIDs),
		ua.TextAnswer,
		ua.IsCorrect,
		ua.ID,
	)
	return err
}

// GetUserAnswerByAttemptAndQuestion — проверяем, отвечал ли студент уже на этот вопрос в этой попытке
func (r *userAnswerRepository) GetUserAnswerByAttemptAndQuestion(ctx context.Context, attemptID, questionID int64) (*domain.UserAnswer, error) {
	ua := &domain.UserAnswer{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, attempt_id, question_id, answer_id, answer_ids, text_answer, is_correct, answered_at
		FROM user_answers
		WHERE attempt_id = $1 AND question_id = $2
	`, attemptID, questionID).Scan(
		&ua.ID, &ua.AttemptID, &ua.QuestionID, &ua.AnswerID,
		pq.Array(&ua.AnswerIDs),
		&ua.TextAnswer, &ua.IsCorrect, &ua.AnsweredAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ua, nil
}

// GetUserAnswersByAttemptID — все ответы попытки (для подсчёта результата и для экрана с разбором)
func (r *userAnswerRepository) GetUserAnswersByAttemptID(ctx context.Context, attemptID int64) ([]*domain.UserAnswer, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, attempt_id, question_id, answer_id, answer_ids, text_answer, is_correct, answered_at
		FROM user_answers
		WHERE attempt_id = $1
		ORDER BY answered_at
	`, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userAnswers []*domain.UserAnswer
	for rows.Next() {
		ua := &domain.UserAnswer{}
		if err := rows.Scan(
			&ua.ID, &ua.AttemptID, &ua.QuestionID, &ua.AnswerID,
			pq.Array(&ua.AnswerIDs),
			&ua.TextAnswer, &ua.IsCorrect, &ua.AnsweredAt,
		); err != nil {
			return nil, err
		}
		userAnswers = append(userAnswers, ua)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userAnswers, nil
}
