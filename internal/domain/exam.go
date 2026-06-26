package domain

import (
	"context"
	"time"
)

type Exam struct {
	ID              int64     `json:"id"`
	SubjectID       int64     `json:"subject_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description,omitempty"`
	DurationMinutes int       `json:"duration_minutes"`
	Year            int       `json:"year,omitempty"`
	IsPremium       bool      `json:"is_premium"`
	CreatedAt       time.Time `json:"created_at"`
}

type ExamRepository interface {
	GetExamsBySubjectID(ctx context.Context, subjectID int64) ([]*Exam, error)
	GetExamByID(ctx context.Context, id int64) (*Exam, error)
}
