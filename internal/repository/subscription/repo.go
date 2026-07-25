package subscription

import (
	"time"

	"stackbridge-home-task/internal/db/sqlc/storage"
	"stackbridge-home-task/internal/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	queries *storage.Queries
}

func New(db *pgxpool.Pool) *repo {
	return &repo{
		queries: storage.New(db),
	}
}

func toModelSubscription(u storage.Subscription) *models.Subscription {
	return &models.Subscription{
		ID:          u.ID,
		ServiceName: u.ServiceName,
		Price:       int(u.Price),
		UserID:      u.UserID,
		StartDate:   u.StartDate.Time,
		EndDate:     u.EndDate,
	}
}

func toPgDate(t time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  t,
		Valid: !t.IsZero(),
	}
}