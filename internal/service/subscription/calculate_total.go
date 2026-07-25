package subscription

import (
	"context"
	"strings"

	"stackbridge-home-task/internal/errorz"

	"github.com/google/uuid"
)

func (s *service) CalculateTotalCost(ctx context.Context, userID uuid.UUID, serviceName, startPeriod, endPeriod string) (int, error) {
	serviceName = strings.TrimSpace(serviceName)
	if userID == uuid.Nil || serviceName == "" {
		return 0, errorz.ErrInvalidSubscriptionData
	}

	startDate, err := parseMonth(startPeriod)
	if err != nil {
		return 0, err
	}

	endDate, err := parseMonth(endPeriod)
	if err != nil {
		return 0, err
	}

	if endDate.Before(startDate) {
		return 0, errorz.ErrInvalidSubscriptionData
	}

	return s.repo.CalculateTotalCost(ctx, userID, serviceName, startDate, endDate)
}
