package dto

import "Test_App/internal/domain"

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// что возвращаем после успешного входа/регистрации
type AuthResponse struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}
