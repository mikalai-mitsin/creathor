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
	"github.com/018bf/example/pkg/postgresql"
	"github.com/018bf/example/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type ArchRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewArchRepository(
	database *sqlx.DB,
	logger log.Logger,
) repositories.ArchRepository {
	return &ArchRepository{
		database: database,
		logger:   logger,
	}
}

func (r *ArchRepository) Create(
	ctx context.Context,
	arch *models.Arch,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.arches").
		Columns(
			"name",
			"release",
			"tested",
			"updated_at",
			"created_at",
		).
		Values(
			arch.Name,
			arch.Release,
			arch.Tested,
			arch.UpdatedAt,
			arch.CreatedAt,
		).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(arch); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}

func (r *ArchRepository) Get(
	ctx context.Context,
	id models.UUID,
) (*models.Arch, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	arch := &models.Arch{}
	q := sq.Select(
		"arches.id",
		"arches.name",
		"arches.release",
		"arches.tested",
		"arches.updated_at",
		"arches.created_at",
	).
		From("public.arches").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, arch, query, args...); err != nil {
		e := errs.FromPostgresError(err).
			WithParam("arch_id", string(id))
		return nil, e
	}
	return arch, nil
}

func (r *ArchRepository) List(
	ctx context.Context,
	filter *models.ArchFilter,
) ([]*models.Arch, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var listArches []*models.Arch
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = utils.Pointer(pageSize)
	}
	q := sq.Select(
		"arches.id",
		"arches.name",
		"arches.release",
		"arches.tested",
		"arches.updated_at",
		"arches.created_at",
	).
		From("public.arches").
		Limit(pageSize)
	if filter.Search != nil {
		q = q.Where(postgresql.Search{
			Lang:  "english",
			Query: *filter.Search,
			Fields: []string{
				"name",
			},
		})
	}
	// TODO: add filtering
	if filter.PageNumber != nil && *filter.PageNumber > 1 {
		q = q.Offset((*filter.PageNumber - 1) * *filter.PageSize)
	}
	q = q.Limit(*filter.PageSize)
	if len(filter.OrderBy) > 0 {
		q = q.OrderBy(filter.OrderBy...)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.SelectContext(ctx, &listArches, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return listArches, nil
}

func (r *ArchRepository) Update(
	ctx context.Context,
	arch *models.Arch,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.arches").
		Where(sq.Eq{"id": arch.ID}).
		Set("arches.name", arch.Name).
		Set("arches.release", arch.Release).
		Set("arches.tested", arch.Tested).
		Set("updated_at", arch.UpdatedAt)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("arch_id", fmt.Sprint(arch.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).
			WithParam("arch_id", fmt.Sprint(arch.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("arch_id", fmt.Sprint(arch.ID))
		return e
	}
	return nil
}

func (r *ArchRepository) Delete(
	ctx context.Context,
	id models.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.arches").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("arch_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("arch_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("arch_id", fmt.Sprint(id))
		return e
	}
	return nil
}

func (r *ArchRepository) Count(
	ctx context.Context,
	filter *models.ArchFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.arches")
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
