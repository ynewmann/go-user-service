package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"go-user-service/src/repository"
	"go-user-service/src/repository/models"
)

var testDB *sql.DB

func TestPostgres(t *testing.T) {
	ctx := context.Background()

	testUser := "testuser"
	testPass := "testpass"
	testDb := "testdb"

	pgContainer, err := testpostgres.Run(
		ctx,
		"postgres:alpine",
		testpostgres.WithDatabase(testDb),
		testpostgres.WithUsername(testUser),
		testpostgres.WithPassword(testPass),
		testpostgres.BasicWaitStrategies(),
		testcontainers.WithLogger(testcontainers.Logger),
	)
	require.NoError(t, err, err)
	dur := time.Second
	defer pgContainer.Stop(ctx, &dur)

	dbHost, err := pgContainer.Host(ctx)
	require.NoError(t, err, err)

	dbPortStr, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Could not get mapped port: %v", err)
	}

	db, err := NewRepository(repository.Config{
		Host:     dbHost,
		Port:     dbPortStr.Int(),
		User:     testUser,
		Password: testPass,
		DbName:   testDb,
		SslMode:  "disable",
	})
	require.NoError(t, err, err)

	t.Run("Create", func(t *testing.T) {
		newValidUser := &models.User{
			Id:    0,
			Email: "test@email.com",
			Name:  "testname",
		}

		t.Run("Valid", func(t *testing.T) {
			id, err := db.Create(ctx, newValidUser)
			require.NoError(t, err, err)
			require.Greater(t, id, 0)
		})

		t.Run("SameEmail", func(t *testing.T) {
			id, err := db.Create(ctx, newValidUser)
			require.Error(t, err, err)
			require.Equal(t, id, 0)
		})

		t.Run("EmptyName", func(t *testing.T) {
			user := &models.User{
				Email: "test2@email.com",
			}
			id, err := db.Create(ctx, user)
			require.NoError(t, err, err)
			require.Greater(t, id, 0)
		})

		t.Run("EmptyEmail", func(t *testing.T) {
			user := newValidUser
			user.Email = ""
			id, err := db.Create(ctx, user)
			require.NoError(t, err, err)
			require.Greater(t, id, 0)
		})
	})

	t.Run("Get", func(t *testing.T) {
		validUser := &models.User{
			Id:    0,
			Email: "toget@email.com",
			Name:  "toget",
		}

		id, err := db.Create(ctx, validUser)
		require.NoError(t, err, err)
		require.Greater(t, id, 0)
		validUser.Id = id

		t.Run("Valid", func(t *testing.T) {
			u, err := db.Get(ctx, id)
			require.NoError(t, err, err)
			require.Equal(t, validUser.Id, u.Id)
			require.Equal(t, validUser.Name, u.Name)
			require.Equal(t, validUser.Email, u.Email)
		})

		t.Run("BadId", func(t *testing.T) {
			u, err := db.Get(ctx, 0)
			require.True(t, errors.Is(err, ErrNotFound))
			require.Nil(t, u)
		})
	})

	t.Run("Update", func(t *testing.T) {
		validUser := &models.User{
			Id:    0,
			Email: "toupdate@email.com",
			Name:  "toupdate",
		}

		id, err := db.Create(ctx, validUser)
		require.NoError(t, err, err)
		require.Greater(t, id, 0)
		validUser.Id = id

		newEmail := "new@email.com"
		t.Run("ValidEmail", func(t *testing.T) {
			err := db.UpdateEmail(ctx, validUser.Id, newEmail)
			require.NoError(t, err, err)

			u, err := db.Get(ctx, validUser.Id)
			require.Equal(t, validUser.Id, u.Id)
			require.Equal(t, validUser.Name, u.Name)
			require.Equal(t, newEmail, u.Email)
			validUser.Email = newEmail
		})

		t.Run("BadIdEmail", func(t *testing.T) {
			err := db.UpdateEmail(ctx, 0, newEmail)
			require.Error(t, err, ErrNotFound.Error())
		})

		updatedUser := &models.User{
			Id:    validUser.Id,
			Email: "updated@email.com",
			Name:  "updated",
		}

		t.Run("Valid", func(t *testing.T) {
			err := db.Update(ctx, updatedUser)
			require.NoError(t, err, err)

			u, err := db.Get(ctx, updatedUser.Id)
			require.Equal(t, updatedUser.Id, u.Id)
			require.Equal(t, updatedUser.Name, u.Name)
			require.Equal(t, updatedUser.Email, u.Email)
		})

		t.Run("BadIdEmail", func(t *testing.T) {
			err := db.UpdateEmail(ctx, 0, newEmail)
			require.Error(t, err, ErrNotFound.Error())
		})
	})

	t.Run("Delete", func(t *testing.T) {
		validUser := &models.User{
			Id:    0,
			Email: "todelete@email.com",
			Name:  "todelete",
		}

		id, err := db.Create(ctx, validUser)
		require.NoError(t, err, err)
		require.Greater(t, id, 0)
		validUser.Id = id

		t.Run("BadId", func(t *testing.T) {
			err := db.Delete(ctx, 0)
			require.True(t, errors.Is(err, ErrNotFound))
		})

		t.Run("Valid", func(t *testing.T) {
			err := db.Delete(ctx, id)
			require.NoError(t, err, err)

			u, err := db.Get(ctx, id)
			require.True(t, errors.Is(err, ErrNotFound))
			require.Nil(t, u)
		})
	})
}
