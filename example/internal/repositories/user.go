package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/018bf/example/pkg/log"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewUserRepository(
	database *sqlx.DB,
	logger log.Logger,
) repositories.UserRepository {
	return &UserRepository{
		database: database,
		logger:   logger,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *models.User,
) *errs.Error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.users").
		Columns(
			"updated_at",
			"created_at",
		).
		Values(
			user.UpdatedAt,
			user.CreatedAt,
		).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(user); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}

func (r *UserRepository) Get(
	ctx context.Context,
	id string,
) (*models.User, *errs.Error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	user := &models.User{}
	q := sq.Select(
		"users.id",
		"users.updated_at",
		"users.created_at",
	).
		From("public.users").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, user, query, args...); err != nil {
		e := errs.FromPostgresError(err).
			WithParam("user_id", id)
		return nil, e
	}
	return user, nil
}

func (r *UserRepository) List(
	ctx context.Context,
	filter *models.UserFilter,
) ([]*models.User, *errs.Error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var users []*models.User
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = utils.Pointer(pageSize)
	}
	q := sq.Select(
		"users.id",
		"users.updated_at",
		"users.created_at",
	).
		From("public.users").
		Limit(pageSize)
	// TODO: add filtering
	if filter.PageNumber != nil && *filter.PageNumber > 1 {
		q = q.Offset((*filter.PageNumber - 1) * *filter.PageSize)
	}
	q = q.Limit(*filter.PageSize)
	if len(filter.OrderBy) > 0 {
		q = q.OrderBy(filter.OrderBy...)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.SelectContext(ctx, &users, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return users, nil
}

func (r *UserRepository) Update(
	ctx context.Context,
	user *models.User,
) *errs.Error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.users").
		Where(sq.Eq{"id": user.ID}).
		Set("updated_at", user.UpdatedAt)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("user_id", fmt.Sprint(user.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).
			WithParam("user_id", fmt.Sprint(user.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("user_id", fmt.Sprint(user.ID))
		return e
	}
	return nil
}

func (r *UserRepository) Delete(
	ctx context.Context,
	id string,
) *errs.Error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.users").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("user_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("user_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("user_id", fmt.Sprint(id))
		return e
	}
	return nil
}

func (r *UserRepository) Count(
	ctx context.Context,
	filter *models.UserFilter,
) (uint64, *errs.Error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.users")
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
