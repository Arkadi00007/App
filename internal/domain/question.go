package domain

import (
	"context"
	"time"
)

type Question struct {
	ID           int64     `json:"id"`
	SubjectID    int64     `json:"subject_id"`
	QuestionText string    `json:"question_text"`
	QuestionType string    `json:"question_type"` // single_choice, multiple_choice, short_answer
	ImageURL     string    `json:"image_url,omitempty"`
	Explanation  string    `json:"explanation,omitempty"`
	Difficulty   int       `json:"difficulty"`
	Points       int       `json:"points"`
	CreatedAt    time.Time `json:"created_at"`
}

const (
	QuestionTypeSingleChoice   = "single_choice"
	QuestionTypeMultipleChoice = "multiple_choice"
	QuestionTypeShortAnswer    = "short_answer"
)

type QuestionRepository interface {
	GetQuestionsByTestID(ctx context.Context, testID int64) ([]*Question, error)
	GetQuestionByID(ctx context.Context, id int64) (*Question, error)
}
