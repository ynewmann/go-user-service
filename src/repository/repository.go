package repository

import (
	"context"

	"go-user-service/src/repository/models"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
	SslMode  bool   `mapstructure:"sslmode"`
}

type Repository interface {
	Create(ctx context.Context, user models.User) (int, error)
	Get(ctx context.Context, id int) (models.User, error)
	UpdateEmail(ctx context.Context, id int, email string) error
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id int) error
}
