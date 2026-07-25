package subscription

import (
	"context"
	"errors"

	"stackbridge-home-task/internal/db/sqlc/storage"
	"stackbridge-home-task/internal/errorz"
	"stackbridge-home-task/internal/models"
	"stackbridge-home-task/internal/repository"

	"github.com/jackc/pgx/v5"
)

func (r *repo) Update(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	q := repository.Queries(ctx, r.queries)

	updated, err := q.UpdateSubscription(ctx, storage.UpdateSubscriptionParams{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       int32(sub.Price),
		UserID:      sub.UserID,
		StartDate:   toPgDate(sub.StartDate),
		EndDate:     sub.EndDate,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorz.ErrSubscriptionNotFound
		}
		return nil, err
	}

	return toModelSubscription(updated), nil
}
