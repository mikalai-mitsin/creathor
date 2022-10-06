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
	"github.com/lib/pq"
)

type PostgresUserRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewPostgresUserRepository(database *sqlx.DB, logger log.Logger) repositories.UserRepository {
	return &PostgresUserRepository{database: database, logger: logger}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.users").
		Columns(). // TODO: add columns
		Values().  // TODO: add values
		Suffix("RETURNING \"id\"")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).Scan(&user.ID); err != nil {
		e := errs.NewUnexpectedBehaviorError(err.Error())
		return e
	}
	return nil
}

func (r *PostgresUserRepository) Get(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	user := &models.User{}
	q := sq.Select("*").
		From("public.users").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, &user, query, args...); err != nil {
		e := errs.NewUnexpectedBehaviorError(err.Error())
		return nil, e
	}
	return user, nil
}

func (r *PostgresUserRepository) List(ctx context.Context, filter *models.UserFilter) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var users []*models.User
	const pageSize = 10
	q := sq.Select("*").
		From("public.users").
		Limit(pageSize) //
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
	if err := r.database.SelectContext(ctx, &users, query, args...); err != nil {
		e := errs.NewUnexpectedBehaviorError(err.Error())
		return nil, e
	}
	return users, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.users").Where(sq.Eq{"id": user.ID}).Set("", "") // TODO: set values
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.NewUnexpectedBehaviorError(err.Error())
		pgError, ok := err.(*pq.Error)
		if ok {
			switch pgError.Code {
			case "23505":
				e = errs.NewInvalidFormError()
				e.AddParam("phone", "The phone field has already been taken.")
			default:
				e = errs.NewUnexpectedBehaviorError(pgError.Detail)
			}
		}
		e.AddParam("user_id", fmt.Sprint(user.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.NewUnexpectedBehaviorError(err.Error())
	}
	if affected == 0 {
		e := errs.NewUserNotFound()
		e.AddParam("user_id", fmt.Sprint(user.ID))
		return e
	}
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.users").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.NewUnexpectedBehaviorError(err.Error())
		e.AddParam("user_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.NewUnexpectedBehaviorError(err.Error())
		e.AddParam("user_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewUserNotFound()
		e.AddParam("user_id", fmt.Sprint(id))
		return e
	}
	return nil
}
