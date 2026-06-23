package domain

import (
	"context"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	IsVerified   bool      `json:"is_verified"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// что принимаем при регистрации
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// что принимаем при входе
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// что возвращаем после успешного входа/регистрации
type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
}

type UserUseCase interface {
	Register(ctx context.Context, email, password, name string) error
	VerifyEmail(ctx context.Context, email, code string) (*User, string, error)
	Login(ctx context.Context, email, password string) (*User, string, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email, code, newPassword string) error
	ResendCode(ctx context.Context, emailAddr, codeType string) error
}
