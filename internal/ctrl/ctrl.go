package ctrl

import (
	"context"
	"time"
)

type AppRepo interface{ subscriptionRepo }

type AppCtrl interface{ subscriptionCtrl }

type CacheService interface {
	Close(ctx context.Context) error
	GetToStruct(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, t time.Duration, key string, val any)
	Delete(ctx context.Context, key string)
}

type Controller struct {
	repo  AppRepo
	cache CacheService
}

func New(repo AppRepo, cache CacheService) *Controller {
	return &Controller{
		repo:  repo,
		cache: cache,
	}
}
