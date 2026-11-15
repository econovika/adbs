package get_readings

import (
	"context"
	"energosbyt/service/internal/models"
)

type domain interface {
	GetRealty(ctx context.Context, id uint64) (*models.Realty, error)
	GetMeterReadings(ctx context.Context, realtyID uint64, fromDate, toDate string) ([]models.MeterReading, error)
}
