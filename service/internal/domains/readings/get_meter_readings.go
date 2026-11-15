package readings

import (
	"context"
	"energosbyt/service/internal/models"
	"time"
)

func (d *Domain) GetMeterReadings(ctx context.Context, realtyID uint64, fromDate, toDate string) ([]models.MeterReading, error) {
	query := `
		SELECT reading_date, user_id, tariff_id, gauge_id, day_kw, night_kw
		FROM fact_meter_readings FINAL
		WHERE realty_id = ? AND reading_date >= ? AND reading_date <= ?
		ORDER BY reading_date ASC`
	rows, err := d.db.QueryContext(ctx, query, realtyID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	readings := make([]models.MeterReading, 0)
	for rows.Next() {
		var mr models.MeterReading
		var d time.Time
		if err := rows.Scan(&d, &mr.UserID, &mr.TariffID, &mr.GaugeID, &mr.DayKW, &mr.NightKW); err != nil {
			return nil, err
		}
		mr.ReadingDate = d.Format("2006-01-02")
		readings = append(readings, mr)
	}
	return readings, nil
}
