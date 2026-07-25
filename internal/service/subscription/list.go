package subscription

import (
	"context"

	"stackbridge-home-task/internal/models"

	"github.com/google/uuid"
)

func (s *service) List(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error) {
	return s.repo.List(ctx, userID)
}