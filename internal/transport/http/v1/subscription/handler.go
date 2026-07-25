package subscription

import (
	"context"
	"time"

	"stackbridge-home-task/internal/models"
	v1 "stackbridge-home-task/internal/transport/http/v1"

	"github.com/google/uuid"
)

const monthLayout = "01-2006"

type service interface {
	Create(ctx context.Context, input *models.SubscriptionDTO) (*models.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	List(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, input *models.SubscriptionDTO) (*models.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CalculateTotalCost(ctx context.Context, userID uuid.UUID, serviceName, startPeriod, endPeriod string) (int, error)
}

type Handler struct {
	svc service
}

func New(svc service) *Handler {
	return &Handler{svc: svc}
}

func toV1Subscription(sub *models.Subscription) v1.Subscription {
	return v1.Subscription{
		Id:          sub.ID,
		UserId:      sub.UserID,
		Price:       sub.Price,
		ServiceName: sub.ServiceName,
		StartDate:   formatMonth(sub.StartDate),
		EndDate:     formatOptionalMonth(sub.EndDate),
	}
}

func toV1Subscriptions(subs []*models.Subscription) []v1.Subscription {
	result := make([]v1.Subscription, 0, len(subs))
	for _, sub := range subs {
		result = append(result, toV1Subscription(sub))
	}
	return result
}

func formatMonth(t time.Time) string {
	return t.Format(monthLayout)
}

func formatOptionalMonth(t *time.Time) *string {
	if t == nil {
		return nil
	}

	value := t.Format(monthLayout)
	return &value
}

func badRequest(message string) v1.BadRequestJSONResponse {
	return v1.BadRequestJSONResponse{Message: message}
}

func notFound(message string) v1.NotFoundJSONResponse {
	return v1.NotFoundJSONResponse{Message: message}
}