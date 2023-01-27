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

type DayRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewDayRepository(
	database *sqlx.DB,
	logger log.Logger,
) repositories.DayRepository {
	return &DayRepository{
		database: database,
		logger:   logger,
	}
}

func (r *DayRepository) Create(
	ctx context.Context,
	day *models.Day,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.days").
		Columns(
			"name",
			"repeat",
			"equipment_id",
			"updated_at",
			"created_at",
		).
		Values(
			day.Name,
			day.Repeat,
			day.EquipmentID,
			day.UpdatedAt,
			day.CreatedAt,
		).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(day); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}

func (r *DayRepository) Get(
	ctx context.Context,
	id models.UUID,
) (*models.Day, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	day := &models.Day{}
	q := sq.Select(
		"days.id",
		"days.name",
		"days.repeat",
		"days.equipment_id",
		"days.updated_at",
		"days.created_at",
	).
		From("public.days").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, day, query, args...); err != nil {
		e := errs.FromPostgresError(err).
			WithParam("day_id", string(id))
		return nil, e
	}
	return day, nil
}

func (r *DayRepository) List(
	ctx context.Context,
	filter *models.DayFilter,
) ([]*models.Day, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var listDays []*models.Day
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = utils.Pointer(pageSize)
	}
	q := sq.Select(
		"days.id",
		"days.name",
		"days.repeat",
		"days.equipment_id",
		"days.updated_at",
		"days.created_at",
	).
		From("public.days").
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
	if err := r.database.SelectContext(ctx, &listDays, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return listDays, nil
}

func (r *DayRepository) Update(
	ctx context.Context,
	day *models.Day,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.days").
		Where(sq.Eq{"id": day.ID}).
		Set("days.name", day.Name).
		Set("days.repeat", day.Repeat).
		Set("days.equipment_id", day.EquipmentID).
		Set("updated_at", day.UpdatedAt)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("day_id", fmt.Sprint(day.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).
			WithParam("day_id", fmt.Sprint(day.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("day_id", fmt.Sprint(day.ID))
		return e
	}
	return nil
}

func (r *DayRepository) Delete(
	ctx context.Context,
	id models.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.days").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("day_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("day_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("day_id", fmt.Sprint(id))
		return e
	}
	return nil
}

func (r *DayRepository) Count(
	ctx context.Context,
	filter *models.DayFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.days")
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
