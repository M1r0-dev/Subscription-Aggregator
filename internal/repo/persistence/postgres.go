package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/M1r0-dev/Subscription-Aggregator/internal/entity"
	"github.com/M1r0-dev/Subscription-Aggregator/pkg/postgres"
	"github.com/Masterminds/squirrel"
)

type SubscriptionRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *SubscriptionRepo {
	return &SubscriptionRepo{
		pg,
	}
}

func (r *SubscriptionRepo) Store(ctx context.Context, sub *entity.Subscription) error {
	const op = "subscriptionRepo.Store"
	sql, args, err := r.Builder.
		Insert("subscriptions").
		Columns("service_name", "price", "user_id", "start_date", "end_date").
		Values(sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: build query: %w", op, err)
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&sub.Id)
	if err != nil {
		return fmt.Errorf("%s: execute query %w", op, err)
	}

	return nil
}

func (r *SubscriptionRepo) Get(ctx context.Context, id int) (*entity.Subscription, error) {
	const op = "subscriptionRepo.Get"
	sql, args, err := r.Builder.
		Select("id", "service_name", "price", "user_id", "start_date", "end_date").
		From("subscriptions").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: build query: %w", op, err)
	}
	sub := &entity.Subscription{}
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&sub.Id, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	if err != nil {
		return nil, fmt.Errorf("%s: execute query %w", op, err)
	}

	return sub, nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, sub *entity.Subscription) error {
	const op = "subscriptionRepo.Update"

	sql, args, err := r.Builder.
		Update("subscriptions").
		Set("service_name", sub.ServiceName).
		Set("price", sub.Price).
		Set("user_id", sub.UserID).
		Set("start_date", sub.StartDate).
		Set("end_date", sub.EndDate).
		Where(squirrel.Eq{"id": sub.Id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: build query: %w", op, err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: execute query: %w", op, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("%s: no rows affected", op)
	}

	return nil
}

func (r *SubscriptionRepo) Delete(ctx context.Context, id int) error {
	const op = "subscriptionRepo.Delete"

	sql, args, err := r.Builder.
		Delete("subscriptions").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: build query: %w", op, err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: execute query: %w", op, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("%s: subscription not found", op)
	}

	return nil
}

func (r *SubscriptionRepo) List(ctx context.Context, opts ...ListOption) ([]*entity.Subscription, error) {
	const op = "subscriptionRepo.List"
	options := &ListOptions{
		Limit:     50,
		SortBy:    "start_date",
		SortOrder: "desc",
	}

	for _, opt := range opts {
		opt(options)
	}

	builder := r.Builder.
		Select("id", "service_name", "price", "user_id", "start_date", "end_date").
		From("subscriptions")

	if options.UserID != nil {
		builder = builder.Where(squirrel.Eq{"user_id": *options.UserID})
	}

	if options.ServiceName != nil {
		builder = builder.Where(squirrel.Eq{"service_name": *options.ServiceName})
	}

	if options.Price != nil {
		builder = builder.Where(squirrel.Eq{"price": *options.Price})
	}

	if options.StartDateFrom != nil {
		builder = builder.Where(squirrel.GtOrEq{"start_date": *options.StartDateFrom})
	}

	if options.StartDateTo != nil {
		builder = builder.Where(squirrel.LtOrEq{"start_date": *options.StartDateTo})
	}

	builder = builder.OrderBy(fmt.Sprintf("%s %s", options.SortBy, options.SortOrder))

	if options.Limit > 0 {
		builder = builder.Limit(uint64(options.Limit))
	}

	if options.Offset > 0 {
		builder = builder.Offset(uint64(options.Offset))
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build query: %w", op, err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: execute query: %w", op, err)
	}
	defer rows.Close()

	var subscriptions []*entity.Subscription
	for rows.Next() {
		var sub entity.Subscription
		err := rows.Scan(
			&sub.Id, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		subscriptions = append(subscriptions, &sub)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepo) Count(ctx context.Context, opts ...ListOption) (int, error) {
	const op = "subscriptionRepo.Count"

	options := &ListOptions{}
	for _, opt := range opts {
		opt(options)
	}

	builder := r.Builder.Select("COUNT(*)").
		From("subscriptions")

	if options.UserID != nil {
		builder = builder.Where(squirrel.Eq{"user_id": *options.UserID})
	}
	if options.ServiceName != nil {
		builder = builder.Where(squirrel.Eq{"service_name": *options.ServiceName})
	}
	if options.Price != nil {
		builder = builder.Where(squirrel.Eq{"price": *options.Price})
	}
	if options.StartDateFrom != nil {
		builder = builder.Where(squirrel.GtOrEq{"start_date": *options.StartDateFrom})
	}
	if options.StartDateTo != nil {
		builder = builder.Where(squirrel.LtOrEq{"start_date": *options.StartDateTo})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: build query: %w", op, err)
	}

	var count int
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return count, nil
}

func (r *SubscriptionRepo) GetTotalCost(ctx context.Context, userID *string, serviceName *string, startDate, endDate string) (uint64, error) {
	const op = "subscriptionRepo.GetTotalCost"

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return 0, fmt.Errorf("%s: invalid start date format: %w", op, err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return 0, fmt.Errorf("%s: invalid end date format: %w", op, err)
	}

	builder := r.Builder.
		Select("COALESCE(SUM(price), 0)").
		From("subscriptions").
		Where("start_date <= ?", end).                                                            
		Where("(end_date >= ? OR end_date IS NULL OR end_date = '0001-01-01'::timestamp)", start)

	if userID != nil && *userID != "" {
		builder = builder.Where(squirrel.Eq{"user_id": *userID})
	}

	if serviceName != nil && *serviceName != "" {
		builder = builder.Where(squirrel.Eq{"service_name": *serviceName})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: build query: %w", op, err)
	}

	var total uint64
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return total, nil
}
