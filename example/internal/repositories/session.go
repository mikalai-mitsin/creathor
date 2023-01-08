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

type SessionRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewSessionRepository(
	database *sqlx.DB,
	logger log.Logger,
) repositories.SessionRepository {
	return &SessionRepository{
		database: database,
		logger:   logger,
	}
}

func (r *SessionRepository) Create(
	ctx context.Context,
	session *models.Session,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.sessions").
		Columns(
			"updated_at",
			"created_at",
		).
		Values(
			session.UpdatedAt,
			session.CreatedAt,
		).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(session); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}

func (r *SessionRepository) Get(
	ctx context.Context,
	id string,
) (*models.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	session := &models.Session{}
	q := sq.Select(
		"sessions.id",
		"sessions.updated_at",
		"sessions.created_at",
	).
		From("public.sessions").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, session, query, args...); err != nil {
		e := errs.FromPostgresError(err).
			WithParam("session_id", id)
		return nil, e
	}
	return session, nil
}

func (r *SessionRepository) List(
	ctx context.Context,
	filter *models.SessionFilter,
) ([]*models.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var sessions []*models.Session
	const pageSize = 10
	q := sq.Select(
		"sessions.id",
		"sessions.updated_at",
		"sessions.created_at",
	).
		From("public.sessions").
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
	if err := r.database.SelectContext(ctx, &sessions, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return sessions, nil
}

func (r *SessionRepository) Update(
	ctx context.Context,
	session *models.Session,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.sessions").
		Where(sq.Eq{"id": session.ID}).
		Set("updated_at", session.UpdatedAt)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("session_id", fmt.Sprint(session.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).
			WithParam("session_id", fmt.Sprint(session.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("session_id", fmt.Sprint(session.ID))
		return e
	}
	return nil
}

func (r *SessionRepository) Delete(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.sessions").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("session_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("session_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("session_id", fmt.Sprint(id))
		return e
	}
	return nil
}

func (r *SessionRepository) Count(
	ctx context.Context,
	filter *models.SessionFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.sessions")
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
