package subscription

import (
	"errors"

	"stackbridge-home-task/internal/errorz"
)

func isBadRequest(err error) bool {
	return errors.Is(err, errorz.ErrInvalidSubscriptionData) ||
		errors.Is(err, errorz.ErrSubscriptionAlreadyExists)
}

func isNotFound(err error) bool {
	return errors.Is(err, errorz.ErrSubscriptionNotFound)
}
