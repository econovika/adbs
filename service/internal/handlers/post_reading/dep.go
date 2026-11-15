package post_reading

import (
	"context"
	"energosbyt/service/internal/models"
)

type domain interface {
	GetUser(ctx context.Context, id uint64) (*models.User, error)
	GaugeExists(ctx context.Context, id uint64) (bool, error)
	TariffExists(ctx context.Context, id uint64) (bool, error)
	AddMeterReading(ctx context.Context, realtyID uint64, req *models.PostReadingRequest) error
	GetRealty(ctx context.Context, id uint64) (*models.Realty, error)
}
