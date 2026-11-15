package users

import (
	"context"
	"database/sql"
	"time"
)

func (d *Domain) UpdateUser(ctx context.Context, id uint64, fio *string, phone *string) error {
	u, err := d.GetUser(ctx, id)
	if err != nil {
		return err
	}
	if u == nil {
		return sql.ErrNoRows
	}

	newFIO := u.FIO
	if fio != nil {
		newFIO = *fio
	}

	newPhone := u.Phone
	newPhonePrev := u.PhonePrev
	if phone != nil && *phone != u.Phone {
		newPhonePrev = u.Phone
		newPhone = *phone
	}

	createdAt := u.CreatedAt
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	if createdAt == "" {
		createdAt = now
	}

	insertQuery := `INSERT INTO dim_users (id, fio, phone, phone_prev, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = d.db.ExecContext(ctx, insertQuery, u.ID, newFIO, newPhone, newPhonePrev, createdAt, now)
	return err
}
