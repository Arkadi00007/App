package postgres

import (
	"Test_App/internal/domain"
	"context"
	"database/sql"
	"time"
)

type verificationCodeRepository struct {
	db *sql.DB
}

func NewVerificationCodeRepository(db *sql.DB) domain.VerificationCodeRepository {
	return &verificationCodeRepository{db: db}
}

func (r *verificationCodeRepository) CreateCode(ctx context.Context, vc *domain.VerificationCode) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO verification_codes (user_id, code, type, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`, vc.UserID, vc.Code, vc.Type, vc.ExpiresAt).
		Scan(&vc.ID, &vc.CreatedAt)
}

func (r *verificationCodeRepository) GetActiveCode(ctx context.Context, userID int64, codeType string) (*domain.VerificationCode, error) {
	vc := &domain.VerificationCode{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, code, type, expires_at, created_at
		FROM verification_codes
		WHERE user_id = $1
		  AND type = $2
		  AND used_at IS NULL
		  AND expires_at > $3
		ORDER BY created_at DESC
		LIMIT 1
	`, userID, codeType, time.Now()).Scan(
		&vc.ID,
		&vc.UserID,
		&vc.Code,
		&vc.Type,
		&vc.ExpiresAt,
		&vc.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return vc, nil
}

func (r *verificationCodeRepository) GetLatestCode(ctx context.Context, userID int64, codeType string) (*domain.VerificationCode, error) {
	vc := &domain.VerificationCode{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, code, type, expires_at, created_at
		FROM verification_codes
		WHERE user_id = $1 AND type = $2
		ORDER BY created_at DESC
		LIMIT 1
	`, userID, codeType).Scan(
		&vc.ID, &vc.UserID, &vc.Code, &vc.Type, &vc.ExpiresAt, &vc.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return vc, nil
}

func (r *verificationCodeRepository) MarkAsUsedCode(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE verification_codes
		SET used_at = $1
		WHERE id = $2
	`, time.Now(), id)
	return err
}
