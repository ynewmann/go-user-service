package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"go-user-service/src/controllers"
	"go-user-service/src/repository/models"
)

var (
	ErrBadUserId      = errors.New("bad user id")
	ErrInternal       = errors.New("internal error")
	ErrBadUserPayload = errors.New("bad user payload")
	ErrNoEmail        = errors.New("no email provided")
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type CreateResponse struct {
	Id int `json:"id"`
}

type UpdateRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UpdateEmailRequest struct {
	Email *string `json:"email"`
}

type Handler struct {
	controller *controllers.Controller
}

func New(controller *controllers.Controller) *Handler {
	return &Handler{controller: controller}
}

func (h *Handler) Create(c *fiber.Ctx) error {
	req := CreateRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(ErrBadUserPayload.Error())
	}

	id, err := h.controller.Create(c.UserContext(), models.User{Email: req.Email, Name: req.Name})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(ErrInternal.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&CreateResponse{Id: id})
}

func (h *Handler) Get(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: ErrBadUserPayload.Error()})
	}

	user, err := h.controller.Get(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: ErrInternal.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(&user)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: ErrBadUserId.Error()})
	}

	req := UpdateRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: ErrBadUserPayload.Error()})
	}

	user := models.User{
		Id:    id,
		Email: req.Email,
		Name:  req.Name,
	}
	err = h.controller.Update(c.UserContext(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: ErrInternal.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *Handler) UpdateEmail(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: ErrBadUserId.Error()})
	}

	req := UpdateEmailRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: ErrBadUserPayload.Error()})
	}

	if req.Email == nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: ErrNoEmail.Error()})
	}

	err = h.controller.UpdateEmail(c.UserContext(), id, *req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: ErrInternal.Error()})
	}

	user, err := h.controller.Get(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: ErrInternal.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: ErrBadUserId.Error()})
	}

	err = h.controller.Delete(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: ErrInternal.Error()})
	}

	return c.Status(fiber.StatusOK).Send(nil)
}
