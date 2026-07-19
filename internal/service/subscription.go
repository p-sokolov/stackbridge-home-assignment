package service

import (
	"context"
	"github.com/google/uuid"
	"stackbridge-home-task/internal/domain"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, sub *domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, userID uuid.UUID) ([]*domain.Subscription, error)
	CalculateTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, startPeriod, endPeriod time.Time) (int, error)
}

type SubscriptionService interface {
}