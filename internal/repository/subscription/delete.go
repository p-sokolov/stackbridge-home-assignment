package subscription

import (
	"context"

	"stackbridge-home-task/internal/errorz"
	"stackbridge-home-task/internal/repository"

	"github.com/google/uuid"
)

func (r *repo) Delete(ctx context.Context, id uuid.UUID) error {
	q := repository.Queries(ctx, r.queries)

	rowsAffected, err := q.DeleteSubscription(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errorz.ErrSubscriptionNotFound
	}

	return nil
}
