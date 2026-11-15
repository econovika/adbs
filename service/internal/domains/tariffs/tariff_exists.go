package tariffs

import "context"

func (d *Domain) TariffExists(ctx context.Context, id uint64) (bool, error) {
	var count int
	err := d.db.QueryRowContext(ctx, `SELECT count() FROM dim_tariffs FINAL WHERE id = ?`, id).Scan(&count)
	return count > 0, err
}
