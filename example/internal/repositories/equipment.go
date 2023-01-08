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

type EquipmentRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewEquipmentRepository(
	database *sqlx.DB,
	logger log.Logger,
) repositories.EquipmentRepository {
	return &EquipmentRepository{
		database: database,
		logger:   logger,
	}
}

func (r *EquipmentRepository) Create(
	ctx context.Context,
	equipment *models.Equipment,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.equipment").
		Columns(
			"updated_at",
			"created_at",
		).
		Values(
			equipment.UpdatedAt,
			equipment.CreatedAt,
		).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(equipment); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}

func (r *EquipmentRepository) Get(
	ctx context.Context,
	id string,
) (*models.Equipment, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	equipment := &models.Equipment{}
	q := sq.Select(
		"equipment.id",
		"equipment.updated_at",
		"equipment.created_at",
	).
		From("public.equipment").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, equipment, query, args...); err != nil {
		e := errs.FromPostgresError(err).
			WithParam("equipment_id", id)
		return nil, e
	}
	return equipment, nil
}

func (r *EquipmentRepository) List(
	ctx context.Context,
	filter *models.EquipmentFilter,
) ([]*models.Equipment, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var equipment []*models.Equipment
	const pageSize = 10
	q := sq.Select(
		"equipment.id",
		"equipment.updated_at",
		"equipment.created_at",
	).
		From("public.equipment").
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
	if err := r.database.SelectContext(ctx, &equipment, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return equipment, nil
}

func (r *EquipmentRepository) Update(
	ctx context.Context,
	equipment *models.Equipment,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.equipment").
		Where(sq.Eq{"id": equipment.ID}).
		Set("updated_at", equipment.UpdatedAt)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("equipment_id", fmt.Sprint(equipment.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).
			WithParam("equipment_id", fmt.Sprint(equipment.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("equipment_id", fmt.Sprint(equipment.ID))
		return e
	}
	return nil
}

func (r *EquipmentRepository) Delete(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.equipment").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("equipment_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("equipment_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("equipment_id", fmt.Sprint(id))
		return e
	}
	return nil
}

func (r *EquipmentRepository) Count(
	ctx context.Context,
	filter *models.EquipmentFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.equipment")
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
