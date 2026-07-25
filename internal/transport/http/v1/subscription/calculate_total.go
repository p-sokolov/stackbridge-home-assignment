package subscription

import (
	"context"
	"errors"

	"stackbridge-home-task/internal/errorz"
	v1 "stackbridge-home-task/internal/transport/http/v1"
)

func (h *Handler) CalculateTotalCost(
	ctx context.Context,
	req v1.CalculateTotalCostRequestObject,
) (v1.CalculateTotalCostResponseObject, error) {
	total, err := h.svc.CalculateTotalCost(
		ctx,
		req.Params.UserId,
		req.Params.ServiceName,
		req.Params.StartPeriod,
		req.Params.EndPeriod,
	)
	if err != nil {
		if errors.Is(err, errorz.ErrInvalidSubscriptionData) {
			return v1.CalculateTotalCost400JSONResponse{
				BadRequestJSONResponse: badRequest("invalid subscription data"),
			}, nil
		}
		return nil, err
	}

	return v1.CalculateTotalCost200JSONResponse{
		TotalCost: total,
	}, nil
}