package db

import (
	"context"

	"github.com/JMURv/golang-clean-template/internal/dto"
	md "github.com/JMURv/golang-clean-template/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (r *Repository) Create(
	ctx context.Context,
	req dto.CreateSubscriptionRequest,
) (md.Subscription, error) {
	const op = "sub.Create.repo"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := builder.
		Insert("subscriptions").
		Columns(
			"service_name",
			"price",
			"user_id",
			"start_date",
			"end_date",
		).
		Values(
			req.ServiceName,
			req.Price,
			req.UserID,
			req.StartDate,
			req.EndDate,
		).
		Suffix(`
			RETURNING
				id,
				service_name,
				price,
				user_id,
				start_date,
				end_date,
				created_at,
				updated_at
		`)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return md.Subscription{}, err
	}

	var res md.Subscription
	err = r.conn.GetContext(
		ctx,
		&res,
		sqlQuery,
		args...,
	)
	if err != nil {
		zap.L().Error(
			"failed to create subscription",
			zap.String("op", op),
			zap.Error(err),
		)

		return md.Subscription{}, err
	}

	return res, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (md.Subscription, error) {
	const op = "sub.GetByID.repo"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := builder.
		Select(
			"id",
			"service_name",
			"price",
			"user_id",
			"start_date",
			"end_date",
			"created_at",
			"updated_at",
		).
		From("subscriptions").
		Where(sq.Eq{"id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return md.Subscription{}, err
	}

	res := md.Subscription{}
	err = r.conn.GetContext(ctx, &res, sqlQuery, args...)
	if err != nil {
		return md.Subscription{}, err
	}

	return res, nil
}

func (r *Repository) List(ctx context.Context) ([]md.Subscription, error) {
	const op = "sub.List.repo"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := builder.
		Select(
			"id",
			"service_name",
			"price",
			"user_id",
			"start_date",
			"end_date",
			"created_at",
			"updated_at",
		).
		From("subscriptions").
		OrderBy("created_at DESC")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var res []md.Subscription
	err = r.conn.SelectContext(ctx, &res, sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) Update(
	ctx context.Context,
	id uuid.UUID,
	req dto.UpdateSubscriptionRequest,
) (md.Subscription, error) {
	const op = "sub.Update.repo"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	builder := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("subscriptions").
		Where(sq.Eq{"id": id}).
		Set("updated_at", sq.Expr("NOW()"))

	if req.ServiceName != nil {
		builder = builder.Set("service_name", *req.ServiceName)
	}

	if req.Price != nil {
		builder = builder.Set("price", *req.Price)
	}

	if req.StartDate != nil {
		builder = builder.Set("start_date", *req.StartDate)
	}

	if req.EndDate != nil {
		builder = builder.Set("end_date", *req.EndDate)
	}

	builder = builder.Suffix(`
		RETURNING
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date,
			created_at,
			updated_at
	`)

	sqlQuery, args, err := builder.ToSql()
	if err != nil {
		return md.Subscription{}, err
	}

	res := md.Subscription{}
	err = r.conn.GetContext(ctx, &res, sqlQuery, args...)
	if err != nil {
		return md.Subscription{}, err
	}

	return res, nil
}

func (r *Repository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "sub.Delete.repo"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Delete("subscriptions").
		Where(sq.Eq{"id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.conn.ExecContext(ctx, sqlQuery, args...)
	return err
}

func (r *Repository) CalculateTotalCost(
	ctx context.Context,
	req dto.CalculateTotalCostRequest,
) (int64, error) {
	const op = "sub.CalculateTotalCost.repo"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("COALESCE(SUM(price), 0)").
		From("subscriptions").
		Where(sq.Eq{
			"user_id":      req.UserID,
			"service_name": req.ServiceName,
		}).
		Where(sq.LtOrEq{
			"start_date": req.To,
		}).
		Where(
			sq.Or{
				sq.Expr("end_date IS NULL"),
				sq.GtOrEq{
					"end_date": req.From,
				},
			},
		)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	var total int64
	err = r.conn.GetContext(ctx, &total, sqlQuery, args...)
	if err != nil {
		return 0, err
	}

	return total, nil
}
