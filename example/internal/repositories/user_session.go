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

type UserSessionRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewUserSessionRepository(
	database *sqlx.DB,
	logger log.Logger,
) repositories.UserSessionRepository {
	return &UserSessionRepository{
		database: database,
		logger:   logger,
	}
}

func (r *UserSessionRepository) Create(
	ctx context.Context,
	userSession *models.UserSession,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.user_sessions").
		Columns(
			"updated_at",
			"created_at",
		).
		Values(
			userSession.UpdatedAt,
			userSession.CreatedAt,
		).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(userSession); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}

func (r *UserSessionRepository) Get(
	ctx context.Context,
	id string,
) (*models.UserSession, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	userSession := &models.UserSession{}
	q := sq.Select(
		"user_sessions.id",
		"user_sessions.updated_at",
		"user_sessions.created_at",
	).
		From("public.user_sessions").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, userSession, query, args...); err != nil {
		e := errs.FromPostgresError(err).
			WithParam("user_session_id", id)
		return nil, e
	}
	return userSession, nil
}

func (r *UserSessionRepository) List(
	ctx context.Context,
	filter *models.UserSessionFilter,
) ([]*models.UserSession, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var userSessions []*models.UserSession
	const pageSize = 10
	q := sq.Select(
		"user_sessions.id",
		"user_sessions.updated_at",
		"user_sessions.created_at",
	).
		From("public.user_sessions").
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
	if err := r.database.SelectContext(ctx, &userSessions, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return userSessions, nil
}

func (r *UserSessionRepository) Update(
	ctx context.Context,
	userSession *models.UserSession,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.user_sessions").
		Where(sq.Eq{"id": userSession.ID}).
		Set("updated_at", userSession.UpdatedAt)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("user_session_id", fmt.Sprint(userSession.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).
			WithParam("user_session_id", fmt.Sprint(userSession.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("user_session_id", fmt.Sprint(userSession.ID))
		return e
	}
	return nil
}

func (r *UserSessionRepository) Delete(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.user_sessions").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("user_session_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("user_session_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("user_session_id", fmt.Sprint(id))
		return e
	}
	return nil
}

func (r *UserSessionRepository) Count(
	ctx context.Context,
	filter *models.UserSessionFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.user_sessions")
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
