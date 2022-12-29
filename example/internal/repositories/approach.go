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

type PostgresApproachRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewPostgresApproachRepository(database *sqlx.DB, logger log.Logger) repositories.ApproachRepository {
	return &PostgresApproachRepository{database: database, logger: logger}
}

func (r *PostgresApproachRepository) Create(ctx context.Context, approach *models.Approach) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.approachs").
		Columns(). // TODO: add columns
		Values().  // TODO: add values
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(approach); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}

func (r *PostgresApproachRepository) Get(ctx context.Context, id string) (*models.Approach, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	approach := &models.Approach{}
	q := sq.Select("*").
		From("public.approachs").
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

func (r *PostgresApproachRepository) List(ctx context.Context, filter *models.ApproachFilter) ([]*models.Approach, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var approachs []*models.Approach
	const pageSize = 10
	q := sq.Select("*").
		From("public.approachs").
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
	if err := r.database.SelectContext(ctx, &approachs, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return approachs, nil
}

func (r *PostgresApproachRepository) Update(ctx context.Context, approach *models.Approach) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.approachs").Where(sq.Eq{"id": approach.ID}).Set("", "") // TODO: set values
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

func (r *PostgresApproachRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.approachs").Where(sq.Eq{"id": id})
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
