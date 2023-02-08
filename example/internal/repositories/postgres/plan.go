package postgres

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

type PlanDTO struct {
	ID          string    `db:"id,omitempty"`
	Name        string    `db:"name"`
	Repeat      int64     `db:"repeat"`
	EquipmentID string    `db:"equipment_id"`
	UpdatedAt   time.Time `db:"updated_at,omitempty"`
	CreatedAt   time.Time `db:"created_at,omitempty"`
}

func NewPlanDTOFromModel(plan *models.Plan) *PlanDTO {
	dto := &PlanDTO{
		ID:          string(plan.ID),
		Name:        string(plan.Name),
		Repeat:      int64(plan.Repeat),
		EquipmentID: string(plan.EquipmentID),
		UpdatedAt:   plan.UpdatedAt,
		CreatedAt:   plan.CreatedAt,
	}
	return dto
}

func (dto *PlanDTO) ToModel() *models.Plan {
	model := &models.Plan{
		ID:          models.UUID(dto.ID),
		Name:        string(dto.Name),
		Repeat:      uint64(dto.Repeat),
		EquipmentID: string(dto.EquipmentID),
		UpdatedAt:   dto.UpdatedAt,
		CreatedAt:   dto.CreatedAt,
	}
	return model
}

type PlanListDTO []*PlanDTO

func (list PlanListDTO) ToModels() []*models.Plan {
	listPlans := make([]*models.Plan, len(list))
	for i := range list {
		listPlans[i] = list[i].ToModel()
	}
	return listPlans
}

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
	dto := NewPlanDTOFromModel(plan)
	q := sq.Insert("public.plans").
		Columns(
			"name",
			"repeat",
			"equipment_id",
			"updated_at",
			"created_at",
		).
		Values(
			dto.Name,
			dto.Repeat,
			dto.EquipmentID,
			dto.UpdatedAt,
			dto.CreatedAt,
		).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(dto); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	plan.ID = models.UUID(dto.ID)
	return nil
}

func (r *PlanRepository) Get(
	ctx context.Context,
	id models.UUID,
) (*models.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &PlanDTO{}
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
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).
			WithParam("plan_id", string(id))
		return nil, e
	}
	return dto.ToModel(), nil
}

func (r *PlanRepository) List(
	ctx context.Context,
	filter *models.PlanFilter,
) ([]*models.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto PlanListDTO
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

func (r *PlanRepository) Update(
	ctx context.Context,
	plan *models.Plan,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewPlanDTOFromModel(plan)
	q := sq.Update("public.plans").
		Where(sq.Eq{"id": plan.ID}).
		Set("plans.name", dto.Name).
		Set("plans.repeat", dto.Repeat).
		Set("plans.equipment_id", dto.EquipmentID).
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
