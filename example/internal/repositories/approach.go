package repositories

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"time"

	"github.com/018bf/example/pkg/log"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/jmoiron/sqlx"
)

type ApproachRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewApproachRepository(
	database *sqlx.DB,
	logger log.Logger,
) repositories.ApproachRepository {
	return &ApproachRepository{
		database: database,
		logger:   logger,
	}
}

func (r *ApproachRepository) Create(
	ctx context.Context,
	approach *models.Approach,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.approaches").
		Columns(
			"updated_at",
			"created_at",
		).
		Values(
			approach.UpdatedAt,
			approach.CreatedAt,
		).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(approach); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}

func (r *ApproachRepository) Get(
	ctx context.Context,
	id string,
) (*models.Approach, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	approach := &models.Approach{}
	q := sq.Select(
		"approaches.id",
		"approaches.updated_at",
		"approaches.created_at",
	).
		From("public.approaches").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, approach, query, args...); err != nil {
		e := errs.FromPostgresError(err).
			WithParam("approach_id", id)
		return nil, e
	}
	return approach, nil
}

func (r *ApproachRepository) List(
	ctx context.Context,
	filter *models.ApproachFilter,
) ([]*models.Approach, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var approaches []*models.Approach
	const pageSize = 10
	q := sq.Select(
		"approaches.id",
		"approaches.updated_at",
		"approaches.created_at",
	).
		From("public.approaches").
		Limit(pageSize)
	// TODO: add filtering
	if filter.PageNumber != nil && *filter.PageNumber > 1 {
		q = q.Offset((*filter.PageNumber - 1) * *filter.PageSize)
	}
	if filter.PageSize != nil {
		q = q.Limit(*filter.PageSize)
	}
	if len(filter.OrderBy) > 0 {
		q = q.OrderBy(filter.OrderBy...)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.SelectContext(ctx, &approaches, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return approaches, nil
}

func (r *ApproachRepository) Update(
	ctx context.Context,
	approach *models.Approach,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.approaches").
		Where(sq.Eq{"id": approach.ID}).
		Set("updated_at", approach.UpdatedAt)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("approach_id", fmt.Sprint(approach.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).
			WithParam("approach_id", fmt.Sprint(approach.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("approach_id", fmt.Sprint(approach.ID))
		return e
	}
	return nil
}

func (r *ApproachRepository) Delete(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.approaches").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("approach_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("approach_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("approach_id", fmt.Sprint(id))
		return e
	}
	return nil
}

func (r *ApproachRepository) Count(
	ctx context.Context,
	filter *models.ApproachFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.approaches")
	// TODO: add filtering
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result := r.database.QueryRowxContext(ctx, query, args...)
	if err := result.Err(); err != nil {
		e := errs.FromPostgresError(err)
		return 0, e
	}
	var count uint64
	if err := result.Scan(&count); err != nil {
		e := errs.FromPostgresError(err)
		return 0, e
	}
	return count, nil
}
