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

type PostgresEquipmentRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewPostgresEquipmentRepository(database *sqlx.DB, logger log.Logger) repositories.EquipmentRepository {
	return &PostgresEquipmentRepository{database: database, logger: logger}
}

func (r *PostgresEquipmentRepository) Create(ctx context.Context, equipment *models.Equipment) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.equipments").
		Columns(). // TODO: add columns
		Values().  // TODO: add values
		Suffix("RETURNING \"id\"")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).Scan(&equipment.ID); err != nil {
		e := errs.NewUnexpectedBehaviorError(err.Error())
		return e
	}
	return nil
}

func (r *PostgresEquipmentRepository) Get(ctx context.Context, id string) (*models.Equipment, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	equipment := &models.Equipment{}
	q := sq.Select("*").
		From("public.equipments").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, &equipment, query, args...); err != nil {
		e := errs.NewUnexpectedBehaviorError(err.Error())
		return nil, e
	}
	return equipment, nil
}

func (r *PostgresEquipmentRepository) List(ctx context.Context, filter *models.EquipmentFilter) ([]*models.Equipment, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var equipments []*models.Equipment
	const pageSize = 10
	q := sq.Select("*").
		From("public.equipments").
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
	if err := r.database.SelectContext(ctx, &equipments, query, args...); err != nil {
		e := errs.NewUnexpectedBehaviorError(err.Error())
		return nil, e
	}
	return equipments, nil
}

func (r *PostgresEquipmentRepository) Update(ctx context.Context, equipment *models.Equipment) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.equipments").Where(sq.Eq{"id": equipment.ID}).Set("", "") // TODO: set values
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
		e.AddParam("equipment_id", fmt.Sprint(equipment.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.NewUnexpectedBehaviorError(err.Error())
	}
	if affected == 0 {
		e := errs.NewEquipmentNotFound()
		e.AddParam("equipment_id", fmt.Sprint(equipment.ID))
		return e
	}
	return nil
}

func (r *PostgresEquipmentRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.equipments").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.NewUnexpectedBehaviorError(err.Error())
		e.AddParam("equipment_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.NewUnexpectedBehaviorError(err.Error())
		e.AddParam("equipment_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEquipmentNotFound()
		e.AddParam("equipment_id", fmt.Sprint(id))
		return e
	}
	return nil
}
