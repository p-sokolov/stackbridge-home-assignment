package subscription

import (
	"context"

	v1 "stackbridge-home-task/internal/transport/http/v1"
)

func (h *Handler) ListSubscriptions(
	ctx context.Context,
	req v1.ListSubscriptionsRequestObject,
) (v1.ListSubscriptionsResponseObject, error) {
	limit := 0
	if req.Params.LimitCnt != nil {
		limit = *req.Params.LimitCnt
	}

	offset := 0
	if req.Params.OffsetCnt != nil {
		offset = *req.Params.OffsetCnt
	}

	subs, err := h.svc.List(ctx, req.Params.UserId, limit, offset)
	if err != nil {
		if isBadRequest(err) {
			return v1.ListSubscriptions400JSONResponse{
				BadRequestJSONResponse: badRequest("invalid subscription data"),
			}, nil
		}
		return nil, err
	}

	return v1.ListSubscriptions200JSONResponse(toV1Subscriptions(subs)), nil
}
