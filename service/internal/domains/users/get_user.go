package users

import (
	"context"
	"database/sql"
	"energosbyt/service/internal/models"
	"errors"
	"time"
)

func (d *Domain) GetUser(ctx context.Context, id uint64) (*models.User, error) {
	query := `SELECT id, fio, phone, phone_prev, created_at FROM dim_users FINAL WHERE id = ?`
	row := d.db.QueryRowContext(ctx, query, id)
	var u models.User
	var ca time.Time
	err := row.Scan(&u.ID, &u.FIO, &u.Phone, &u.PhonePrev, &ca)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u.CreatedAt = ca.Format("2006-01-02 15:04:05")
	return &u, nil
}
