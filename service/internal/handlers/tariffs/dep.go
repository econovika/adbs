package tariffs

import (
	"context"
	"energosbyt/service/internal/models"
)

type domain interface {
	GetTariffByDate(ctx context.Context, dateStr string) (*models.Tariff, error)
}
