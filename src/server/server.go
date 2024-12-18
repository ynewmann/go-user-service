package server

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"go-user-service/src/handlers"
)

type Config struct {
	Port string
}

type Server struct {
	cfg Config

	app     *fiber.App
	handler handlers.Handler
}

func New(cfg Config, handler handlers.Handler) *Server {
	app := fiber.New()

	app.Post("/users", handler.Create)
	app.Get("/users/:id", handler.Get)
	app.Put("/users/:id", handler.Update)
	app.Patch("/users/:id", handler.UpdateEmail)
	app.Delete("/users/:id", handler.Delete)

	return &Server{
		cfg:     cfg,
		app:     app,
		handler: handler,
	}
}

func (s *Server) Start() error {
	return s.app.Listen(":" + s.cfg.Port)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
