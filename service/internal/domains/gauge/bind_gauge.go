package gauge

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

func (d *Domain) BindGauge(ctx context.Context, realtyID uint64, gaugeID uint64) error {
	queryActive := `
		SELECT gauge_id, version
		FROM fact_gauge_bindings
		WHERE realty_id = ?
		GROUP BY gauge_id, version
		HAVING sum(sign) > 0
		ORDER BY version DESC
		LIMIT 1`
	row := d.db.QueryRowContext(ctx, queryActive, realtyID)

	var activeGaugeID uint64
	var activeVersion uint64
	hasActive := false
	err := row.Scan(&activeGaugeID, &activeVersion)
	if err == nil {
		hasActive = true
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	now := time.Now().UTC().Format("2006-01-02 15:04:05")

	if hasActive {
		cancelQuery := `INSERT INTO fact_gauge_bindings (gauge_id, realty_id, created_at, sign, version) VALUES (?, ?, ?, -1, ?)`
		if _, err := d.db.ExecContext(ctx, cancelQuery, activeGaugeID, realtyID, now, activeVersion); err != nil {
			return err
		}
	}

	newVersion := uint64(1)
	if hasActive {
		newVersion = activeVersion + 1
	}

	insertQuery := `INSERT INTO fact_gauge_bindings (gauge_id, realty_id, created_at, sign, version) VALUES (?, ?, ?, 1, ?)`
	_, err = d.db.ExecContext(ctx, insertQuery, gaugeID, realtyID, now, newVersion)
	return err
}
