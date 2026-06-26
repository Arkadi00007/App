// internal/domain/subscription.go
package domain

import (
	"context"
	"time"
)

type UserSubscription struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PlanID    int64     `json:"plan_id"`
	Status    string    `json:"status"`
	StartedAt time.Time `json:"started_at"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type SubscriptionRepository interface {
	// единственный метод который нужен для проверки доступа
	HasActiveSubscription(ctx context.Context, userID int64) (bool, error)
}
