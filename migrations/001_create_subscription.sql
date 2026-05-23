-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscription (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_subscription_user_id ON subscription(user_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_subscription_start_date ON subscription(start_date);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_subscription_end_date ON subscription(end_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscription;
-- +goose StatementEnd