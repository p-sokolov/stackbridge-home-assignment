package subscription

import (
	"context"

	v1 "stackbridge-home-task/internal/transport/http/v1"
)

func (h *Handler) GetSubscription(
	ctx context.Context,
	req v1.GetSubscriptionRequestObject,
) (v1.GetSubscriptionResponseObject, error) {
	sub, err := h.svc.GetByID(ctx, req.Id)
	if err != nil {
		if isNotFound(err) {
			return v1.GetSubscription404JSONResponse{
				NotFoundJSONResponse: notFound(err.Error()),
			}, nil
		}
		return nil, err
	}

	return v1.GetSubscription200JSONResponse(toV1Subscription(sub)), nil
}
