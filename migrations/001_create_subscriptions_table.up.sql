-- migrations/001_create_subscriptions_table.up.sql
CREATE TABLE subscriptions (
    id BIGSERIAL PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    price BIGINT NOT NULL CHECK (price >= 0),
    user_id UUID NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    
    -- CONSTRAINT valid_dates CHECK (end_date > start_date),
    CONSTRAINT unique_subscription UNIQUE (user_id, service_name, start_date)
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_service_name ON subscriptions(service_name);
CREATE INDEX idx_subscriptions_dates ON subscriptions(start_date);
CREATE INDEX idx_subscriptions_created_at ON subscriptions(created_at);
