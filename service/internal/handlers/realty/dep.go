package realty

import (
	"context"
	"energosbyt/service/internal/models"
)

type domain interface {
	GetRealty(ctx context.Context, id uint64) (*models.Realty, error)
}
