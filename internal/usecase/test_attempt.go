package usecase

//
//import (
//	"context"
//	"errors"
//	"strings"
//	"time"
//
//	"Test_App/internal/domain"
//)
//
//// ошибки которые usecase может вернуть — handler их поймает и отдаст нужный HTTP статус
//var (
//	ErrTestNotFound    = errors.New("test not found")
//	ErrAttemptNotFound = errors.New("attempt not found")
//	ErrNotYourAttempt  = errors.New("attempt does not belong to this user")
//	ErrAlreadyFinished = errors.New("attempt is already finished")
//	ErrWrongQuestion   = errors.New("question does not belong to this test")
//)
//
//// testAttemptUseCase — реализация интерфейса domain.TestAttemptUseCase
//type testAttemptUseCase struct {
//	testRepo       domain.TestRepository
//	questionRepo   domain.QuestionRepository
//	answerRepo     domain.AnswerRepository
//	attemptRepo    domain.AttemptRepository
//	userAnswerRepo domain.UserAnswerRepository
//}
//
//// NewTestAttemptUseCase — конструктор, принимает все нужные репозитории
//func NewTestAttemptUseCase(
//	testRepo domain.TestRepository,
//	questionRepo domain.QuestionRepository,
//	answerRepo domain.AnswerRepository,
//	attemptRepo domain.AttemptRepository,
//	userAnswerRepo domain.UserAnswerRepository,
//) domain.TestAttemptUseCase {
//	return &testAttemptUseCase{
//		testRepo:       testRepo,
//		questionRepo:   questionRepo,
//		answerRepo:     answerRepo,
//		attemptRepo:    attemptRepo,
//		userAnswerRepo: userAnswerRepo,
//	}
//}
//
//// =============================================================================
//// GetTestForUser — студент открыл тест
//// Возвращает: тест, вопросы, варианты ответов, текущую попытку (если есть)
//// =============================================================================
//
//func (uc *testAttemptUseCase) GetTestForUser(
//	ctx context.Context,
//	userID, testID int64,
//) (*domain.Test, []*domain.Question, []*domain.Answer, *domain.Attempt, error) {
//
//	// 1. Загружаем тест
//	test, err := uc.testRepo.GetTestByID(ctx, testID)
//	if err != nil {
//		return nil, nil, nil, nil, err
//	}
//	if test == nil {
//		return nil, nil, nil, nil, ErrTestNotFound
//	}
//
//	// 2. Загружаем вопросы теста (уже отсортированы по question_order в репозитории)
//	questions, err := uc.questionRepo.GetQuestionsByTestID(ctx, testID)
//	if err != nil {
//		return nil, nil, nil, nil, err
//	}
//
//	// 3. Собираем все ID вопросов чтобы загрузить варианты ответов одним запросом
//	questionIDs := make([]int64, len(questions))
//	for i, q := range questions {
//		questionIDs[i] = q.ID
//	}
//
//	// 4. Загружаем все варианты ответов для всех вопросов сразу
//	// IsCorrect в Answer помечен json:"-" — клиент не увидит правильные ответы
//	answers, err := uc.answerRepo.GetByQuestionIDs(ctx, questionIDs)
//	if err != nil {
//		return nil, nil, nil, nil, err
//	}
//
//	// 5. Ищем незавершённую попытку этого студента по этому тесту
//	// Может быть nil — значит студент открывает тест впервые
//	attempt, err := uc.attemptRepo.GetInProgressAttemptByTestID(ctx, userID, testID)
//	if err != nil {
//		return nil, nil, nil, nil, err
//	}
//
//	return test, questions, answers, attempt, nil
//}
//
//// =============================================================================
//// SubmitAnswer — студент ответил на вопрос
//// Создаёт попытку если её нет, сохраняет ответ, возвращает результат
//// =============================================================================
//
//func (uc *testAttemptUseCase) SubmitAnswer(
//	ctx context.Context,
//	userID, testID int64,
//	input domain.SubmitAnswerInput,
//) (*domain.SubmitAnswerResult, error) {
//
//	// 1. Загружаем тест — нужен show_answer_mode
//	test, err := uc.testRepo.GetTestByID(ctx, testID)
//	if err != nil {
//		return nil, err
//	}
//	if test == nil {
//		return nil, ErrTestNotFound
//	}
//
//	// 2. Загружаем вопрос — нужен question_type, explanation, subject_id
//	question, err := uc.questionRepo.GetQuestionByID(ctx, input.QuestionID)
//	if err != nil {
//		return nil, err
//	}
//	if question == nil {
//		return nil, ErrWrongQuestion
//	}
//
//	// 3. Находим или создаём попытку
//	attempt, err := uc.getOrCreateAttempt(ctx, userID, testID, question.SubjectID)
//	if err != nil {
//		return nil, err
//	}
//
//	// 4. Проверяем правильность ответа
//	isCorrect, err := uc.checkAnswer(ctx, question, input)
//	if err != nil {
//		return nil, err
//	}
//
//	// 5. Сохраняем ответ — создаём или обновляем если студент вернулся назад
//	err = uc.saveUserAnswer(ctx, attempt.ID, question.ID, input, isCorrect)
//	if err != nil {
//		return nil, err
//	}
//
//	// 6. Формируем результат
//	// Если режим end_only — скрываем is_correct и explanation до конца теста
//	result := &domain.SubmitAnswerResult{
//		AttemptID:  attempt.ID,
//		QuestionID: question.ID,
//	}
//
//	if test.ShowAnswerMode == domain.ShowAnswerModeImmediate {
//		result.IsCorrect = &isCorrect
//		if !isCorrect {
//			// explanation показываем только при неправильном ответе в immediate режиме
//			result.Explanation = question.Explanation
//		}
//	}
//	// при end_only: IsCorrect остаётся nil — фронт просто переходит к следующему вопросу
//
//	return result, nil
//}
//
//// =============================================================================
//// Restart — студент нажал "начать заново"
//// Помечает старую попытку как abandoned, создаёт новую
//// =============================================================================
//
//func (uc *testAttemptUseCase) Restart(
//	ctx context.Context,
//	userID, testID int64,
//) (*domain.Attempt, error) {
//
//	// 1. Ищем активную попытку
//	existing, err := uc.attemptRepo.GetInProgressAttemptByTestID(ctx, userID, testID)
//	if err != nil {
//		return nil, err
//	}
//
//	// 2. Если есть — помечаем как abandoned
//	if existing != nil {
//		existing.Status = domain.AttemptStatusAbandoned
//		now := time.Now()
//		existing.FinishedAt = &now
//		if err := uc.attemptRepo.UpdateAttempt(ctx, existing); err != nil {
//			return nil, err
//		}
//	}
//
//	// 3. Загружаем тест чтобы получить subject_id через вопросы
//	// (subject_id нужен для попытки, берём из первого вопроса)
//	questions, err := uc.questionRepo.GetQuestionsByTestID(ctx, testID)
//	if err != nil {
//		return nil, err
//	}
//
//	var subjectID int64
//	if len(questions) > 0 {
//		subjectID = questions[0].SubjectID
//	}
//
//	// 4. Создаём новую попытку
//	newAttempt := &domain.Attempt{
//		UserID:    userID,
//		SubjectID: subjectID,
//		TestID:    &testID,
//		Mode:      domain.AttemptModeTest,
//		Status:    domain.AttemptStatusInProgress,
//	}
//
//	if err := uc.attemptRepo.CreateAttempt(ctx, newAttempt); err != nil {
//		return nil, err
//	}
//
//	return newAttempt, nil
//}
//
//// =============================================================================
//// Finish — студент нажал "завершить тест"
//// Считает итоговый результат, помечает попытку completed
//// =============================================================================
//
//func (uc *testAttemptUseCase) Finish(
//	ctx context.Context,
//	userID, attemptID int64,
//) (*domain.AttemptResult, error) {
//
//	// 1. Загружаем попытку
//	attempt, err := uc.attemptRepo.GetAttemptByID(ctx, attemptID)
//	if err != nil {
//		return nil, err
//	}
//	if attempt == nil {
//		return nil, ErrAttemptNotFound
//	}
//
//	// 2. Проверяем что попытка принадлежит этому пользователю
//	if attempt.UserID != userID {
//		return nil, ErrNotYourAttempt
//	}
//
//	// 3. Проверяем что попытка ещё не завершена
//	if attempt.Status != domain.AttemptStatusInProgress {
//		return nil, ErrAlreadyFinished
//	}
//
//	// 4. Загружаем тест — нужен show_answer_mode
//	test, err := uc.testRepo.GetTestByID(ctx, *attempt.TestID)
//	if err != nil {
//		return nil, err
//	}
//
//	// 5. Загружаем все вопросы теста
//	questions, err := uc.questionRepo.GetQuestionsByTestID(ctx, *attempt.TestID)
//	if err != nil {
//		return nil, err
//	}
//
//	// 6. Загружаем все ответы студента по этой попытке
//	userAnswers, err := uc.userAnswerRepo.GetByAttemptID(ctx, attemptID)
//	if err != nil {
//		return nil, err
//	}
//
//	// 7. Если режим end_only — нужно проверить правильность всех ответов
//	// (при immediate это уже было сделано в SubmitAnswer)
//	if test.ShowAnswerMode == domain.ShowAnswerModeEndOnly {
//		userAnswers, err = uc.checkAndUpdateAllAnswers(ctx, userAnswers, questions)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	// 8. Считаем score
//	score, maxScore := uc.calculateScore(userAnswers, questions)
//
//	// 9. Считаем процент
//	var percentage float64
//	if maxScore > 0 {
//		percentage = float64(score) / float64(maxScore) * 100
//	}
//
//	// 10. Обновляем попытку
//	now := time.Now()
//	attempt.Status = domain.AttemptStatusCompleted
//	attempt.FinishedAt = &now
//	attempt.Score = score
//	attempt.MaxScore = maxScore
//	attempt.Percentage = percentage
//
//	if err := uc.attemptRepo.UpdateAttempt(ctx, attempt); err != nil {
//		return nil, err
//	}
//
//	// 11. Собираем детальный разбор
//	result, err := uc.buildAttemptResult(ctx, attempt, questions, userAnswers)
//	if err != nil {
//		return nil, err
//	}
//
//	return result, nil
//}
//
//// =============================================================================
//// GetResult — получить результат уже завершённой попытки
//// Используется когда студент хочет пересмотреть разбор из истории
//// =============================================================================
//
//func (uc *testAttemptUseCase) GetResult(
//	ctx context.Context,
//	userID, attemptID int64,
//) (*domain.AttemptResult, error) {
//
//	// 1. Загружаем попытку
//	attempt, err := uc.attemptRepo.GetAttemptByID(ctx, attemptID)
//	if err != nil {
//		return nil, err
//	}
//	if attempt == nil {
//		return nil, ErrAttemptNotFound
//	}
//
//	// 2. Проверяем владельца
//	if attempt.UserID != userID {
//		return nil, ErrNotYourAttempt
//	}
//
//	// 3. Попытка должна быть завершена
//	if attempt.Status != domain.AttemptStatusCompleted {
//		return nil, errors.New("attempt is not completed yet")
//	}
//
//	// 4. Загружаем вопросы и ответы студента
//	questions, err := uc.questionRepo.GetQuestionsByTestID(ctx, *attempt.TestID)
//	if err != nil {
//		return nil, err
//	}
//
//	userAnswers, err := uc.userAnswerRepo.GetByAttemptID(ctx, attemptID)
//	if err != nil {
//		return nil, err
//	}
//
//	// 5. Собираем разбор
//	return uc.buildAttemptResult(ctx, attempt, questions, userAnswers)
//}
//
//// =============================================================================
//// Вспомогательные методы (приватные)
//// =============================================================================
//
//// getOrCreateAttempt — возвращает активную попытку или создаёт новую
//func (uc *testAttemptUseCase) getOrCreateAttempt(
//	ctx context.Context,
//	userID, testID, subjectID int64,
//) (*domain.Attempt, error) {
//
//	// ищем существующую незавершённую попытку
//	attempt, err := uc.attemptRepo.GetInProgressAttemptByTestID(ctx, userID, testID)
//	if err != nil {
//		return nil, err
//	}
//
//	// нашли — возвращаем её
//	if attempt != nil {
//		return attempt, nil
//	}
//
//	// не нашли — создаём новую
//	newAttempt := &domain.Attempt{
//		UserID:    userID,
//		SubjectID: subjectID,
//		TestID:    &testID,
//		Mode:      domain.AttemptModeTest,
//		Status:    domain.AttemptStatusInProgress,
//	}
//
//	if err := uc.attemptRepo.CreateAttempt(ctx, newAttempt); err != nil {
//		return nil, err
//	}
//
//	return newAttempt, nil
//}
//
//// checkAnswer — проверяет правильность ответа в зависимости от типа вопроса
//func (uc *testAttemptUseCase) checkAnswer(
//	ctx context.Context,
//	question *domain.Question,
//	input domain.SubmitAnswerInput,
//) (bool, error) {
//
//	switch question.QuestionType {
//
//	case domain.QuestionTypeSingleChoice:
//		// студент должен передать один AnswerID
//		if input.AnswerID == nil {
//			return false, nil
//		}
//		answer, err := uc.answerRepo.GetByID(ctx, *input.AnswerID)
//		if err != nil {
//			return false, err
//		}
//		if answer == nil {
//			return false, nil
//		}
//		// проверяем что этот ответ относится к этому вопросу (защита от подмены)
//		if answer.QuestionID != question.ID {
//			return false, ErrWrongQuestion
//		}
//		return answer.IsCorrect, nil
//
//	case domain.QuestionTypeMultipleChoice:
//		// загружаем все правильные ответы для этого вопроса
//		allAnswers, err := uc.answerRepo.GetByQuestionID(ctx, question.ID)
//		if err != nil {
//			return false, err
//		}
//
//		// собираем множество правильных ID
//		correctIDs := make(map[int64]bool)
//		for _, a := range allAnswers {
//			if a.IsCorrect {
//				correctIDs[a.ID] = true
//			}
//		}
//
//		// собираем множество выбранных студентом ID
//		selectedIDs := make(map[int64]bool)
//		for _, id := range input.AnswerIDs {
//			selectedIDs[id] = true
//		}
//
//		// строгая проверка: множества должны совпадать полностью
//		if len(selectedIDs) != len(correctIDs) {
//			return false, nil
//		}
//		for id := range correctIDs {
//			if !selectedIDs[id] {
//				return false, nil
//			}
//		}
//		return true, nil
//
//	case domain.QuestionTypeShortAnswer:
//		// загружаем эталонный ответ (первый is_correct=true)
//		allAnswers, err := uc.answerRepo.GetByQuestionID(ctx, question.ID)
//		if err != nil {
//			return false, err
//		}
//		for _, a := range allAnswers {
//			if a.IsCorrect {
//				// сравниваем без учёта регистра и лишних пробелов
//				expected := strings.TrimSpace(strings.ToLower(a.AnswerText))
//				got := strings.TrimSpace(strings.ToLower(input.TextAnswer))
//				return expected == got, nil
//			}
//		}
//		return false, nil
//	}
//
//	return false, nil
//}
//
//// saveUserAnswer — создаёт или обновляет ответ студента
//func (uc *testAttemptUseCase) saveUserAnswer(
//	ctx context.Context,
//	attemptID, questionID int64,
//	input domain.SubmitAnswerInput,
//	isCorrect bool,
//) error {
//
//	// проверяем — отвечал ли уже на этот вопрос
//	existing, err := uc.userAnswerRepo.GetByAttemptAndQuestion(ctx, attemptID, questionID)
//	if err != nil {
//		return err
//	}
//
//	if existing == nil {
//		// первый ответ — создаём
//		ua := &domain.UserAnswer{
//			AttemptID:  attemptID,
//			QuestionID: questionID,
//			AnswerID:   input.AnswerID,
//			AnswerIDs:  input.AnswerIDs,
//			TextAnswer: input.TextAnswer,
//			IsCorrect:  &isCorrect,
//		}
//		return uc.userAnswerRepo.Create(ctx, ua)
//	}
//
//	// студент вернулся назад и поменял ответ — обновляем
//	existing.AnswerID = input.AnswerID
//	existing.AnswerIDs = input.AnswerIDs
//	existing.TextAnswer = input.TextAnswer
//	existing.IsCorrect = &isCorrect
//	return uc.userAnswerRepo.Update(ctx, existing)
//}
//
//// checkAndUpdateAllAnswers — для режима end_only
//// проверяем все ответы студента и обновляем is_correct в БД
//func (uc *testAttemptUseCase) checkAndUpdateAllAnswers(
//	ctx context.Context,
//	userAnswers []*domain.UserAnswer,
//	questions []*domain.Question,
//) ([]*domain.UserAnswer, error) {
//
//	// строим map вопросов для быстрого доступа по ID
//	questionMap := make(map[int64]*domain.Question, len(questions))
//	for _, q := range questions {
//		questionMap[q.ID] = q
//	}
//
//	for _, ua := range userAnswers {
//		question, ok := questionMap[ua.QuestionID]
//		if !ok {
//			continue
//		}
//
//		// формируем input из сохранённых данных
//		input := domain.SubmitAnswerInput{
//			QuestionID: ua.QuestionID,
//			AnswerID:   ua.AnswerID,
//			AnswerIDs:  ua.AnswerIDs,
//			TextAnswer: ua.TextAnswer,
//		}
//
//		isCorrect, err := uc.checkAnswer(ctx, question, input)
//		if err != nil {
//			return nil, err
//		}
//
//		ua.IsCorrect = &isCorrect
//
//		// обновляем в БД
//		if err := uc.userAnswerRepo.Update(ctx, ua); err != nil {
//			return nil, err
//		}
//	}
//
//	return userAnswers, nil
//}
//
//// calculateScore — считает набранные и максимальные баллы
//func (uc *testAttemptUseCase) calculateScore(
//	userAnswers []*domain.UserAnswer,
//	questions []*domain.Question,
//) (score, maxScore int) {
//
//	// строим map вопросов чтобы знать points каждого
//	questionMap := make(map[int64]*domain.Question, len(questions))
//	for _, q := range questions {
//		questionMap[q.ID] = q
//		maxScore += q.Points
//	}
//
//	for _, ua := range userAnswers {
//		if ua.IsCorrect != nil && *ua.IsCorrect {
//			if q, ok := questionMap[ua.QuestionID]; ok {
//				score += q.Points
//			}
//		}
//	}
//
//	return score, maxScore
//}
//
//// buildAttemptResult — собирает детальный разбор попытки
//func (uc *testAttemptUseCase) buildAttemptResult(
//	ctx context.Context,
//	attempt *domain.Attempt,
//	questions []*domain.Question,
//	userAnswers []*domain.UserAnswer,
//) (*domain.AttemptResult, error) {
//
//	// строим map ответов студента по question_id для быстрого доступа
//	userAnswerMap := make(map[int64]*domain.UserAnswer, len(userAnswers))
//	for _, ua := range userAnswers {
//		userAnswerMap[ua.QuestionID] = ua
//	}
//
//	details := make([]*domain.AttemptAnswerDetail, 0, len(questions))
//
//	for _, q := range questions {
//		detail := &domain.AttemptAnswerDetail{
//			Question:   q,
//			UserAnswer: userAnswerMap[q.ID], // nil если студент не ответил на вопрос
//		}
//
//		// загружаем правильный ответ для показа в разборе
//		// для single_choice — один правильный вариант
//		// для multiple_choice и short_answer — тоже нужен для показа
//		allAnswers, err := uc.answerRepo.GetByQuestionID(ctx, q.ID)
//		if err != nil {
//			return nil, err
//		}
//		for _, a := range allAnswers {
//			if a.IsCorrect {
//				detail.CorrectAnswer = a
//				break
//			}
//		}
//
//		details = append(details, detail)
//	}
//
//	return &domain.AttemptResult{
//		Attempt: attempt,
//		Answers: details,
//	}, nil
//}
