package invoices

import (
	"context"
	"energosbyt/service/internal/models"
	"time"
)

func (d *Domain) GetInvoices(ctx context.Context, realtyID uint64, limit, offset int) ([]models.Invoice, int, error) {
	var total int
	err := d.db.QueryRowContext(ctx, `SELECT count() FROM fact_invoice WHERE realty_id = ?`, realtyID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT tariff_id, gauge_id, day_kw, night_kw, price, deadline, created_at
		FROM fact_invoice
		WHERE realty_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`
	rows, err := d.db.QueryContext(ctx, query, realtyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	invoices := make([]models.Invoice, 0)
	for rows.Next() {
		var inv models.Invoice
		var dl time.Time
		var ca time.Time
		if err := rows.Scan(&inv.TariffID, &inv.GaugeID, &inv.DayKW, &inv.NightKW, &inv.Price, &dl, &ca); err != nil {
			return nil, 0, err
		}
		inv.Deadline = dl.Format("2006-01-02")
		inv.CreatedAt = ca.Format("2006-01-02 15:04:05")
		invoices = append(invoices, inv)
	}
	return invoices, total, nil
}
