package subscription

import (
	"context"
	"errors"
	"testing"
	"time"

	"stackbridge-home-task/internal/errorz"
	"stackbridge-home-task/internal/models"
	v1 "stackbridge-home-task/internal/transport/http/v1"

	"github.com/google/uuid"
)

type fakeService struct {
	createFn             func(context.Context, *models.SubscriptionDTO) (*models.Subscription, error)
	getByIDFn            func(context.Context, uuid.UUID) (*models.Subscription, error)
	listFn               func(context.Context, uuid.UUID) ([]*models.Subscription, error)
	updateFn             func(context.Context, uuid.UUID, *models.SubscriptionDTO) (*models.Subscription, error)
	deleteFn             func(context.Context, uuid.UUID) error
	calculateTotalCostFn func(context.Context, uuid.UUID, string, string, string) (int, error)
}

func (s fakeService) Create(ctx context.Context, input *models.SubscriptionDTO) (*models.Subscription, error) {
	return s.createFn(ctx, input)
}

func (s fakeService) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	return s.getByIDFn(ctx, id)
}

func (s fakeService) List(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error) {
	return s.listFn(ctx, userID)
}

func (s fakeService) Update(ctx context.Context, id uuid.UUID, input *models.SubscriptionDTO) (*models.Subscription, error) {
	return s.updateFn(ctx, id, input)
}

func (s fakeService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.deleteFn(ctx, id)
}

func (s fakeService) CalculateTotalCost(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	startPeriod string,
	endPeriod string,
) (int, error) {
	return s.calculateTotalCostFn(ctx, userID, serviceName, startPeriod, endPeriod)
}

func TestCreateSubscriptionSuccess(t *testing.T) {
	userID := uuid.New()
	subID := uuid.New()

	h := New(fakeService{
		createFn: func(ctx context.Context, input *models.SubscriptionDTO) (*models.Subscription, error) {
			if input.ServiceName != "Yandex Plus" {
				t.Fatalf("unexpected service name: %q", input.ServiceName)
			}
			if input.Price != 400 {
				t.Fatalf("unexpected price: %d", input.Price)
			}
			if input.UserID != userID {
				t.Fatalf("unexpected user id: %s", input.UserID)
			}
			if input.StartDate != "07-2025" {
				t.Fatalf("unexpected start date: %q", input.StartDate)
			}

			return &models.Subscription{
				ID:          subID,
				ServiceName: input.ServiceName,
				Price:       input.Price,
				UserID:      input.UserID,
				StartDate:   mustMonth(t, "07-2025"),
			}, nil
		},
	})

	body := v1.CreateSubscriptionJSONRequestBody{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserId:      userID,
		StartDate:   "07-2025",
	}

	res, err := h.CreateSubscription(context.Background(), v1.CreateSubscriptionRequestObject{Body: &body})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := res.(v1.CreateSubscription201JSONResponse)
	if !ok {
		t.Fatalf("expected 201 response, got %T", res)
	}

	if created.Id != subID {
		t.Fatalf("unexpected subscription id: %s", created.Id)
	}
	if created.StartDate != "07-2025" {
		t.Fatalf("unexpected start date: %q", created.StartDate)
	}
}

