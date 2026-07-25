package subscription

import (
	"context"

	"stackbridge-home-task/internal/models"
)

func (s *service) Create(ctx context.Context, input *models.SubscriptionDTO) (*models.Subscription, error) {
	sub, err := buildSubscription(input)
	if err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, sub)
}
