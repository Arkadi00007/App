package domain

import (
	"context"
	"time"
)

type UserAnswer struct {
	ID         int64   `json:"id"`
	AttemptID  int64   `json:"attempt_id"`
	QuestionID int64   `json:"question_id"`
	AnswerID   *int64  `json:"answer_id,omitempty"`
	AnswerIDs  []int64 `json:"selected_ids,omitempty"` // для multiple_choice

	TextAnswer string    `json:"text_answer,omitempty"`
	IsCorrect  *bool     `json:"is_correct,omitempty"`
	AnsweredAt time.Time `json:"answered_at"`
}

type UserAnswerRepository interface {
	CreateUserAnswer(ctx context.Context, answer *UserAnswer) error
	// обновляем существующий ответ (если студент поменял ответ на вопрос)
	UpdateUserAnswer(ctx context.Context, answer *UserAnswer) error
	// ищем — отвечал ли уже студент на этот вопрос в этой попытке
	GetUserAnswerByAttemptAndQuestion(ctx context.Context, attemptID, questionID int64) (*UserAnswer, error)
	GetUserAnswersByAttemptID(ctx context.Context, attemptID int64) ([]*UserAnswer, error)
}
