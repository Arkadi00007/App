package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"Test_App/internal/domain"
)

var (
	ErrTestNotFound    = errors.New("test not found")
	ErrAttemptNotFound = errors.New("attempt not found")
	ErrNotYourAttempt  = errors.New("attempt does not belong to this user")
	ErrAlreadyFinished = errors.New("attempt is already finished")
	ErrWrongQuestion   = errors.New("question does not belong to this test")
	ErrPremiumRequired = errors.New("для доступа к этому тесту требуется подписка")
)

type testAttemptUseCase struct {
	testRepo         domain.TestRepository
	questionRepo     domain.QuestionRepository
	answerRepo       domain.AnswerRepository
	attemptRepo      domain.AttemptRepository
	userAnswerRepo   domain.UserAnswerRepository
	subscriptionRepo domain.SubscriptionRepository
}

func NewTestAttemptUseCase(
	testRepo domain.TestRepository,
	questionRepo domain.QuestionRepository,
	answerRepo domain.AnswerRepository,
	attemptRepo domain.AttemptRepository,
	userAnswerRepo domain.UserAnswerRepository,
	subscriptionRepo domain.SubscriptionRepository,
) domain.TestAttemptUseCase {
	return &testAttemptUseCase{
		testRepo:       testRepo,
		questionRepo:   questionRepo,
		answerRepo:     answerRepo,
		attemptRepo:    attemptRepo,
		userAnswerRepo: userAnswerRepo,
	}
}

// =============================================================================
// GetTestForUser — студент открыл тест
// =============================================================================

