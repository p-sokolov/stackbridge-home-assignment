package subscription

import (
	"context"
	"errors"

	"stackbridge-home-task/internal/errorz"
	"stackbridge-home-task/internal/models"
	v1 "stackbridge-home-task/internal/transport/http/v1"
)

func (h *Handler) CreateSubscription(
	ctx context.Context,
	req v1.CreateSubscriptionRequestObject,
) (v1.CreateSubscriptionResponseObject, error) {
	if req.Body == nil {
		return v1.CreateSubscription400JSONResponse{
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

	sub, err := h.svc.Create(ctx, input)
	if err != nil {
		if errors.Is(err, errorz.ErrInvalidSubscriptionData) {
			return v1.CreateSubscription400JSONResponse{
				BadRequestJSONResponse: badRequest("invalid subscription data"),
			}, nil
		}
		if errors.Is(err, errorz.ErrSubscriptionAlreadyExists) {
			return v1.CreateSubscription400JSONResponse{
				BadRequestJSONResponse: badRequest("subscription already exists"),
			}, nil
		}
		return nil, err
	}

	return v1.CreateSubscription201JSONResponse(toV1Subscription(sub)), nil
}