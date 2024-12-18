package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"go-user-service/src/repository"
	"go-user-service/src/repository/models"
)

var _ repository.Repository = (*Repository)(nil)

var (
	ErrDatabase = errors.New("database error")

	ErrNotFound = errors.New("not found")
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
	SslMode  bool   `mapstructure:"sslmode"`
}

type Repository struct {
	conn *sql.DB
}

func NewRepository(cfg Config) (*Repository, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%t",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DbName,
		cfg.SslMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return &Repository{conn: db}, nil
}

func (r *Repository) Create(ctx context.Context, user models.User) (int, error) {
	query := `INSERT INTO users(email, name) VALUES ($1)`

	id := 0
	err := r.conn.QueryRowContext(ctx, query, user.Email, user.Name).Scan(id)
	if err != nil {
		return 0, errors.Join(ErrDatabase, err)
	}

	return id, nil
}

func (r *Repository) Get(ctx context.Context, id int) (models.User, error) {
	query := "SELECT email, name FROM users WHERE id = $1"
	user := models.User{Id: id}

	row := r.conn.QueryRowContext(ctx, query, id)
	if err := row.Scan(&user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.Join(ErrDatabase, ErrNotFound)
		}

		return models.User{}, errors.Join(ErrDatabase, err)
	}

	return user, nil
}

func (r *Repository) UpdateEmail(ctx context.Context, id int, email string) error {
	query := "UPDATE users SET email = $1 WHERE id = $2"
	res, err := r.conn.ExecContext(ctx, query, email, id)
	if err != nil {
		return errors.Join(ErrDatabase, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return errors.Join(ErrDatabase, err)
	}

	if affected == 0 {
		return errors.Join(ErrDatabase, ErrNotFound)
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, user models.User) error {
	query := "UPDATE users SET email = $1, name = $2 WHERE id = $3"
	res, err := r.conn.ExecContext(ctx, query, user.Email, user.Name)
	if err != nil {
		return errors.Join(ErrDatabase, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return errors.Join(ErrDatabase, err)
	}

	if affected == 0 {
		return errors.Join(ErrDatabase, ErrNotFound)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	res, err := r.conn.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Join(ErrDatabase, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return errors.Join(ErrDatabase, err)
	}

	if affected == 0 {
		return errors.Join(ErrDatabase, ErrNotFound)
	}

	return nil
}
