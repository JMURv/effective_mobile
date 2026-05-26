package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/JMURv/effective-mobile/internal/dto"
	md "github.com/JMURv/effective-mobile/internal/models"
	"github.com/JMURv/effective-mobile/internal/repo"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	os.Exit(m.Run())
}

func newTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:18.1-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = pgContainer.Terminate(ctx)
	})

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	conn, err := sqlx.Open("pgx", dsn)
	if err != nil {
		zap.L().Fatal("failed to connect to the database", zap.Error(err))
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		t.Fatal(err)
	}

	err = goose.Up(conn.DB, "../../../migrations")
	if err != nil {
		t.Fatal(err)
	}

	return conn
}

func newTestRepo(t *testing.T) *Repository {
	return &Repository{
		conn: newTestDB(t),
	}
}

func seedSubscription(t *testing.T, r *Repository) md.Subscription {
	t.Helper()

	req := dto.CreateSubscriptionRequest{
		ServiceName: "netflix",
		Price:       100,
		UserID:      uuid.New(),
		StartDate:   time.Now().Add(-time.Hour),
		EndDate:     nil,
	}

	res, err := r.Create(context.Background(), req)
	require.NoError(t, err)

	return res
}

func TestRepository_Create(t *testing.T) {
	r := newTestRepo(t)

	req := dto.CreateSubscriptionRequest{
		ServiceName: "netflix",
		Price:       200,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
		EndDate:     nil,
	}

	res, err := r.Create(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "netflix", res.ServiceName)
	assert.Equal(t, 200, res.Price)
	assert.NotZero(t, res.ID)
}

func TestRepository_GetByID(t *testing.T) {
	r := newTestRepo(t)

	sub := seedSubscription(t, r)

	t.Run("success", func(t *testing.T) {
		res, err := r.GetByID(context.Background(), sub.ID)

		require.NoError(t, err)
		assert.Equal(t, sub.ID, res.ID)
		assert.Equal(t, sub.ServiceName, res.ServiceName)
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := r.GetByID(context.Background(), uuid.New())
		assert.ErrorIs(t, err, repo.ErrNotFound)
	})
}

func TestRepository_List(t *testing.T) {
	r := newTestRepo(t)

	_ = seedSubscription(t, r)
	_ = seedSubscription(t, r)

	res, err := r.List(context.Background())

	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(res), 2)

	assert.True(t, res[0].CreatedAt.After(res[1].CreatedAt) ||
		res[0].CreatedAt.Equal(res[1].CreatedAt))
}

func TestRepository_Update(t *testing.T) {
	r := newTestRepo(t)

	sub := seedSubscription(t, r)

	req := dto.UpdateSubscriptionRequest{
		ServiceName: new("spotify"),
		Price:       new(999),
	}

	t.Run("success", func(t *testing.T) {
		res, err := r.Update(context.Background(), sub.ID, req)

		require.NoError(t, err)
		assert.Equal(t, "spotify", res.ServiceName)
		assert.Equal(t, 999, res.Price)
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := r.Update(context.Background(), uuid.New(), req)
		assert.ErrorIs(t, err, repo.ErrNotFound)
	})
}

func TestRepository_Delete(t *testing.T) {
	r := newTestRepo(t)

	sub := seedSubscription(t, r)

	err := r.Delete(context.Background(), sub.ID)
	require.NoError(t, err)

	_, err = r.GetByID(context.Background(), sub.ID)
	assert.ErrorIs(t, err, repo.ErrNotFound)
}

func TestRepository_CalculateTotalCost(t *testing.T) {
	r := newTestRepo(t)

	userID := uuid.New()

	_, _ = r.Create(context.Background(), dto.CreateSubscriptionRequest{
		ServiceName: "netflix",
		Price:       100,
		UserID:      userID,
		StartDate:   time.Now().Add(-2 * time.Hour),
		EndDate:     nil,
	})

	_, _ = r.Create(context.Background(), dto.CreateSubscriptionRequest{
		ServiceName: "netflix",
		Price:       200,
		UserID:      userID,
		StartDate:   time.Now().Add(-1 * time.Hour),
		EndDate:     nil,
	})

	req := dto.CalculateTotalCostRequest{
		ServiceName: "netflix",
		UserID:      userID,
		From:        time.Now().Add(-3 * time.Hour),
		To:          time.Now().Add(1 * time.Hour),
	}

	total, err := r.CalculateTotalCost(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, int64(300), total)
}
