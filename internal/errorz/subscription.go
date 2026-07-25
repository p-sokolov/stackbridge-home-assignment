package errorz

import "errors"

var (
	ErrInvalidSubscriptionData 	  = errors.New("invalid subscription data")
	ErrSubscriptionAlreadyExists  = errors.New("subscription already exists")
	ErrSubscriptionNotFound       = errors.New("subscription not found")
)
