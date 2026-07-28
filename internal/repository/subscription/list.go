package subscription

import (
	"context"

	"stackbridge-home-task/internal/db/sqlc/storage"
	"stackbridge-home-task/internal/models"
	"stackbridge-home-task/internal/repository"

	"github.com/google/uuid"
)

func (r *repo) List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Subscription, error) {
	q := repository.Queries(ctx, r.queries)

	subs, err := q.ListSubscriptions(ctx, storage.ListSubscriptionsParams{
		UserID:    userID,
		LimitCnt:  int32(limit),
		OffsetCnt: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*models.Subscription, 0, len(subs))
	for _, sub := range subs {
		result = append(result, toModelSubscription(sub))
	}

	return result, nil
}
