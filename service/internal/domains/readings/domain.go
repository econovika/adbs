package readings

import (
	"context"
	"database/sql"
	"energosbyt/service/internal/models"
)

type userGetter interface {
	GetUser(ctx context.Context, id uint64) (*models.User, error)
}

type gaugeChecker interface {
	GaugeExists(ctx context.Context, id uint64) (bool, error)
}

type tariffChecker interface {
	TariffExists(ctx context.Context, id uint64) (bool, error)
}

type realtyGetter interface {
	GetRealty(ctx context.Context, id uint64) (*models.Realty, error)
}

type Domain struct {
	db           *sql.DB
	usersDomain  userGetter
	gaugesDomain gaugeChecker
	tariffDomain tariffChecker
	realtyDomain realtyGetter
}

func NewDomain(
	db *sql.DB,
	usersDomain userGetter,
	gaugesDomain gaugeChecker,
	tariffDomain tariffChecker,
	realtyDomain realtyGetter,
) *Domain {
	return &Domain{
		db:           db,
		usersDomain:  usersDomain,
		gaugesDomain: gaugesDomain,
		tariffDomain: tariffDomain,
		realtyDomain: realtyDomain,
	}
}

func (d *Domain) GetUser(ctx context.Context, id uint64) (*models.User, error) {
	return d.usersDomain.GetUser(ctx, id)
}

func (d *Domain) GaugeExists(ctx context.Context, id uint64) (bool, error) {
	return d.gaugesDomain.GaugeExists(ctx, id)
}

func (d *Domain) TariffExists(ctx context.Context, id uint64) (bool, error) {
	return d.tariffDomain.TariffExists(ctx, id)
}

func (d *Domain) GetRealty(ctx context.Context, id uint64) (*models.Realty, error) {
	return d.realtyDomain.GetRealty(ctx, id)
}
