package repository

import (
	"context"

	"go-user-service/src/repository/models"
)

type Repository interface {
	Create(ctx context.Context, user models.User) (int, error)
	Get(ctx context.Context, id int) (models.User, error)
	UpdateEmail(ctx context.Context, id int, email string) error
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id int) error
}
