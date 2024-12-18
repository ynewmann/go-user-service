package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"go.uber.org/zap"

	"go-user-service/src/controllers"
	"go-user-service/src/handlers"
	"go-user-service/src/repository"
	"go-user-service/src/repository/models"
	"go-user-service/src/repository/postgres"
)

func setupApp() (*Server, error) {
	ctx := context.Background()

	testUser := "user"
	testPass := "pass"
	testDb := "db"

	pgContainer, err := testpostgres.Run(
		ctx,
		"postgres:alpine",
		testpostgres.WithDatabase(testDb),
		testpostgres.WithUsername(testUser),
		testpostgres.WithPassword(testPass),
		testpostgres.BasicWaitStrategies(),
		testcontainers.WithLogger(testcontainers.Logger),
	)
	if err != nil {
		return nil, err
	}

	dur := time.Second
	defer pgContainer.Stop(ctx, &dur)

	dbHost, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	dbPortStr, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	repo, err := postgres.NewRepository(repository.Config{
		Host:     dbHost,
		Port:     dbPortStr.Int(),
		User:     testUser,
		Password: testPass,
		DbName:   testDb,
		SslMode:  "disable",
	})
	if err != nil {
		return nil, err
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	userController := controllers.New(repo)
	userHandler := handlers.New(userController)
	microservice := New(Config{Port: "8081"}, logger, userHandler)
	go microservice.Start()

	return microservice, nil
}

func TestUser(t *testing.T) {
	server, err := setupApp()
	require.NoError(t, err, err)

	user := handlers.CreateRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}
	userId := 0
	t.Run("Create", func(t *testing.T) {
		body, _ := json.Marshal(user)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var createResp handlers.CreateResponse
		err = json.NewDecoder(resp.Body).Decode(&createResp)
		require.NoError(t, err)
		assert.NotZero(t, createResp.Id)
		userId = createResp.Id
	})

	t.Run("Get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", userId), nil)
		resp, err := server.app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var user models.User
		err = json.NewDecoder(resp.Body).Decode(&user)
		require.NoError(t, err)
		assert.Equal(t, 1, user.Id)
	})

	t.Run("Update", func(t *testing.T) {
		user := handlers.UpdateRequest{
			Email: "updated@example.com",
			Name:  "Updated User",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/users/%d", userId), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var updatedUser models.User
		err = json.NewDecoder(resp.Body).Decode(&updatedUser)
		require.NoError(t, err)
		assert.Equal(t, "updated@example.com", updatedUser.Email)
		assert.Equal(t, "Updated User", updatedUser.Name)
	})

	t.Run("UpdateEmail", func(t *testing.T) {
		email := "new-email@example.com"
		user := handlers.UpdateEmailRequest{
			Email: &email,
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%d", userId), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var updatedUser models.User
		err = json.NewDecoder(resp.Body).Decode(&updatedUser)
		require.NoError(t, err)
		assert.Equal(t, "new-email@example.com", updatedUser.Email)
	})

	t.Run("Delete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%d", userId), nil)
		resp, err := server.app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}
