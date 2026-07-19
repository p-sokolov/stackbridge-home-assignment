package repository

import (
	"context"

	"github.com/google/uuid"
	"stackbridge-home-task/intenal/domain"
)

type Repository interface {
    Create(ctx context.Context, s *domain.Subscription) error
    Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
    Update(ctx context.Context, id uuid.UUID) (error)
    Delete(ctx context.Context, id uuid.UUID) (error)
    List(ctx context.Context, id uuid.UUID) (*[]domain.Subscription, error)
    Sum(ctx context.Context, id uuid.UUID) (int, error)
}