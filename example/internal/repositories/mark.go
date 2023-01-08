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

type MarkRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewMarkRepository(
	database *sqlx.DB,
	logger log.Logger,
) repositories.MarkRepository {
	return &MarkRepository{
		database: database,
		logger:   logger,
	}
}

func (r *MarkRepository) Create(
	ctx context.Context,
	mark *models.Mark,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.marks").
		Columns(
			"name",
			"title",
			"weight",
			"updated_at",
			"created_at",
		).
		Values(
			mark.Name,
			mark.Title,
			mark.Weight,
			mark.UpdatedAt,
			mark.CreatedAt,
		).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(mark); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}

func (r *MarkRepository) Get(
	ctx context.Context,
	id string,
) (*models.Mark, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	mark := &models.Mark{}
	q := sq.Select(
		"marks.id",
		"marks.name",
		"marks.title",
		"marks.weight",
		"marks.updated_at",
		"marks.created_at",
	).
		From("public.marks").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, mark, query, args...); err != nil {
		e := errs.FromPostgresError(err).
			WithParam("mark_id", id)
		return nil, e
	}
	return mark, nil
}

func (r *MarkRepository) List(
	ctx context.Context,
	filter *models.MarkFilter,
) ([]*models.Mark, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var marks []*models.Mark
	const pageSize = 10
	q := sq.Select(
		"marks.id",
		"marks.name",
		"marks.title",
		"marks.weight",
		"marks.updated_at",
		"marks.created_at",
	).
		From("public.marks").
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
	if err := r.database.SelectContext(ctx, &marks, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return marks, nil
}

func (r *MarkRepository) Update(
	ctx context.Context,
	mark *models.Mark,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.marks").
		Where(sq.Eq{"id": mark.ID}).
		Set("marks.name", mark.Name).
		Set("marks.title", mark.Title).
		Set("marks.weight", mark.Weight).
		Set("updated_at", mark.UpdatedAt)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("mark_id", fmt.Sprint(mark.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).
			WithParam("mark_id", fmt.Sprint(mark.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("mark_id", fmt.Sprint(mark.ID))
		return e
	}
	return nil
}

func (r *MarkRepository) Delete(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.marks").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("mark_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("mark_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("mark_id", fmt.Sprint(id))
		return e
	}
	return nil
}

func (r *MarkRepository) Count(
	ctx context.Context,
	filter *models.MarkFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.marks")
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
