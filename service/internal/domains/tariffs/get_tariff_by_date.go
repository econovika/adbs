package tariffs

import (
	"context"
	"database/sql"
	"energosbyt/service/internal/models"
	"errors"
	"fmt"
	"time"
)

func (d *Domain) GetTariffByDate(ctx context.Context, dateStr string) (*models.Tariff, error) {
	dt, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		dt, err = time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date")
		}
	}
	dtStr := dt.Format("2006-01-02 15:04:05")

	query := `
		SELECT id, name, price_day_kwph, price_night_kwph, valid_from, valid_to
		FROM dim_tariffs FINAL
		WHERE valid_from <= ? AND valid_to > ?
		ORDER BY updated_at DESC
		LIMIT 1`
	row := d.db.QueryRowContext(ctx, query, dtStr, dtStr)
	var t models.Tariff
	var vf, vt time.Time
	err = row.Scan(&t.ID, &t.Name, &t.PriceDayKWPH, &t.PriceNightKWPH, &vf, &vt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	t.ValidFrom = vf.Format("2006-01-02 15:04:05")
	t.ValidTo = vt.Format("2006-01-02 15:04:05")
	return &t, nil
}
