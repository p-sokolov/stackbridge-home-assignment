package subscription

import (
	"context"

	"stackbridge-home-task/internal/models"

	"github.com/google/uuid"
)

func (s *service) Update(ctx context.Context, id uuid.UUID, input *models.SubscriptionDTO) (*models.Subscription, error) {
	sub, err := buildSubscription(input)
	if err != nil {
		return nil, err
	}

	sub.ID = id

	return s.repo.Update(ctx, sub)
}
