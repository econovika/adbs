package get_invoices

import (
	"context"
	"energosbyt/service/internal/models"
)

type domain interface {
	GetRealty(ctx context.Context, id uint64) (*models.Realty, error)
	GetInvoices(ctx context.Context, realtyID uint64, limit, offset int) ([]models.Invoice, int, error)
}
