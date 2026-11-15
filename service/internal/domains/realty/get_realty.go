package realty

import (
	"context"
	"database/sql"
	"energosbyt/service/internal/models"
	"errors"
)

func (d *Domain) GetRealty(ctx context.Context, id uint64) (*models.Realty, error) {
	query := `SELECT id, address, square, type, number, rooms, floor, porch FROM dim_realty WHERE id = ? LIMIT 1`
	row := d.db.QueryRowContext(ctx, query, id)
	var re models.Realty
	err := row.Scan(&re.ID, &re.Address, &re.Square, &re.Type, &re.Number, &re.Rooms, &re.Floor, &re.Porch)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &re, nil
}
