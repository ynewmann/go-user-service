package controllers

import (
	"context"
	"errors"

	"go-user-service/src/repository"
	"go-user-service/src/repository/models"
)

var (
	ErrBadName  = errors.New("bad email")
	ErrBadEmail = errors.New("bad name")
)

type Controller struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) Create(ctx context.Context, user models.User) (int, error) {
	err := validateUser(user)
	if err != nil {
		return 0, err
	}

	return c.repo.Create(ctx, user)
}

func (c *Controller) Get(ctx context.Context, id int) (models.User, error) {
	return c.repo.Get(ctx, id)
}

func (c *Controller) Update(ctx context.Context, user models.User) error {
	err := validateUser(user)
	if err != nil {
		return err
	}

	return c.repo.Update(ctx, user)
}

func (c *Controller) UpdateEmail(ctx context.Context, id int, email string) error {
	if email == "" {
		return ErrBadEmail
	}

	return c.repo.UpdateEmail(ctx, id, email)
}

func (c *Controller) Delete(ctx context.Context, id int) error {
	return c.repo.Delete(ctx, id)
}

func validateUser(user models.User) error {
	if user.Email == "" {
		return ErrBadEmail
	}

	if user.Name == "" {
		return ErrBadName
	}

	return nil
}
