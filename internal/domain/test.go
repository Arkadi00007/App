package domain

import (
	"context"
	"time"
)

type Test struct {
	ID             int64     `json:"id"`
	SectionID      int64     `json:"section_id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	ShowAnswerMode string    `json:"show_answer_mode"` // immediate / end_only
	IsPremium      bool      `json:"is_premium"`
	CreatedAt      time.Time `json:"created_at"`
}

const (
	ShowAnswerModeImmediate = "immediate"
	ShowAnswerModeEndOnly   = "end_only"
)

type TestRepository interface {
	GetTestsBySectionID(ctx context.Context, sectionID int64) ([]*Test, error)
	GetTestByID(ctx context.Context, id int64) (*Test, error)
}
