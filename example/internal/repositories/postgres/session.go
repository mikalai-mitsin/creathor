package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/pkg/log"
	"github.com/018bf/example/pkg/postgresql"
	"github.com/018bf/example/pkg/utils"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewSessionRepository(database *sqlx.DB, logger log.Logger) repositories.SessionRepository {
	return &SessionRepository{database: database, logger: logger}
}
func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewSessionDTOFromModel(session)
	q := sq.Insert("public.sessions").
		Columns("updated_at", "created_at", "title", "description").
		Values(dto.UpdatedAt, dto.CreatedAt, dto.Title, dto.Description).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(dto); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	session.ID = models.UUID(dto.ID)
	return nil
}
func (r *SessionRepository) Get(ctx context.Context, id models.UUID) (*models.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &SessionDTO{}
	q := sq.Select("sessions.id", "sessions.updated_at", "sessions.created_at", "sessions.title", "sessions.description").
		From("public.sessions").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("session_id", string(id))
		return nil, e
	}
	return dto.ToModel(), nil
}

func (r *SessionRepository) List(
	ctx context.Context,
	filter *models.SessionFilter,
) ([]*models.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto SessionListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = utils.Pointer(pageSize)
	}
	q := sq.Select("sessions.id", "sessions.updated_at", "sessions.created_at", "sessions.title", "sessions.description").
		From("public.sessions").
		Limit(pageSize)
	if filter.Search != nil {
		q = q.Where(
			postgresql.Search{
				Lang:   "english",
				Query:  *filter.Search,
				Fields: []string{"description"},
			},
		)
	}
	if len(filter.IDs) > 0 {
		q = q.Where(sq.Eq{"id": filter.IDs})
	}
	if filter.PageNumber != nil && *filter.PageNumber > 1 {
		q = q.Offset((*filter.PageNumber - 1) * *filter.PageSize)
	}
	q = q.Limit(*filter.PageSize)
	if len(filter.OrderBy) > 0 {
		q = q.OrderBy(filter.OrderBy...)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.SelectContext(ctx, &dto, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return dto.ToModels(), nil
}

func (r *SessionRepository) Count(
	ctx context.Context,
	filter *models.SessionFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.sessions")
	if filter.Search != nil {
		q = q.Where(
			postgresql.Search{
				Lang:   "english",
				Query:  *filter.Search,
				Fields: []string{"description"},
			},
		)
	}
	if len(filter.IDs) > 0 {
		q = q.Where(sq.Eq{"id": filter.IDs})
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
func (r *SessionRepository) Update(ctx context.Context, session *models.Session) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewSessionDTOFromModel(session)
	q := sq.Update("public.sessions").Where(sq.Eq{"id": session.ID})
	{
		q = q.Set("sessions.updated_at", dto.UpdatedAt)
		q = q.Set("sessions.title", dto.Title)
		q = q.Set("sessions.description", dto.Description)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("session_id", fmt.Sprint(session.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("session_id", fmt.Sprint(session.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().WithParam("session_id", fmt.Sprint(session.ID))
		return e
	}
	return nil
}
func (r *SessionRepository) Delete(ctx context.Context, id models.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.sessions").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("session_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("session_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().WithParam("session_id", fmt.Sprint(id))
		return e
	}
	return nil
}

type SessionDTO struct {
	ID          string    `db:"id,omitempty"`
	UpdatedAt   time.Time `db:"updated_at,omitempty"`
	CreatedAt   time.Time `db:"created_at,omitempty"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
}
type SessionListDTO []*SessionDTO

func (list SessionListDTO) ToModels() []*models.Session {
	listSessions := make([]*models.Session, len(list))
	for i := range list {
		listSessions[i] = list[i].ToModel()
	}
	return listSessions
}
func NewSessionDTOFromModel(session *models.Session) *SessionDTO {
	dto := &SessionDTO{
		ID:          string(session.ID),
		UpdatedAt:   session.UpdatedAt,
		CreatedAt:   session.CreatedAt,
		Title:       session.Title,
		Description: session.Description,
	}
	return dto
}
func (dto *SessionDTO) ToModel() *models.Session {
	model := &models.Session{
		ID:          models.UUID(dto.ID),
		UpdatedAt:   dto.UpdatedAt,
		CreatedAt:   dto.CreatedAt,
		Title:       dto.Title,
		Description: dto.Description,
	}
	return model
}
