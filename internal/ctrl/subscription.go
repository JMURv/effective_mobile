package ctrl

import (
	"context"
	"fmt"

	"github.com/JMURv/golang-clean-template/internal/cache"
	"github.com/JMURv/golang-clean-template/internal/config"
	"github.com/JMURv/golang-clean-template/internal/dto"
	md "github.com/JMURv/golang-clean-template/internal/models"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

type subscriptionCtrl interface {
	Create(ctx context.Context, req dto.CreateSubscriptionRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (md.Subscription, error)
	List(ctx context.Context) ([]md.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateSubscriptionRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	CalculateTotalCost(
		ctx context.Context,
		req dto.CalculateTotalCostRequest,
	) (dto.CalculateTotalCostResponse, error)
}

type subscriptionRepo interface {
	Create(ctx context.Context, req dto.CreateSubscriptionRequest) (md.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (md.Subscription, error)
	List(ctx context.Context) ([]md.Subscription, error)
	Update(
		ctx context.Context,
		id uuid.UUID,
		req dto.UpdateSubscriptionRequest,
	) (md.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CalculateTotalCost(
		ctx context.Context,
		req dto.CalculateTotalCostRequest,
	) (int64, error)
}

func (c *Controller) Create(ctx context.Context, req dto.CreateSubscriptionRequest) error {
	const op = "sub.Create.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := c.repo.Create(ctx, req)
	if err != nil {
		return err
	}

	go func() {
		c.cache.Set(ctx, config.DefaultCacheTime, fmt.Sprintf(cache.SubKey, res.ID), res)
		c.cache.Delete(ctx, cache.ListSubKey)
	}()
	return nil
}

func (c *Controller) GetByID(ctx context.Context, id uuid.UUID) (md.Subscription, error) {
	const op = "sub.GetByID.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := md.Subscription{}
	if err := c.cache.GetToStruct(ctx, fmt.Sprintf(cache.SubKey, id), &res); err == nil {
		return res, nil
	}

	res, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return res, err
	}

	go c.cache.Set(ctx, config.DefaultCacheTime, fmt.Sprintf(cache.SubKey, res.ID), res)
	return res, nil
}

func (c *Controller) List(ctx context.Context) ([]md.Subscription, error) {
	const op = "sub.List.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := make([]md.Subscription, 0, config.DefaultSize)
	if err := c.cache.GetToStruct(ctx, cache.ListSubKey, res); err == nil {
		return res, nil
	}

	res, err := c.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	go c.cache.Set(ctx, config.DefaultCacheTime, cache.ListSubKey, res)
	return res, nil
}

func (c *Controller) Update(
	ctx context.Context,
	id uuid.UUID,
	req dto.UpdateSubscriptionRequest,
) error {
	const op = "sub.Update.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := c.repo.Update(ctx, id, req)
	if err != nil {
		return err
	}

	go c.cache.Set(ctx, config.DefaultCacheTime, fmt.Sprintf(cache.SubKey, id), res)
	return nil
}

func (c *Controller) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "sub.Delete.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	err := c.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	go c.cache.Delete(ctx, cache.ListSubKey)
	return nil
}

func (c *Controller) CalculateTotalCost(
	ctx context.Context,
	req dto.CalculateTotalCostRequest,
) (dto.CalculateTotalCostResponse, error) {
	const op = "sub.CalculateTotalCost.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := c.repo.CalculateTotalCost(ctx, req)
	if err != nil {
		return dto.CalculateTotalCostResponse{}, err
	}

	return dto.CalculateTotalCostResponse{
		TotalCost: res,
	}, nil
}
