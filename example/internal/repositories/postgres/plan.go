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

type PlanRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewPlanRepository(database *sqlx.DB, logger log.Logger) repositories.PlanRepository {
	return &PlanRepository{database: database, logger: logger}
}
func (r *PlanRepository) Create(ctx context.Context, plan *models.Plan) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewPlanDTOFromModel(plan)
	q := sq.Insert("public.plans").
		Columns("updated_at", "created_at", "name", "repeat", "equipment_id").
		Values(dto.UpdatedAt, dto.CreatedAt, dto.Name, dto.Repeat, dto.EquipmentID).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(dto); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	plan.ID = models.UUID(dto.ID)
	return nil
}
func (r *PlanRepository) Get(ctx context.Context, id models.UUID) (*models.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &PlanDTO{}
	q := sq.Select("plans.id", "plans.updated_at", "plans.created_at", "plans.name", "plans.repeat", "plans.equipment_id").
		From("public.plans").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("plan_id", string(id))
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
	q := sq.Select("plans.id", "plans.updated_at", "plans.created_at", "plans.name", "plans.repeat", "plans.equipment_id").
		From("public.plans").
		Limit(pageSize)
	if filter.Search != nil {
		q = q.Where(
			postgresql.Search{Lang: "english", Query: *filter.Search, Fields: []string{"name"}},
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
func (r *PlanRepository) Count(ctx context.Context, filter *models.PlanFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.plans")
	if filter.Search != nil {
		q = q.Where(
			postgresql.Search{Lang: "english", Query: *filter.Search, Fields: []string{"name"}},
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
func (r *PlanRepository) Update(ctx context.Context, plan *models.Plan) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewPlanDTOFromModel(plan)
	q := sq.Update("public.plans").Where(sq.Eq{"id": plan.ID})
	{
		q = q.Set("plans.updated_at", dto.UpdatedAt)
		q = q.Set("plans.name", dto.Name)
		q = q.Set("plans.repeat", dto.Repeat)
		q = q.Set("plans.equipment_id", dto.EquipmentID)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("plan_id", fmt.Sprint(plan.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("plan_id", fmt.Sprint(plan.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().WithParam("plan_id", fmt.Sprint(plan.ID))
		return e
	}
	return nil
}
func (r *PlanRepository) Delete(ctx context.Context, id models.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.plans").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("plan_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("plan_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().WithParam("plan_id", fmt.Sprint(id))
		return e
	}
	return nil
}

type PlanDTO struct {
	ID          string    `db:"id,omitempty"`
	UpdatedAt   time.Time `db:"updated_at,omitempty"`
	CreatedAt   time.Time `db:"created_at,omitempty"`
	Name        string    `db:"name"`
	Repeat      int64     `db:"repeat"`
	EquipmentID string    `db:"equipment_id"`
}
type PlanListDTO []*PlanDTO

func (list PlanListDTO) ToModels() []*models.Plan {
	listPlans := make([]*models.Plan, len(list))
	for i := range list {
		listPlans[i] = list[i].ToModel()
	}
	return listPlans
}
func NewPlanDTOFromModel(plan *models.Plan) *PlanDTO {
	dto := &PlanDTO{
		ID:          string(plan.ID),
		UpdatedAt:   plan.UpdatedAt,
		CreatedAt:   plan.CreatedAt,
		Name:        plan.Name,
		Repeat:      int64(plan.Repeat),
		EquipmentID: plan.EquipmentID,
	}
	return dto
}
func (dto *PlanDTO) ToModel() *models.Plan {
	model := &models.Plan{
		ID:          models.UUID(dto.ID),
		UpdatedAt:   dto.UpdatedAt,
		CreatedAt:   dto.CreatedAt,
		Name:        dto.Name,
		Repeat:      uint64(dto.Repeat),
		EquipmentID: dto.EquipmentID,
	}
	return model
}
