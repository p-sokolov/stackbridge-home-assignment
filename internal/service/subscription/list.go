package subscription

import (
	"context"

	"stackbridge-home-task/internal/errorz"
	"stackbridge-home-task/internal/models"

	"github.com/google/uuid"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

func (s *service) List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Subscription, error) {
	if userID == uuid.Nil || offset < 0 {
		return nil, errorz.ErrInvalidSubscriptionData
	}

	if limit == 0 {
		limit = defaultLimit
	}
	if limit < 0 || limit > maxLimit {
		return nil, errorz.ErrInvalidSubscriptionData
	}

	return s.repo.List(ctx, userID, limit, offset)
}
