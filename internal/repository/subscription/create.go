package subscription

import (
	"context"
	
	"stackbridge-home-task/internal/db/sqlc/storage"
	"stackbridge-home-task/internal/errorz"
	"stackbridge-home-task/internal/models"
	"stackbridge-home-task/internal/repository"
	"stackbridge-home-task/pkg/dberrors"
)

func (r *repo) Create(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	args := storage.CreateSubscriptionParams{
		ServiceName: sub.ServiceName,
		Price: int32(sub.Price),
		UserID: sub.UserID,
		StartDate: toPgDate(sub.StartDate),
		EndDate: sub.EndDate,
	}
		
	q := repository.Queries(ctx, r.queries)

	sqlcSub, err := q.CreateSubscription(ctx, args)
	if err != nil {
		if dberrors.IsUniqueViolation(err) {
			return nil, errorz.ErrSubscriptionAlreadyExists
		}
		return nil, err
	}

	return toModelSubscription(sqlcSub), nil
}
