package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/018bf/example/pkg/postgresql"
	"time"

	sq "github.com/Masterminds/squirrel"

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

func (r *PostgresUserRepository) Count(ctx context.Context, filter *models.UserFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.users")
	if filter.Search != nil {
		q = q.Where(postgresql.Search{
			Lang:   "english",
			Fields: []string{"first_name", "last_name", "email"},
			Query:  *filter.Search,
		})
	}
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

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.users").
		Columns(
			"first_name",
			"last_name",
			"password",
			"email",
			"group_id",
		).
		Values(
			user.FirstName,
			user.LastName,
			user.Password,
			user.Email,
			user.GroupID,
		).
		Suffix(`RETURNING id, created_at`)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(user); err != nil {
		e := errs.FromPostgresError(err)
		if errors.Is(e, errs.NewEntityNotFound()) {
			e = errs.NewInvalidFormError()
			e.AddParam("email", "The email field has already been taken.")
		}
		return e
	}
	return nil
}

func (r *PostgresUserRepository) Get(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	user := &models.User{}
	q := sq.Select(
		"id",
		"first_name",
		"last_name",
		"password",
		"email",
		"group_id",
		"created_at",
		"updated_at",
	).
		From("public.users").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, user, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		e.AddParam("user_id", fmt.Sprint(id))
		return nil, e
	}
	return user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	raw := &models.User{}
	q := sq.Select(
		"id",
		"first_name",
		"last_name",
		"password",
		"email",
		"group_id",
		"created_at",
		"updated_at",
	).
		From("public.users").
		Where(sq.Eq{"email": email}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, raw, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		e.AddParam("user_email", fmt.Sprint(email))
		return nil, e
	}
	return raw, nil
}

func (r *PostgresUserRepository) List(ctx context.Context, filter *models.UserFilter) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var users []*models.User
	const pageSize = 10
	q := sq.Select(
		"id",
		"first_name",
		"last_name",
		"password",
		"email",
		"group_id",
		"created_at",
		"updated_at",
	).
		From("public.users").
		Limit(pageSize)
	if filter.Search != nil {
		q = q.Where(postgresql.Search{
			Lang:   "english",
			Fields: []string{"first_name", "last_name", "email"},
			Query:  *filter.Search,
		})
	}
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
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return users, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.users").
		Where(sq.Eq{"id": user.ID}).
		Set("first_name", user.FirstName).
		Set("last_name", user.LastName).
		Set("password", user.Password).
		Set("email", user.Email).
		Set("group_id", user.GroupID)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			e = errs.NewInvalidFormError()
			e.AddParam("email", "The email field has already been taken.")
		}
		e.AddParam("user_id", fmt.Sprint(user.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err)
	}
	if affected == 0 {
		e := errs.NewEntityNotFound()
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
		e := errs.FromPostgresError(err)
		e.AddParam("user_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err)
		e.AddParam("user_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound()
		e.AddParam("user_id", fmt.Sprint(id))
		return e
	}
	return nil
}
