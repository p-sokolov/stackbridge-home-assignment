package subscription

import (
	"context"

	"stackbridge-home-task/internal/models"
	"stackbridge-home-task/internal/repository"

	"github.com/google/uuid"
)

func (r *repo) List(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error) {
	q := repository.Queries(ctx, r.queries)

	subs, err := q.ListSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Subscription, 0, len(subs))
	for _, sub := range subs {
		result = append(result, toModelSubscription(sub))
	}

	return result, nil
}