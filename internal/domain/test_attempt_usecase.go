package domain

import "context"

// что приходит от фронтенда при ответе на вопрос
type SubmitAnswerInput struct {
	QuestionID int64
	AnswerID   *int64  // для single_choice
	AnswerIDs  []int64 // multiple_choice
	TextAnswer string  // для short_answer
}

// что возвращаем после ответа на вопрос
type SubmitAnswerResult struct {
	AttemptID   int64  `json:"attempt_id"`
	QuestionID  int64  `json:"question_id"`
	IsCorrect   *bool  `json:"is_correct,omitempty"` // nil если режим end_only
	Explanation string `json:"explanation,omitempty"`
}

// полный результат завершённой попытки
type AttemptResult struct {
	Attempt *Attempt               `json:"attempt"`
	Answers []*AttemptAnswerDetail `json:"answers"`
}

// детальный разбор одного вопроса в результатах
type AttemptAnswerDetail struct {
	Question      *Question   `json:"question"`
	UserAnswer    *UserAnswer `json:"user_answer"`
	CorrectAnswer *Answer     `json:"correct_answer,omitempty"`
}

type TestAttemptUseCase interface {
	// отдаёт тест с вопросами + текущее состояние попытки (если есть)
	GetTestForUser(ctx context.Context, userID, testID int64) (*Test, []*Question, []*Answer, *Attempt, error)

	// ответить на вопрос — создаёт попытку если её нет
	SubmitAnswer(ctx context.Context, userID, testID int64, input SubmitAnswerInput) (*SubmitAnswerResult, error)

	// начать тест с нуля (помечает старую попытку abandoned)
	Restart(ctx context.Context, userID, testID int64) (*Attempt, error)

	// завершить попытку — считает итоговый результат
	Finish(ctx context.Context, userID, attemptID int64) (*AttemptResult, error)

	// получить результат уже завершённой попытки
	GetResult(ctx context.Context, userID, attemptID int64) (*AttemptResult, error)
}
