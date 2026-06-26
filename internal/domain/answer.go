package domain

import "context"

type Answer struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	AnswerText string `json:"answer_text"`
	IsCorrect  bool   `json:"-"` // никогда не уходит в JSON клиенту
}

type AnswerRepository interface {
	GetAnswersByQuestionIDs(ctx context.Context, questionIDs []int64) ([]*Answer, error)
	GetAnswersByQuestionID(ctx context.Context, questionID int64) ([]*Answer, error)
	GetAnswerByID(ctx context.Context, id int64) (*Answer, error)
}
