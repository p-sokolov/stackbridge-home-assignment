package subscription

import (
	"context"

	v1 "stackbridge-home-task/internal/transport/http/v1"
)

func (h *Handler) ListSubscriptions(
	ctx context.Context,
	req v1.ListSubscriptionsRequestObject,
) (v1.ListSubscriptionsResponseObject, error) {
	subs, err := h.svc.List(ctx, req.Params.UserId)
	if err != nil {
		return nil, err
	}

	return v1.ListSubscriptions200JSONResponse(toV1Subscriptions(subs)), nil
}