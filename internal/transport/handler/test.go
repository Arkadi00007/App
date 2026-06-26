package handler

import (
	"Test_App/internal/domain"
	"Test_App/internal/usecase"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	testUC     domain.TestAttemptUseCase
	testListUC domain.TestListUseCase
}

func NewTestHandler(
	testUC domain.TestAttemptUseCase,
	testListUC domain.TestListUseCase,
) *TestHandler {
	return &TestHandler{
		testUC:     testUC,
		testListUC: testListUC,
	}
}

// ===================== структуры запросов =====================

// submitAnswerRequest — тело запроса при ответе на вопрос
// Заполняется только одно из трёх полей в зависимости от типа вопроса:
// AnswerID   — для single_choice
// AnswerIDs  — для multiple_choice
// TextAnswer — для short_answer
type submitAnswerRequest struct {
	QuestionID int64   `json:"question_id" binding:"required"`
	AnswerID   *int64  `json:"answer_id"`
	AnswerIDs  []int64 `json:"answer_ids"`
	TextAnswer string  `json:"text_answer"`
}

// ===================== вспомогательная функция =====================

// getUserID — достаёт userID из контекста Gin (туда его кладёт AuthMiddleware)
func getUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	userID, ok := val.(int64)
	return userID, ok
}

// ===================== хэндлеры =====================

// GetTestsBySection godoc
// GET /sections/:id/tests
// Возвращает список тестов внутри раздела
func (h *TestHandler) GetTestsBySection(c *gin.Context) {
	sectionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный id раздела"})
		return
	}

	tests, err := h.testListUC.GetTestsBySectionID(c.Request.Context(), sectionID)
	if err != nil {
		if errors.Is(err, usecase.ErrSectionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения тестов"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tests})
}

// GetTest godoc
// GET /tests/:id
// Открывает конкретный тест:
// - возвращает сам тест с его настройками (show_answer_mode, is_premium)
// - возвращает все вопросы теста в правильном порядке
// - возвращает варианты ответов для всех вопросов (без is_correct)
// - возвращает текущую незавершённую попытку если она есть (nil если студент заходит впервые)
func (h *TestHandler) GetTest(c *gin.Context) {
	testID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный id теста"})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	test, questions, answers, attempt, err := h.testUC.GetTestForUser(
		c.Request.Context(),
		userID,
		testID,
	)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrTestNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrPremiumRequired):
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()}) // 402
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения теста"})
		}
		return
	}
	//if err != nil {
	//	if errors.Is(err, usecase.ErrTestNotFound) {
	//		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	//		return
	//	}
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения теста"})
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{
		"test":      test,
		"questions": questions,
		"answers":   answers,
		"attempt":   attempt,
	})
}

// SubmitAnswer godoc
// POST /tests/:id/answer
// Студент ответил на вопрос.
// Если попытки ещё нет — создаётся автоматически (при первом ответе).
// Если студент вернулся назад и меняет ответ — обновляется существующий.
// Возвращает:
//   - attempt_id (нужен фронтенду для последующих запросов finish/result)
//   - is_correct и explanation (только если show_answer_mode == "immediate")
func (h *TestHandler) SubmitAnswer(c *gin.Context) {
	testID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный id теста"})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	var req submitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := domain.SubmitAnswerInput{
		QuestionID: req.QuestionID,
		AnswerID:   req.AnswerID,
		AnswerIDs:  req.AnswerIDs,
		TextAnswer: req.TextAnswer,
	}

	result, err := h.testUC.SubmitAnswer(c.Request.Context(), userID, testID, input)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrTestNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrWrongQuestion):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сохранения ответа"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// Restart godoc
// POST /tests/:id/restart
// Студент нажал "начать заново".
// Текущая незавершённая попытка помечается как abandoned.
// Создаётся новая пустая попытка.
func (h *TestHandler) Restart(c *gin.Context) {
	testID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный id теста"})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	attempt, err := h.testUC.Restart(c.Request.Context(), userID, testID)
	if err != nil {
		if errors.Is(err, usecase.ErrTestNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка перезапуска теста"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": attempt})
}

// Finish godoc
// POST /attempts/:id/finish
// Студент нажал "завершить тест".
// Для режима end_only — здесь происходит проверка всех ответов.
// Для режима immediate — просто подсчитывается итоговый score.
// Возвращает полный детальный разбор попытки.
func (h *TestHandler) Finish(c *gin.Context) {
	attemptID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный id попытки"})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	result, err := h.testUC.Finish(c.Request.Context(), userID, attemptID)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrAttemptNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrNotYourAttempt):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrAlreadyFinished):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка завершения теста"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetResult godoc
// GET /attempts/:id/result
// Получить результат уже завершённой попытки.
// Используется когда студент хочет пересмотреть разбор из истории попыток.
func (h *TestHandler) GetResult(c *gin.Context) {
	attemptID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный id попытки"})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	result, err := h.testUC.GetResult(c.Request.Context(), userID, attemptID)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrAttemptNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrNotYourAttempt):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения результата"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
