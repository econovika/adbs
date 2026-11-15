package readings

import (
	"context"
	"energosbyt/service/internal/models"
	"time"
)

func (d *Domain) AddMeterReading(ctx context.Context, realtyID uint64, req *models.PostReadingRequest) error {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	query := `INSERT INTO fact_meter_readings (realty_id, user_id, tariff_id, gauge_id, reading_date, day_kw, night_kw, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := d.db.ExecContext(ctx, query, realtyID, req.UserID, req.TariffID, req.GaugeID, req.ReadingDate, req.DayKW, req.NightKW, now)
	return err
}
