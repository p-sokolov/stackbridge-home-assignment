package subscription

import (
	"context"
	"errors"
	"testing"
	"time"

	"stackbridge-home-task/internal/errorz"
	"stackbridge-home-task/internal/models"

	"github.com/google/uuid"
)

type fakeRepo struct {
	createCalled bool
	createSub    *models.Subscription

	totalCalled bool
	totalUserID uuid.UUID
	totalName   string
	totalStart  time.Time
	totalEnd    time.Time
}

func (r *fakeRepo) Create(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	r.createCalled = true
	r.createSub = sub
	sub.ID = uuid.New()
	return sub, nil
}

func (r *fakeRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	return nil, nil
}

func (r *fakeRepo) List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Subscription, error) {
	return nil, nil
}

func (r *fakeRepo) Update(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	return sub, nil
}

func (r *fakeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *fakeRepo) CalculateTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, startPeriod, endPeriod time.Time) (int, error) {
	r.totalCalled = true
	r.totalUserID = userID
	r.totalName = serviceName
	r.totalStart = startPeriod
	r.totalEnd = endPeriod
	return 1200, nil
}

func TestCreateBuildsValidSubscription(t *testing.T) {
	userID := uuid.New()
	endDate := "12-2025"
	repo := &fakeRepo{}
	svc := New(repo)

	sub, err := svc.Create(context.Background(), &models.SubscriptionDTO{
		ServiceName: "  Yandex Plus  ",
		Price:       400,
		UserID:      userID,
		StartDate:   "07-2025",
		EndDate:     &endDate,
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if !repo.createCalled {
		t.Fatal("expected repo.Create to be called")
	}
	if sub.ServiceName != "Yandex Plus" {
		t.Fatalf("expected trimmed service name, got %q", sub.ServiceName)
	}
	if sub.StartDate.Format(monthLayout) != "07-2025" {
		t.Fatalf("unexpected start date: %s", sub.StartDate.Format(monthLayout))
	}
	if sub.EndDate == nil || sub.EndDate.Format(monthLayout) != "12-2025" {
		t.Fatalf("unexpected end date: %v", sub.EndDate)
	}
}

func TestCreateRejectsInvalidSubscription(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name  string
		input *models.SubscriptionDTO
	}{
		{
			name: "empty service name",
			input: &models.SubscriptionDTO{
				ServiceName: "",
				Price:       400,
				UserID:      userID,
				StartDate:   "07-2025",
			},
		},
		{
			name: "zero price",
			input: &models.SubscriptionDTO{
				ServiceName: "Yandex Plus",
				Price:       0,
				UserID:      userID,
				StartDate:   "07-2025",
			},
		},
		{
			name: "nil user id",
			input: &models.SubscriptionDTO{
				ServiceName: "Yandex Plus",
				Price:       400,
				UserID:      uuid.Nil,
				StartDate:   "07-2025",
			},
		},
		{
			name: "invalid start date",
			input: &models.SubscriptionDTO{
				ServiceName: "Yandex Plus",
				Price:       400,
				UserID:      userID,
				StartDate:   "2025-07",
			},
		},
		{
			name: "end date before start date",
			input: &models.SubscriptionDTO{
				ServiceName: "Yandex Plus",
				Price:       400,
				UserID:      userID,
				StartDate:   "07-2025",
				EndDate:     ptr("06-2025"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{}
			svc := New(repo)

			_, err := svc.Create(context.Background(), tt.input)
			if !errors.Is(err, errorz.ErrInvalidSubscriptionData) {
				t.Fatalf("expected ErrInvalidSubscriptionData, got %v", err)
			}
			if repo.createCalled {
				t.Fatal("repo.Create should not be called for invalid input")
			}
		})
	}
}

func TestCalculateTotalCostParsesInputAndCallsRepo(t *testing.T) {
	userID := uuid.New()
	repo := &fakeRepo{}
	svc := New(repo)

	total, err := svc.CalculateTotalCost(context.Background(), userID, "  Yandex Plus  ", "01-2025", "03-2025")
	if err != nil {
		t.Fatalf("CalculateTotalCost returned error: %v", err)
	}
	if total != 1200 {
		t.Fatalf("expected total 1200, got %d", total)
	}
	if !repo.totalCalled {
		t.Fatal("expected repo.CalculateTotalCost to be called")
	}
	if repo.totalName != "Yandex Plus" {
		t.Fatalf("expected trimmed service name, got %q", repo.totalName)
	}
	if repo.totalStart.Format(monthLayout) != "01-2025" || repo.totalEnd.Format(monthLayout) != "03-2025" {
		t.Fatalf("unexpected period: %s - %s", repo.totalStart.Format(monthLayout), repo.totalEnd.Format(monthLayout))
	}
}

func ptr(s string) *string {
	return &s
}
