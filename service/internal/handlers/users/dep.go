package users

import (
	"context"
	"energosbyt/service/internal/models"
)

type domain interface {
	GetUser(ctx context.Context, id uint64) (*models.User, error)
	UpdateUser(ctx context.Context, id uint64, fio *string, phone *string) error
}
