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

type PlanRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewPlanRepository(
	database *sqlx.DB,
	logger log.Logger,
) repositories.PlanRepository {
	return &PlanRepository{
		database: database,
		logger:   logger,
	}
}

func (r *PlanRepository) Create(
	ctx context.Context,
	plan *models.Plan,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.plans").
		Columns(
			"name",
			"repeat",
			"equipment_id",
			"updated_at",
			"created_at",
		).
		Values(
			plan.Name,
			plan.Repeat,
			plan.EquipmentID,
			plan.UpdatedAt,
			plan.CreatedAt,
		).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(plan); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}

func (r *PlanRepository) Get(
	ctx context.Context,
	id models.UUID,
) (*models.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	plan := &models.Plan{}
	q := sq.Select(
		"plans.id",
		"plans.name",
		"plans.repeat",
		"plans.equipment_id",
		"plans.updated_at",
		"plans.created_at",
	).
		From("public.plans").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, plan, query, args...); err != nil {
		e := errs.FromPostgresError(err).
			WithParam("plan_id", string(id))
		return nil, e
	}
	return plan, nil
}

func (r *PlanRepository) List(
	ctx context.Context,
	filter *models.PlanFilter,
) ([]*models.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var listPlans []*models.Plan
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = utils.Pointer(pageSize)
	}
	q := sq.Select(
		"plans.id",
		"plans.name",
		"plans.repeat",
		"plans.equipment_id",
		"plans.updated_at",
		"plans.created_at",
	).
		From("public.plans").
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
	if err := r.database.SelectContext(ctx, &listPlans, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return listPlans, nil
}

func (r *PlanRepository) Update(
	ctx context.Context,
	plan *models.Plan,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.plans").
		Where(sq.Eq{"id": plan.ID}).
		Set("plans.name", plan.Name).
		Set("plans.repeat", plan.Repeat).
		Set("plans.equipment_id", plan.EquipmentID).
		Set("updated_at", plan.UpdatedAt)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("plan_id", fmt.Sprint(plan.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).
			WithParam("plan_id", fmt.Sprint(plan.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("plan_id", fmt.Sprint(plan.ID))
		return e
	}
	return nil
}

func (r *PlanRepository) Delete(
	ctx context.Context,
	id models.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.plans").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("plan_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("plan_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("plan_id", fmt.Sprint(id))
		return e
	}
	return nil
}

func (r *PlanRepository) Count(
	ctx context.Context,
	filter *models.PlanFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.plans")
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
