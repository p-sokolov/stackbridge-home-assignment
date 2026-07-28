package subscription

import (
	"context"
	"strings"
	"time"

	"stackbridge-home-task/internal/errorz"
	"stackbridge-home-task/internal/models"

	"github.com/google/uuid"
)

const monthLayout = "01-2006"

type repo interface {
	Create(ctx context.Context, sub *models.Subscription) (*models.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) (*models.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CalculateTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, startPeriod, endPeriod time.Time) (int, error)
}

type service struct {
	repo repo
}

func New(repo repo) *service {
	return &service{repo: repo}
}

func parseMonth(value string) (time.Time, error) {
	t, err := time.Parse(monthLayout, value)
	if err != nil {
		return time.Time{}, errorz.ErrInvalidSubscriptionData
	}

	return t, nil
}

func buildSubscription(input *models.SubscriptionDTO) (*models.Subscription, error) {
	serviceName := strings.TrimSpace(input.ServiceName)
	if serviceName == "" || input.Price <= 0 || input.UserID == uuid.Nil {
		return nil, errorz.ErrInvalidSubscriptionData
	}

	startDate, err := parseMonth(input.StartDate)
	if err != nil {
		return nil, err
	}

	var endDate *time.Time
	if input.EndDate != nil && strings.TrimSpace(*input.EndDate) != "" {
		parsedEndDate, err := parseMonth(*input.EndDate)
		if err != nil {
			return nil, err
		}
		if parsedEndDate.Before(startDate) {
			return nil, errorz.ErrInvalidSubscriptionData
		}
		endDate = &parsedEndDate
	}

	return &models.Subscription{
		ServiceName: serviceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}
