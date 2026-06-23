CREATE TABLE subscription_plans (
                                    id            BIGSERIAL PRIMARY KEY,
                                    name          VARCHAR(100) NOT NULL,
                                    price         NUMERIC(10,2) NOT NULL,
                                    currency      VARCHAR(3) DEFAULT 'AMD',
                                    duration_days INTEGER NOT NULL,
                                    is_active     BOOLEAN DEFAULT TRUE
);

CREATE TABLE user_subscriptions (
                                    id         BIGSERIAL PRIMARY KEY,
                                    user_id    BIGINT NOT NULL REFERENCES users(id),
                                    plan_id    BIGINT NOT NULL REFERENCES subscription_plans(id),
                                    status     VARCHAR(20) NOT NULL CHECK (status IN ('active', 'expired', 'cancelled')),
                                    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                    expires_at TIMESTAMP NOT NULL,
                                    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE payments (
                          id                  BIGSERIAL PRIMARY KEY,
                          user_id             BIGINT NOT NULL REFERENCES users(id),
                          plan_id             BIGINT NOT NULL REFERENCES subscription_plans(id),
                          amount              NUMERIC(10,2) NOT NULL,
                          currency            VARCHAR(3) DEFAULT 'AMD',
                          status              VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'success', 'failed')),
                          provider            VARCHAR(30) NOT NULL,
                          provider_payment_id VARCHAR(255),
                          created_at          TIMESTAMP DEFAULT NOW(),
                          updated_at          TIMESTAMP DEFAULT NOW()
);