func TestCreateSubscriptionBadRequest(t *testing.T) {
	h := New(fakeService{})

	res, err := h.CreateSubscription(context.Background(), v1.CreateSubscriptionRequestObject{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := res.(v1.CreateSubscription400JSONResponse); !ok {
		t.Fatalf("expected 400 response, got %T", res)
	}
}

func TestCreateSubscriptionValidationError(t *testing.T) {
	h := New(fakeService{
		createFn: func(ctx context.Context, input *models.SubscriptionDTO) (*models.Subscription, error) {
			return nil, errorz.ErrInvalidSubscriptionData
		},
	})

	body := v1.CreateSubscriptionJSONRequestBody{
		ServiceName: "Yandex Plus",
		Price:       -1,
		UserId:      uuid.New(),
		StartDate:   "07-2025",
	}

	res, err := h.CreateSubscription(context.Background(), v1.CreateSubscriptionRequestObject{Body: &body})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := res.(v1.CreateSubscription400JSONResponse); !ok {
		t.Fatalf("expected 400 response, got %T", res)
	}
}

func TestGetSubscriptionSuccess(t *testing.T) {
	subID := uuid.New()
	userID := uuid.New()

	h := New(fakeService{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
			if id != subID {
				t.Fatalf("unexpected id: %s", id)
			}

			return &models.Subscription{
				ID:          subID,
				ServiceName: "Yandex Plus",
				Price:       400,
				UserID:      userID,
				StartDate:   mustMonth(t, "07-2025"),
			}, nil
		},
	})

	res, err := h.GetSubscription(context.Background(), v1.GetSubscriptionRequestObject{Id: subID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, ok := res.(v1.GetSubscription200JSONResponse)
	if !ok {
		t.Fatalf("expected 200 response, got %T", res)
	}

	if got.Id != subID || got.UserId != userID {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestGetSubscriptionNotFound(t *testing.T) {
	h := New(fakeService{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
			return nil, errorz.ErrSubscriptionNotFound
		},
	})

	res, err := h.GetSubscription(context.Background(), v1.GetSubscriptionRequestObject{Id: uuid.New()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := res.(v1.GetSubscription404JSONResponse); !ok {
		t.Fatalf("expected 404 response, got %T", res)
	}
}

func TestListSubscriptionsSuccess(t *testing.T) {
	userID := uuid.New()

	h := New(fakeService{
		listFn: func(ctx context.Context, id uuid.UUID) ([]*models.Subscription, error) {
			if id != userID {
				t.Fatalf("unexpected user id: %s", id)
			}

			return []*models.Subscription{
				{
					ID:          uuid.New(),
					ServiceName: "Yandex Plus",
					Price:       400,
					UserID:      userID,
					StartDate:   mustMonth(t, "07-2025"),
				},
				{
					ID:          uuid.New(),
					ServiceName: "Netflix",
					Price:       700,
					UserID:      userID,
					StartDate:   mustMonth(t, "08-2025"),
				},
			}, nil
		},
	})

	res, err := h.ListSubscriptions(context.Background(), v1.ListSubscriptionsRequestObject{
		Params: v1.ListSubscriptionsParams{UserId: userID},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, ok := res.(v1.ListSubscriptions200JSONResponse)
	if !ok {
		t.Fatalf("expected 200 response, got %T", res)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 subscriptions, got %d", len(list))
	}
}

func TestUpdateSubscriptionSuccess(t *testing.T) {
	subID := uuid.New()
	userID := uuid.New()

	h := New(fakeService{
		updateFn: func(ctx context.Context, id uuid.UUID, input *models.SubscriptionDTO) (*models.Subscription, error) {
			if id != subID {
				t.Fatalf("unexpected id: %s", id)
			}

			return &models.Subscription{
				ID:          id,
				ServiceName: input.ServiceName,
				Price:       input.Price,
				UserID:      input.UserID,
				StartDate:   mustMonth(t, input.StartDate),
			}, nil
		},
	})

	body := v1.UpdateSubscriptionJSONRequestBody{
		ServiceName: "Yandex Plus",
		Price:       500,
		UserId:      userID,
		StartDate:   "07-2025",
	}

	res, err := h.UpdateSubscription(context.Background(), v1.UpdateSubscriptionRequestObject{
		Id:   subID,
		Body: &body,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, ok := res.(v1.UpdateSubscription200JSONResponse)
	if !ok {
		t.Fatalf("expected 200 response, got %T", res)
	}
	if updated.Price != 500 {
		t.Fatalf("expected price 500, got %d", updated.Price)
	}
}

func TestUpdateSubscriptionNotFound(t *testing.T) {
	h := New(fakeService{
		updateFn: func(ctx context.Context, id uuid.UUID, input *models.SubscriptionDTO) (*models.Subscription, error) {
			return nil, errorz.ErrSubscriptionNotFound
		},
	})

	body := v1.UpdateSubscriptionJSONRequestBody{
		ServiceName: "Yandex Plus",
		Price:       500,
		UserId:      uuid.New(),
		StartDate:   "07-2025",
	}

	res, err := h.UpdateSubscription(context.Background(), v1.UpdateSubscriptionRequestObject{
		Id:   uuid.New(),
		Body: &body,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := res.(v1.UpdateSubscription404JSONResponse); !ok {
		t.Fatalf("expected 404 response, got %T", res)
	}
}

func TestUpdateSubscriptionBadRequestWithoutBody(t *testing.T) {
	h := New(fakeService{})

	res, err := h.UpdateSubscription(context.Background(), v1.UpdateSubscriptionRequestObject{Id: uuid.New()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := res.(v1.UpdateSubscription400JSONResponse); !ok {
		t.Fatalf("expected 400 response, got %T", res)
	}
}

func TestDeleteSubscriptionSuccess(t *testing.T) {
	subID := uuid.New()

	h := New(fakeService{
		deleteFn: func(ctx context.Context, id uuid.UUID) error {
			if id != subID {
				t.Fatalf("unexpected id: %s", id)
			}
			return nil
		},
	})

	res, err := h.DeleteSubscription(context.Background(), v1.DeleteSubscriptionRequestObject{Id: subID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := res.(v1.DeleteSubscription204Response); !ok {
		t.Fatalf("expected 204 response, got %T", res)
	}
}

func TestDeleteSubscriptionNotFound(t *testing.T) {
	h := New(fakeService{
		deleteFn: func(ctx context.Context, id uuid.UUID) error {
			return errorz.ErrSubscriptionNotFound
		},
	})

	res, err := h.DeleteSubscription(context.Background(), v1.DeleteSubscriptionRequestObject{Id: uuid.New()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := res.(v1.DeleteSubscription404JSONResponse); !ok {
		t.Fatalf("expected 404 response, got %T", res)
	}
}

func TestCalculateTotalCostSuccess(t *testing.T) {
	userID := uuid.New()

	h := New(fakeService{
		calculateTotalCostFn: func(ctx context.Context, id uuid.UUID, serviceName, startPeriod, endPeriod string) (int, error) {
			if id != userID {
				t.Fatalf("unexpected user id: %s", id)
			}
			if serviceName != "Yandex Plus" {
				t.Fatalf("unexpected service name: %q", serviceName)
			}
			if startPeriod != "07-2025" || endPeriod != "12-2025" {
				t.Fatalf("unexpected period: %s - %s", startPeriod, endPeriod)
			}

			return 3000, nil
		},
	})

	res, err := h.CalculateTotalCost(context.Background(), v1.CalculateTotalCostRequestObject{
		Params: v1.CalculateTotalCostParams{
			UserId:      userID,
			ServiceName: "Yandex Plus",
			StartPeriod: "07-2025",
			EndPeriod:   "12-2025",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total, ok := res.(v1.CalculateTotalCost200JSONResponse)
	if !ok {
		t.Fatalf("expected 200 response, got %T", res)
	}
	if total.TotalCost != 3000 {
		t.Fatalf("expected total 3000, got %d", total.TotalCost)
	}
}

func TestCalculateTotalCostBadRequest(t *testing.T) {
	h := New(fakeService{
		calculateTotalCostFn: func(ctx context.Context, id uuid.UUID, serviceName, startPeriod, endPeriod string) (int, error) {
			return 0, errorz.ErrInvalidSubscriptionData
		},
	})

	res, err := h.CalculateTotalCost(context.Background(), v1.CalculateTotalCostRequestObject{
		Params: v1.CalculateTotalCostParams{
			UserId:      uuid.New(),
			ServiceName: "",
			StartPeriod: "bad-date",
			EndPeriod:   "12-2025",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := res.(v1.CalculateTotalCost400JSONResponse); !ok {
		t.Fatalf("expected 400 response, got %T", res)
	}
}

func TestHandlerReturnsUnexpectedErrors(t *testing.T) {
	wantErr := errors.New("database unavailable")

	h := New(fakeService{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
			return nil, wantErr
		},
	})

	res, err := h.GetSubscription(context.Background(), v1.GetSubscriptionRequestObject{Id: uuid.New()})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected unexpected error, got %v", err)
	}
	if res != nil {
		t.Fatalf("expected nil response, got %T", res)
	}
}

func mustMonth(t *testing.T, value string) time.Time {
	t.Helper()

	month, err := time.Parse(monthLayout, value)
	if err != nil {
		t.Fatalf("failed to parse test month %q: %v", value, err)
	}

	return month
}
