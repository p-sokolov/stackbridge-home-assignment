package subscription

import (
	"context"
	"errors"

	"stackbridge-home-task/internal/errorz"
	"stackbridge-home-task/internal/models"
	v1 "stackbridge-home-task/internal/transport/http/v1"
)

func (h *Handler) UpdateSubscription(
	ctx context.Context,
	req v1.UpdateSubscriptionRequestObject,
) (v1.UpdateSubscriptionResponseObject, error) {
	if req.Body == nil {
		return v1.UpdateSubscription400JSONResponse{
			BadRequestJSONResponse: badRequest("request body is required"),
		}, nil
	}

	input := &models.SubscriptionDTO{
		ServiceName: req.Body.ServiceName,
		Price:       req.Body.Price,
		UserID:      req.Body.UserId,
		StartDate:   req.Body.StartDate,
		EndDate:     req.Body.EndDate,
	}

	sub, err := h.svc.Update(ctx, req.Id, input)
	if err != nil {
		if errors.Is(err, errorz.ErrInvalidSubscriptionData) {
			return v1.UpdateSubscription400JSONResponse{
				BadRequestJSONResponse: badRequest("invalid subscription data"),
			}, nil
		}
		if errors.Is(err, errorz.ErrSubscriptionNotFound) {
			return v1.UpdateSubscription404JSONResponse{
				NotFoundJSONResponse: notFound("subscription not found"),
			}, nil
		}
		return nil, err
	}

	return v1.UpdateSubscription200JSONResponse(toV1Subscription(sub)), nil
}