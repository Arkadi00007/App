package domain

import (
	"context"
	"time"
)

type VerificationCode struct {
	ID        int64
	UserID    int64
	Code      string
	Type      string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

const (
	CodeTypeEmailVerify   = "email_verify"
	CodeTypeResetPassword = "reset_password"

	ResendCooldown = 60 * time.Second
)

type VerificationCodeRepository interface {
	CreateCode(ctx context.Context, vc *VerificationCode) error
	GetActiveCode(ctx context.Context, userID int64, codeType string) (*VerificationCode, error)
	GetLatestCode(ctx context.Context, userID int64, codeType string) (*VerificationCode, error)
	MarkAsUsedCode(ctx context.Context, id int64) error
}
