package bind_gauge

import (
	"context"
	"energosbyt/service/internal/models"
)

type domain interface {
	GetRealty(ctx context.Context, id uint64) (*models.Realty, error)
	GaugeExists(ctx context.Context, id uint64) (bool, error)
	BindGauge(ctx context.Context, realtyID uint64, gaugeID uint64) error
}
