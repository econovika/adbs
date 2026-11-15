package invoices

import (
	"context"
	"database/sql"
	"energosbyt/service/internal/models"
)

type realtyGetter interface {
	GetRealty(ctx context.Context, id uint64) (*models.Realty, error)
}

type Domain struct {
	db           *sql.DB
	realtyDomain realtyGetter
}

func NewDomain(db *sql.DB, realtyDomain realtyGetter) *Domain {
	return &Domain{
		db:           db,
		realtyDomain: realtyDomain,
	}
}

func (d *Domain) GetRealty(ctx context.Context, id uint64) (*models.Realty, error) {
	return d.realtyDomain.GetRealty(ctx, id)
}
