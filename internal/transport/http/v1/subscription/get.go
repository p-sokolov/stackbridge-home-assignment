package subscription

import (
	"context"
	"errors"

	"stackbridge-home-task/internal/errorz"
	v1 "stackbridge-home-task/internal/transport/http/v1"
)

func (h *Handler) GetSubscription(
	ctx context.Context,
	req v1.GetSubscriptionRequestObject,
) (v1.GetSubscriptionResponseObject, error) {
	sub, err := h.svc.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, errorz.ErrSubscriptionNotFound) {
			return v1.GetSubscription404JSONResponse{
				NotFoundJSONResponse: notFound("subscription not found"),
			}, nil
		}
		return nil, err
	}

	return v1.GetSubscription200JSONResponse(toV1Subscription(sub)), nil
}
