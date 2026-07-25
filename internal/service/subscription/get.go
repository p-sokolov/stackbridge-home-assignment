package subscription

import (
	"context"

	"stackbridge-home-task/internal/models"

	"github.com/google/uuid"
)

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}