package subscription

import (
	"context"
	"errors"

	"stackbridge-home-task/internal/errorz"
	"stackbridge-home-task/internal/models"
	"stackbridge-home-task/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	q := repository.Queries(ctx, r.queries)

	sub, err := q.GetSubscriptionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorz.ErrSubscriptionNotFound
		}
		return nil, err
	}

	return toModelSubscription(sub), nil
}