func (uc *testAttemptUseCase) GetTestForUser(
	ctx context.Context,
	userID, testID int64,
) (*domain.Test, []*domain.Question, []*domain.Answer, *domain.Attempt, error) {

	// 1. загружаем тест
	test, err := uc.testRepo.GetTestByID(ctx, testID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if test == nil {
		return nil, nil, nil, nil, ErrTestNotFound
	}

	if err := uc.checkTestAccess(ctx, userID, test); err != nil {
		return nil, nil, nil, nil, err
	}

	// 2. загружаем вопросы теста
	questions, err := uc.questionRepo.GetQuestionsByTestID(ctx, testID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// 3. собираем ID всех вопросов
	questionIDs := make([]int64, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.ID
	}

	// 4. загружаем варианты ответов одним запросом (is_correct скрыт через json:"-")
	answers, err := uc.answerRepo.GetAnswersByQuestionIDs(ctx, questionIDs)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// 5. ищем незавершённую попытку — nil если студент заходит впервые
	attempt, err := uc.attemptRepo.GetInProgressAttemptByTestID(ctx, userID, testID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return test, questions, answers, attempt, nil
}

// =============================================================================
// SubmitAnswer — студент ответил на вопрос
// =============================================================================

func (uc *testAttemptUseCase) SubmitAnswer(
	ctx context.Context,
	userID, testID int64,
	input domain.SubmitAnswerInput,
) (*domain.SubmitAnswerResult, error) {

	// 1. загружаем тест — нужен show_answer_mode и subject_id
	test, err := uc.testRepo.GetTestByID(ctx, testID)
	if err != nil {
		return nil, err
	}
	if test == nil {
		return nil, ErrTestNotFound
	}

	if err := uc.checkTestAccess(ctx, userID, test); err != nil {
		return nil, err
	}

	// 2. загружаем вопрос — нужен question_type и explanation
	question, err := uc.questionRepo.GetQuestionByID(ctx, input.QuestionID)
	if err != nil {
		return nil, err
	}
	if question == nil {
		return nil, ErrWrongQuestion
	}

	// 3. находим или создаём попытку
	// subject_id берём из теста — не делаем лишний запрос к вопросам
	attempt, err := uc.getOrCreateAttempt(ctx, userID, testID, test.SubjectID)
	if err != nil {
		return nil, err
	}

	// 4. проверяем правильность ответа
	isCorrect, err := uc.checkAnswer(ctx, question, input)
	if err != nil {
		return nil, err
	}

	// 5. сохраняем ответ (создаём или обновляем если студент вернулся назад)
	if err := uc.saveUserAnswer(ctx, attempt.ID, question.ID, input, isCorrect); err != nil {
		return nil, err
	}

	// 6. формируем результат
	result := &domain.SubmitAnswerResult{
		AttemptID:  attempt.ID,
		QuestionID: question.ID,
	}

	if test.ShowAnswerMode == domain.ShowAnswerModeImmediate {
		result.IsCorrect = &isCorrect
		// explanation показываем только при неправильном ответе
		if !isCorrect {
			result.Explanation = question.Explanation
		}
	}
	// end_only: IsCorrect остаётся nil — фронт просто переходит дальше

	return result, nil
}

// =============================================================================
// Restart — студент нажал "начать заново"
// =============================================================================

func (uc *testAttemptUseCase) Restart(
	ctx context.Context,
	userID, testID int64,
) (*domain.Attempt, error) {

	// 1. загружаем тест — нужен subject_id
	test, err := uc.testRepo.GetTestByID(ctx, testID)
	if err != nil {
		return nil, err
	}
	if test == nil {
		return nil, ErrTestNotFound
	}

	// 2. ищем активную попытку и помечаем abandoned если есть
	existing, err := uc.attemptRepo.GetInProgressAttemptByTestID(ctx, userID, testID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		existing.Status = domain.AttemptStatusAbandoned
		now := time.Now()
		existing.FinishedAt = &now
		if err := uc.attemptRepo.UpdateAttempt(ctx, existing); err != nil {
			return nil, err
		}
	}

	// 3. создаём новую попытку
	newAttempt := &domain.Attempt{
		UserID:    userID,
		SubjectID: test.SubjectID, // берём из теста — без лишних запросов
		TestID:    &testID,
		Mode:      domain.AttemptModeTest,
		Status:    domain.AttemptStatusInProgress,
	}
	if err := uc.attemptRepo.CreateAttempt(ctx, newAttempt); err != nil {
		return nil, err
	}

	return newAttempt, nil
}

// =============================================================================
// Finish — студент нажал "завершить тест"
// =============================================================================

func (uc *testAttemptUseCase) Finish(
	ctx context.Context,
	userID, attemptID int64,
) (*domain.AttemptResult, error) {

	// 1. загружаем попытку
	attempt, err := uc.attemptRepo.GetAttemptByID(ctx, attemptID)
	if err != nil {
		return nil, err
	}
	if attempt == nil {
		return nil, ErrAttemptNotFound
	}

	// 2. попытка должна принадлежать этому пользователю
	if attempt.UserID != userID {
		return nil, ErrNotYourAttempt
	}

	// 3. попытка должна быть ещё не завершена
	if attempt.Status != domain.AttemptStatusInProgress {
		return nil, ErrAlreadyFinished
	}

	// 4. загружаем тест — нужен show_answer_mode
	test, err := uc.testRepo.GetTestByID(ctx, *attempt.TestID)
	if err != nil {
		return nil, err
	}

	// 5. загружаем вопросы и ответы студента
	questions, err := uc.questionRepo.GetQuestionsByTestID(ctx, *attempt.TestID)
	if err != nil {
		return nil, err
	}

	userAnswers, err := uc.userAnswerRepo.GetUserAnswersByAttemptID(ctx, attemptID)
	if err != nil {
		return nil, err
	}

	// 6. для режима end_only — проверяем все ответы здесь (до этого is_correct не проставлялся)
	if test.ShowAnswerMode == domain.ShowAnswerModeEndOnly {
		userAnswers, err = uc.checkAndUpdateAllAnswers(ctx, userAnswers, questions)
		if err != nil {
			return nil, err
		}
	}

	// 7. считаем score и max_score
	score, maxScore := uc.calculateScore(userAnswers, questions)

	// 8. считаем процент
	var percentage float64
	if maxScore > 0 {
		percentage = float64(score) / float64(maxScore) * 100
	}

	// 9. закрываем попытку
	now := time.Now()
	attempt.Status = domain.AttemptStatusCompleted
	attempt.FinishedAt = &now
	attempt.Score = score
	attempt.MaxScore = maxScore
	attempt.Percentage = percentage

	if err := uc.attemptRepo.UpdateAttempt(ctx, attempt); err != nil {
		return nil, err
	}

	// 10. собираем детальный разбор
	return uc.buildAttemptResult(ctx, attempt, questions, userAnswers)
}

// =============================================================================
// GetResult — результат уже завершённой попытки
// =============================================================================

func (uc *testAttemptUseCase) GetResult(
	ctx context.Context,
	userID, attemptID int64,
) (*domain.AttemptResult, error) {

	attempt, err := uc.attemptRepo.GetAttemptByID(ctx, attemptID)
	if err != nil {
		return nil, err
	}
	if attempt == nil {
		return nil, ErrAttemptNotFound
	}
	if attempt.UserID != userID {
		return nil, ErrNotYourAttempt
	}
	if attempt.Status != domain.AttemptStatusCompleted {
		return nil, errors.New("attempt is not completed yet")
	}

	questions, err := uc.questionRepo.GetQuestionsByTestID(ctx, *attempt.TestID)
	if err != nil {
		return nil, err
	}

	userAnswers, err := uc.userAnswerRepo.GetUserAnswersByAttemptID(ctx, attemptID)
	if err != nil {
		return nil, err
	}

	return uc.buildAttemptResult(ctx, attempt, questions, userAnswers)
}

// =============================================================================
// Вспомогательные методы
// =============================================================================

// getOrCreateAttempt — возвращает активную попытку или создаёт новую
func (uc *testAttemptUseCase) getOrCreateAttempt(
	ctx context.Context,
	userID, testID, subjectID int64,
) (*domain.Attempt, error) {

	attempt, err := uc.attemptRepo.GetInProgressAttemptByTestID(ctx, userID, testID)
	if err != nil {
		return nil, err
	}
	if attempt != nil {
		return attempt, nil
	}

	newAttempt := &domain.Attempt{
		UserID:    userID,
		SubjectID: subjectID,
		TestID:    &testID,
		Mode:      domain.AttemptModeTest,
		Status:    domain.AttemptStatusInProgress,
	}
	if err := uc.attemptRepo.CreateAttempt(ctx, newAttempt); err != nil {
		return nil, err
	}

	return newAttempt, nil
}

// checkAnswer — проверяет правильность в зависимости от типа вопроса
func (uc *testAttemptUseCase) checkAnswer(
	ctx context.Context,
	question *domain.Question,
	input domain.SubmitAnswerInput,
) (bool, error) {

	switch question.QuestionType {

	case domain.QuestionTypeSingleChoice:
		if input.AnswerID == nil {
			return false, nil
		}
		answer, err := uc.answerRepo.GetAnswerByID(ctx, *input.AnswerID)
		if err != nil {
			return false, err
		}
		if answer == nil {
			return false, nil
		}
		// защита от подмены — ответ должен принадлежать именно этому вопросу
		if answer.QuestionID != question.ID {
			return false, ErrWrongQuestion
		}
		return answer.IsCorrect, nil

	case domain.QuestionTypeMultipleChoice:
		allAnswers, err := uc.answerRepo.GetAnswersByQuestionID(ctx, question.ID)
		if err != nil {
			return false, err
		}

		correctIDs := make(map[int64]bool)
		for _, a := range allAnswers {
			if a.IsCorrect {
				correctIDs[a.ID] = true
			}
		}

		selectedIDs := make(map[int64]bool)
		for _, id := range input.AnswerIDs {
			selectedIDs[id] = true
		}

		if len(selectedIDs) != len(correctIDs) {
			return false, nil
		}
		for id := range correctIDs {
			if !selectedIDs[id] {
				return false, nil
			}
		}
		return true, nil

	case domain.QuestionTypeShortAnswer:
		allAnswers, err := uc.answerRepo.GetAnswersByQuestionID(ctx, question.ID)
		if err != nil {
			return false, err
		}
		got := strings.TrimSpace(strings.ToLower(input.TextAnswer))
		for _, a := range allAnswers {
			if a.IsCorrect {
				expected := strings.TrimSpace(strings.ToLower(a.AnswerText))
				if expected == got {
					return true, nil // нашли совпадение среди всех правильных вариантов
				}
				// не совпало — продолжаем проверять остальные варианты
			}
		}
		return false, nil
	}

	return false, nil
}

// saveUserAnswer — создаёт или обновляет ответ студента
func (uc *testAttemptUseCase) saveUserAnswer(
	ctx context.Context,
	attemptID, questionID int64,
	input domain.SubmitAnswerInput,
	isCorrect bool,
) error {

	existing, err := uc.userAnswerRepo.GetUserAnswerByAttemptAndQuestion(ctx, attemptID, questionID)
	if err != nil {
		return err
	}

	if existing == nil {
		ua := &domain.UserAnswer{
			AttemptID:  attemptID,
			QuestionID: questionID,
			AnswerID:   input.AnswerID,
			AnswerIDs:  input.AnswerIDs,
			TextAnswer: input.TextAnswer,
			IsCorrect:  &isCorrect,
		}
		return uc.userAnswerRepo.CreateUserAnswer(ctx, ua)
	}

	// студент вернулся назад и поменял ответ
	existing.AnswerID = input.AnswerID
	existing.AnswerIDs = input.AnswerIDs
	existing.TextAnswer = input.TextAnswer
	existing.IsCorrect = &isCorrect
	return uc.userAnswerRepo.UpdateUserAnswer(ctx, existing)
}

// checkAndUpdateAllAnswers — для режима end_only
// проверяет все ответы студента и сохраняет is_correct в БД
func (uc *testAttemptUseCase) checkAndUpdateAllAnswers(
	ctx context.Context,
	userAnswers []*domain.UserAnswer,
	questions []*domain.Question,
) ([]*domain.UserAnswer, error) {

	questionMap := make(map[int64]*domain.Question, len(questions))
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	for _, ua := range userAnswers {
		question, ok := questionMap[ua.QuestionID]
		if !ok {
			continue
		}

		input := domain.SubmitAnswerInput{
			QuestionID: ua.QuestionID,
			AnswerID:   ua.AnswerID,
			AnswerIDs:  ua.AnswerIDs,
			TextAnswer: ua.TextAnswer,
		}

		isCorrect, err := uc.checkAnswer(ctx, question, input)
		if err != nil {
			return nil, err
		}

		ua.IsCorrect = &isCorrect
		if err := uc.userAnswerRepo.UpdateUserAnswer(ctx, ua); err != nil {
			return nil, err
		}
	}

	return userAnswers, nil
}

// calculateScore — считает набранные и максимальные баллы
func (uc *testAttemptUseCase) calculateScore(
	userAnswers []*domain.UserAnswer,
	questions []*domain.Question,
) (score, maxScore int) {

	questionMap := make(map[int64]*domain.Question, len(questions))
	for _, q := range questions {
		questionMap[q.ID] = q
		maxScore += q.Points
	}

	for _, ua := range userAnswers {
		if ua.IsCorrect != nil && *ua.IsCorrect {
			if q, ok := questionMap[ua.QuestionID]; ok {
				score += q.Points
			}
		}
	}

	return score, maxScore
}

// buildAttemptResult — собирает детальный разбор попытки
func (uc *testAttemptUseCase) buildAttemptResult(
	ctx context.Context,
	attempt *domain.Attempt,
	questions []*domain.Question,
	userAnswers []*domain.UserAnswer,
) (*domain.AttemptResult, error) {

	userAnswerMap := make(map[int64]*domain.UserAnswer, len(userAnswers))
	for _, ua := range userAnswers {
		userAnswerMap[ua.QuestionID] = ua
	}

	details := make([]*domain.AttemptAnswerDetail, 0, len(questions))

	for _, q := range questions {
		detail := &domain.AttemptAnswerDetail{
			Question:   q,
			UserAnswer: userAnswerMap[q.ID], // nil если студент не ответил
		}

		allAnswers, err := uc.answerRepo.GetAnswersByQuestionID(ctx, q.ID)
		if err != nil {
			return nil, err
		}

		switch q.QuestionType {
		case domain.QuestionTypeSingleChoice, domain.QuestionTypeShortAnswer:
			// один правильный вариант
			for _, a := range allAnswers {
				if a.IsCorrect {
					detail.CorrectAnswer = a
					break
				}
			}
		case domain.QuestionTypeMultipleChoice:
			// несколько правильных вариантов
			for _, a := range allAnswers {
				if a.IsCorrect {
					detail.CorrectAnswers = append(detail.CorrectAnswers, a)
				}
			}
		}

		details = append(details, detail)
	}

	return &domain.AttemptResult{
		Attempt: attempt,
		Answers: details,
	}, nil
}

// checkTestAccess — проверяет имеет ли студент доступ к тесту
// если тест бесплатный — доступ всегда есть
// если тест премиум — нужна активная подписка
func (uc *testAttemptUseCase) checkTestAccess(ctx context.Context, userID int64, test *domain.Test) error {
	if !test.IsPremium {
		return nil // бесплатный тест — проверка не нужна
	}

	hasSubscription, err := uc.subscriptionRepo.HasActiveSubscription(ctx, userID)
	if err != nil {
		return err
	}
	if !hasSubscription {
		return ErrPremiumRequired
	}

	return nil
}
