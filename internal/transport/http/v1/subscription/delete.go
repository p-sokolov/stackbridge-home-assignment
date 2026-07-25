package subscription

import (
	"context"
	"errors"

	"stackbridge-home-task/internal/errorz"
	v1 "stackbridge-home-task/internal/transport/http/v1"
)

func (h *Handler) DeleteSubscription(
	ctx context.Context,
	req v1.DeleteSubscriptionRequestObject,
) (v1.DeleteSubscriptionResponseObject, error) {
	err := h.svc.Delete(ctx, req.Id)
	if err != nil {
		if errors.Is(err, errorz.ErrSubscriptionNotFound) {
			return v1.DeleteSubscription404JSONResponse{
				NotFoundJSONResponse: notFound("subscription not found"),
			}, nil
		}
		return nil, err
	}

	return v1.DeleteSubscription204Response{}, nil
}
