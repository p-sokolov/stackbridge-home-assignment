package subscription

import (
	"context"
	"time"

	"stackbridge-home-task/internal/db/sqlc/storage"
	"stackbridge-home-task/internal/repository"

	"github.com/google/uuid"
)

func (r *repo) CalculateTotalCost(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	startPeriod time.Time,
	endPeriod time.Time,
) (int, error) {
	q := repository.Queries(ctx, r.queries)

	total, err := q.CalculateSubscriptionsTotal(ctx, storage.CalculateSubscriptionsTotalParams{
		UserID:      userID,
		ServiceName: serviceName,
		StartPeriod: toPgDate(endPeriod),
		EndPeriod:   toPgDate(endPeriod),
	})
	if err != nil {
		return 0, err
	}

	return int(total), nil
}
