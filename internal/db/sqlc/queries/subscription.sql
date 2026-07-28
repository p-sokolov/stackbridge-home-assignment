-- name: CreateSubscription :one
INSERT INTO subscriptions (
    service_name,
    price,
    user_id,
    start_date,
    end_date
)
VALUES (
    sqlc.arg(service_name),
    sqlc.arg(price),
    sqlc.arg(user_id),
    sqlc.arg(start_date),
    sqlc.arg(end_date)
)
RETURNING *;


-- name: GetSubscriptionByID :one
SELECT *
FROM subscriptions
WHERE id = sqlc.arg(id);


-- name: ListSubscriptions :many
SELECT *
FROM subscriptions
WHERE user_id = sqlc.arg(user_id)
ORDER BY start_date DESC, service_name ASC
LIMIT sqlc.arg(limit_cnt)
OFFSET sqlc.arg(offset_cnt);


-- name: UpdateSubscription :one
UPDATE subscriptions
SET
    service_name = sqlc.arg(service_name),
    price = sqlc.arg(price),
    user_id = sqlc.arg(user_id),
    start_date = sqlc.arg(start_date),
    end_date = sqlc.arg(end_date)
WHERE id = sqlc.arg(id)
RETURNING *;


-- name: DeleteSubscription :execrows
DELETE
FROM subscriptions
WHERE id = sqlc.arg(id);


-- name: CalculateSubscriptionsTotal :one
SELECT COALESCE(SUM(
    price * (
        (
            EXTRACT(YEAR FROM AGE(
                LEAST(COALESCE(end_date, sqlc.arg(end_period)::date), sqlc.arg(end_period)::date),
                GREATEST(start_date, sqlc.arg(start_period)::date)
            ))::int * 12
        )
        + EXTRACT(MONTH FROM AGE(
            LEAST(COALESCE(end_date, sqlc.arg(end_period)::date), sqlc.arg(end_period)::date),
            GREATEST(start_date, sqlc.arg(start_period)::date)
        ))::int
        + 1
    )
), 0)::int
FROM subscriptions
WHERE
    user_id = sqlc.arg(user_id)
    AND service_name = sqlc.arg(service_name)
    AND start_date <= sqlc.arg(end_period)::date
    AND (end_date IS NULL OR end_date >= sqlc.arg(start_period)::date);