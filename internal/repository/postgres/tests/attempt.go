package tests

import (
	"Test_App/internal/domain"
	"context"
	"database/sql"
)

type attemptRepository struct {
	db *sql.DB
}

func NewAttemptRepository(db *sql.DB) domain.AttemptRepository {
	return &attemptRepository{db: db}
}

// CreateAttempt — создаёт новую попытку прохождения теста
func (r *attemptRepository) CreateAttempt(ctx context.Context, attempt *domain.Attempt) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO test_attempts (user_id, subject_id, test_id, exam_id, mode, status, time_limit_minutes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, started_at
	`,
		attempt.UserID,
		attempt.SubjectID,
		attempt.TestID,
		attempt.ExamID,
		attempt.Mode,
		attempt.Status,
		attempt.TimeLimitMinutes,
	).Scan(&attempt.ID, &attempt.StartedAt)
}

// UpdateAttempt — обновляет попытку (используется при Finish — закрытии попытки с результатом)
func (r *attemptRepository) UpdateAttempt(ctx context.Context, attempt *domain.Attempt) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE test_attempts SET
			status      = $1,
			finished_at = $2,
			score       = $3,
			max_score   = $4,
			percentage  = $5
		WHERE id = $6
	`,
		attempt.Status,
		attempt.FinishedAt,
		attempt.Score,
		attempt.MaxScore,
		attempt.Percentage,
		attempt.ID,
	)
	return err
}

// GetAttemptByID — конкретная попытка по id
func (r *attemptRepository) GetAttemptByID(ctx context.Context, id int64) (*domain.Attempt, error) {
	a := &domain.Attempt{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, subject_id, test_id, exam_id, mode, status,
		       started_at, finished_at, time_limit_minutes, score, max_score, percentage
		FROM test_attempts
		WHERE id = $1
	`, id).Scan(
		&a.ID, &a.UserID, &a.SubjectID, &a.TestID, &a.ExamID, &a.Mode, &a.Status,
		&a.StartedAt, &a.FinishedAt, &a.TimeLimitMinutes, &a.Score, &a.MaxScore, &a.Percentage,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return a, nil
}

// GetInProgressAttemptByTestID — ищем незавершённую попытку юзера по конкретному тесту
// используется в SubmitAnswer, чтобы решить — создавать новую попытку или продолжать старую
func (r *attemptRepository) GetInProgressAttemptByTestID(ctx context.Context, userID, testID int64) (*domain.Attempt, error) {
	a := &domain.Attempt{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, subject_id, test_id, exam_id, mode, status,
		       started_at, finished_at, time_limit_minutes, score, max_score, percentage
		FROM test_attempts
		WHERE user_id = $1
		  AND test_id = $2
		  AND status = $3
		LIMIT 1
	`, userID, testID, domain.AttemptStatusInProgress).Scan(
		&a.ID, &a.UserID, &a.SubjectID, &a.TestID, &a.ExamID, &a.Mode, &a.Status,
		&a.StartedAt, &a.FinishedAt, &a.TimeLimitMinutes, &a.Score, &a.MaxScore, &a.Percentage,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return a, nil
}
