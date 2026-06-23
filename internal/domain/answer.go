package domain

import "context"

type Answer struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	AnswerText string `json:"answer_text"`
	IsCorrect  bool   `json:"-"` // никогда не уходит в JSON клиенту
}

type AnswerRepository interface {
	GetByQuestionIDs(ctx context.Context, questionIDs []int64) ([]*Answer, error)
	GetByQuestionID(ctx context.Context, questionID int64) ([]*Answer, error)
	GetByID(ctx context.Context, id int64) (*Answer, error)
}
