package http

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/JMURv/effective-mobile/internal/cache/redis"
	"github.com/JMURv/effective-mobile/internal/config"
	"github.com/JMURv/effective-mobile/internal/ctrl"
	hdl "github.com/JMURv/effective-mobile/internal/hdl/http"
	"github.com/JMURv/effective-mobile/internal/hdl/validation"
	"github.com/JMURv/effective-mobile/internal/repo/db"
	"github.com/moby/moby/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var rootDir = filepath.Join("..", "..", "..")

func getRedis(ctx context.Context) testcontainers.Container {
	req := testcontainers.ContainerRequest{
		Image:        "redis:8.0.3-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	return c
}

func getPostgres(ctx context.Context) testcontainers.Container {
	pgPort := os.Getenv("POSTGRES_PORT")
	pgPortC := fmt.Sprintf("%s/tcp", pgPort)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:18.1-alpine",
		WaitingFor:   wait.ForHealthCheck(),
		ExposedPorts: []string{pgPortC},
		ConfigModifier: func(conf *container.Config) {
			conf.Healthcheck = &container.HealthConfig{
				Test: []string{
					"CMD-SHELL",
					fmt.Sprintf(
						"pg_isready -U %s -d %s",
						os.Getenv("POSTGRES_USER"),
						os.Getenv("POSTGRES_DB"),
					),
				},
				Interval:    5 * time.Second,
				Timeout:     2 * time.Second,
				Retries:     5,
				StartPeriod: 2 * time.Second,
			}
		},
		Env: map[string]string{
			"POSTGRES_DB":       os.Getenv("POSTGRES_DB"),
			"POSTGRES_USER":     os.Getenv("POSTGRES_USER"),
			"POSTGRES_PASSWORD": os.Getenv("POSTGRES_PASSWORD"),
			"POSTGRES_HOST":     os.Getenv("POSTGRES_HOST"),
			"POSTGRES_PORT":     os.Getenv("POSTGRES_PORT"),
		},
	}

	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	return pgC
}

func getServer(conf config.Config) *httptest.Server {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))

	svc := ctrl.New(db.New(conf), redis.New(conf))
	h := hdl.New(svc)
	return httptest.NewServer(h.Router)
}

type TestEnv struct {
	Cache    testcontainers.Container
	Postgres testcontainers.Container
	Server   *httptest.Server
}

func NewTestEnv(t *testing.T) *TestEnv {
	t.Helper()

	ctx := context.Background()

	t.Setenv("MIGRATIONS_PATH", filepath.ToSlash(
		filepath.Join(rootDir, "migrations"),
	))

	validation.New()

	conf := config.MustLoad(
		filepath.ToSlash(
			filepath.Join(rootDir, "build", "configs", "envs", ".env.integration"),
		),
	)

	env := &TestEnv{}

	errg := errgroup.Group{}

	errg.Go(func() error {
		pg := getPostgres(ctx)
		env.Postgres = pg

		pgHost, err := pg.Host(ctx)
		if err != nil {
			zap.L().Error("failed to get pg host", zap.Error(err))
			return err
		}

		pgPort, err := pg.MappedPort(ctx, "5432/tcp")
		if err != nil {
			zap.L().Error("failed to get pg port", zap.Error(err))
			return err
		}

		conf.DB.Host = pgHost
		conf.DB.Port = int(pgPort.Num())
		return nil
	})

	errg.Go(func() error {
		cache := getRedis(ctx)
		env.Cache = cache

		cacheHost, err := cache.Host(ctx)
		if err != nil {
			zap.L().Error("failed to get cache host", zap.Error(err))
			return err
		}

		cachePort, err := cache.MappedPort(ctx, "6379/tcp")
		if err != nil {
			zap.L().Error("failed to get cache port", zap.Error(err))
			return err
		}

		conf.Redis.Addr = fmt.Sprintf("%s:%s", cacheHost, cachePort.Port())

		return nil
	})

	if err := errg.Wait(); err != nil {
		t.Fatal(err)
	}

	env.Server = getServer(conf)

	t.Cleanup(func() {
		env.Teardown(ctx, t)
	})

	return env
}

func (e *TestEnv) Teardown(ctx context.Context, t *testing.T) {
	t.Helper()

	if e.Server != nil {
		e.Server.Close()
	}

	errg := errgroup.Group{}

	errg.Go(func() error {
		if e.Cache != nil {
			if err := e.Cache.Terminate(ctx); err != nil {
				zap.L().Error("failed to terminate cache container", zap.Error(err))
				return err
			}
		}

		return nil
	})

	errg.Go(func() error {
		if e.Postgres != nil {
			if err := e.Postgres.Terminate(ctx); err != nil {
				zap.L().Error("failed to terminate postgres container", zap.Error(err))
				return err
			}
		}

		return nil
	})

	if err := errg.Wait(); err != nil {
		t.Fatal(err)
	}
}
