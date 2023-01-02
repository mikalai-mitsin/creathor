package repositories

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"time"

	"github.com/018bf/creathor/pkg/log"

	"github.com/018bf/creathor/internal/domain/models"
	"github.com/018bf/creathor/internal/domain/repositories"

	"github.com/018bf/creathor/internal/domain/errs"
	"github.com/jmoiron/sqlx"
)

type UserApproachRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewUserApproachRepository(database *sqlx.DB, logger log.Logger) repositories.UserApproachRepository {
	return &UserApproachRepository{database: database, logger: logger}
}

func (r *UserApproachRepository) Create(ctx context.Context, userApproach *models.UserApproach) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.user_approaches").
		Columns(). // TODO: add columns
		Values().  // TODO: add values
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(userApproach); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}

func (r *UserApproachRepository) Get(ctx context.Context, id string) (*models.UserApproach, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	userApproach := &models.UserApproach{}
	q := sq.Select("user_approaches.id", "user_approaches.updated_at", "user_approaches.created_at").
		From("public.user_approaches").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, userApproach, query, args...); err != nil {
		e := errs.FromPostgresError(err).
			WithParam("user_approach_id", id)
		return nil, e
	}
	return userApproach, nil
}

func (r *UserApproachRepository) List(ctx context.Context, filter *models.UserApproachFilter) ([]*models.UserApproach, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var userApproaches []*models.UserApproach
	const pageSize = 10
	q := sq.Select("user_approaches.id", "user_approaches.updated_at", "user_approaches.created_at").
		From("public.user_approaches").
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
	if err := r.database.SelectContext(ctx, &userApproaches, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return userApproaches, nil
}

func (r *UserApproachRepository) Update(ctx context.Context, userApproach *models.UserApproach) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.user_approaches").Where(sq.Eq{"id": userApproach.ID}).Set("", "") // TODO: set values
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("user_approach_id", fmt.Sprint(userApproach.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).
			WithParam("user_approach_id", fmt.Sprint(userApproach.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("user_approach_id", fmt.Sprint(userApproach.ID))
		return e
	}
	return nil
}

func (r *UserApproachRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.user_approaches").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("user_approach_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("user_approach_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("user_approach_id", fmt.Sprint(id))
		return e
	}
	return nil
}

func (r *UserApproachRepository) Count(ctx context.Context, filter *models.UserApproachFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.user_approaches")
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
