-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_user_service_dates
    ON subscriptions(user_id, service_name, start_date, end_date);    
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_subscriptions_user_service_dates;
DROP INDEX IF EXISTS idx_subscriptions_user_id;
-- +goose StatementEnd
