// internal/repository/postgres/subscription.go
package postgres

import (
	"Test_App/internal/domain"
	"context"
	"database/sql"
	"time"
)

type subscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) domain.SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) HasActiveSubscription(ctx context.Context, userID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
        SELECT EXISTS (
            SELECT 1 FROM user_subscriptions
            WHERE user_id = $1
              AND status = 'active'
              AND expires_at > $2
        )
    `, userID, time.Now()).Scan(&exists)

	if err != nil {
		return false, err
	}
	return exists, nil
}
