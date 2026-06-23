package handler

import (
	"Test_App/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserAuthHandler содержит usecase авторизации
// и предоставляет HTTP-методы для маршрутов /auth/*
type UserAuthHandler struct {
	uc domain.UserUseCase
}

func NewUserAuthHandler(uc domain.UserUseCase) *UserAuthHandler {
	return &UserAuthHandler{uc: uc}
}

// ===================== структуры запросов =====================

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type verifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type resendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Type  string `json:"type" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type resetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ===================== структура ответа =====================

type authResponse struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

// ===================== хэндлеры =====================

// Register godoc
// POST /auth/register
func (h *UserAuthHandler) Register(c *gin.Context) {
	var req registerRequest

	// ShouldBindJSON парсит JSON из тела запроса в req
	// и проверяет binding-теги (required, email, min=6 и т.д.)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// c.Request.Context() — тот же context.Context из net/http
	if err := h.uc.Register(c.Request.Context(), req.Email, req.Password, req.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "код подтверждения отправлен на email"})
}

// VerifyEmail godoc
// POST /auth/verify-email
func (h *UserAuthHandler) VerifyEmail(c *gin.Context) {
	var req verifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.uc.VerifyEmail(c.Request.Context(), req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authResponse{Token: token, User: user})
}

// ResendCode godoc
// POST /auth/resend-code
func (h *UserAuthHandler) ResendCode(c *gin.Context) {
	var req resendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверяем что type — одно из допустимых значений
	if req.Type != domain.CodeTypeEmailVerify && req.Type != domain.CodeTypeResetPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный тип кода"})
		return
	}

	if err := h.uc.ResendCode(c.Request.Context(), req.Email, req.Type); err != nil {
		// 429 Too Many Requests — подходящий статус для cooldown
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "код отправлен повторно"})
}

// Login godoc
// POST /auth/login
func (h *UserAuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.uc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authResponse{Token: token, User: user})
}

// ForgotPassword godoc
// POST /auth/forgot-password
func (h *UserAuthHandler) ForgotPassword(c *gin.Context) {
	var req forgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.ForgotPassword(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "если email существует, код отправлен"})
}

// ResetPassword godoc
// POST /auth/reset-password
func (h *UserAuthHandler) ResetPassword(c *gin.Context) {
	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.ResetPassword(c.Request.Context(), req.Email, req.Code, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "пароль успешно изменён"})
}
