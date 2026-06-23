package domain

import (
	"context"
	"time"
)

type Attempt struct {
	ID               int64      `json:"id"`
	UserID           int64      `json:"user_id"`
	SubjectID        int64      `json:"subject_id"`
	TestID           *int64     `json:"test_id,omitempty"`
	ExamID           *int64     `json:"exam_id,omitempty"`
	Mode             string     `json:"mode"`
	Status           string     `json:"status"`
	StartedAt        time.Time  `json:"started_at"`
	FinishedAt       *time.Time `json:"finished_at,omitempty"`
	TimeLimitMinutes *int       `json:"time_limit_minutes,omitempty"`
	Score            int        `json:"score"`
	MaxScore         int        `json:"max_score"`
	Percentage       float64    `json:"percentage"`
}

// константы для поля Mode
const (
	AttemptModeTest     = "test"
	AttemptModeExam     = "exam"
	AttemptModePractice = "practice"
)

// константы для поля Status
const (
	AttemptStatusInProgress = "in_progress"
	AttemptStatusCompleted  = "completed"
	AttemptStatusAbandoned  = "abandoned"
)

type AttemptRepository interface {
	CreateAttempt(ctx context.Context, attempt *Attempt) error
	UpdateAttempt(ctx context.Context, attempt *Attempt) error
	GetAttemptByID(ctx context.Context, id int64) (*Attempt, error)
	GetInProgressAttemptByTestID(ctx context.Context, userID, testID int64) (*Attempt, error)
}
