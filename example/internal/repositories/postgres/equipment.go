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

type EquipmentRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewEquipmentRepository(database *sqlx.DB, logger log.Logger) repositories.EquipmentRepository {
	return &EquipmentRepository{database: database, logger: logger}
}
func (r *EquipmentRepository) Create(ctx context.Context, equipment *models.Equipment) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewEquipmentDTOFromModel(equipment)
	q := sq.Insert("public.equipment").
		Columns("updated_at", "created_at", "name", "repeat", "weight").
		Values(dto.UpdatedAt, dto.CreatedAt, dto.Name, dto.Repeat, dto.Weight).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(dto); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	equipment.ID = models.UUID(dto.ID)
	return nil
}
func (r *EquipmentRepository) Get(ctx context.Context, id models.UUID) (*models.Equipment, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &EquipmentDTO{}
	q := sq.Select("equipment.id", "equipment.updated_at", "equipment.created_at", "equipment.name", "equipment.repeat", "equipment.weight").
		From("public.equipment").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("equipment_id", string(id))
		return nil, e
	}
	return dto.ToModel(), nil
}

func (r *EquipmentRepository) List(
	ctx context.Context,
	filter *models.EquipmentFilter,
) ([]*models.Equipment, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto EquipmentListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = utils.Pointer(pageSize)
	}
	q := sq.Select("equipment.id", "equipment.updated_at", "equipment.created_at", "equipment.name", "equipment.repeat", "equipment.weight").
		From("public.equipment").
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

func (r *EquipmentRepository) Count(
	ctx context.Context,
	filter *models.EquipmentFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.equipment")
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
func (r *EquipmentRepository) Update(ctx context.Context, equipment *models.Equipment) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewEquipmentDTOFromModel(equipment)
	q := sq.Update("public.equipment").Where(sq.Eq{"id": equipment.ID})
	{
		q = q.Set("equipment.updated_at", dto.UpdatedAt)
		q = q.Set("equipment.name", dto.Name)
		q = q.Set("equipment.repeat", dto.Repeat)
		q = q.Set("equipment.weight", dto.Weight)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("equipment_id", fmt.Sprint(equipment.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("equipment_id", fmt.Sprint(equipment.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().WithParam("equipment_id", fmt.Sprint(equipment.ID))
		return e
	}
	return nil
}
func (r *EquipmentRepository) Delete(ctx context.Context, id models.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.equipment").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("equipment_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("equipment_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().WithParam("equipment_id", fmt.Sprint(id))
		return e
	}
	return nil
}

type EquipmentDTO struct {
	ID        string    `db:"id,omitempty"`
	UpdatedAt time.Time `db:"updated_at,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	Name      string    `db:"name"`
	Repeat    int       `db:"repeat"`
	Weight    int       `db:"weight"`
}
type EquipmentListDTO []*EquipmentDTO

func (list EquipmentListDTO) ToModels() []*models.Equipment {
	listEquipment := make([]*models.Equipment, len(list))
	for i := range list {
		listEquipment[i] = list[i].ToModel()
	}
	return listEquipment
}
func NewEquipmentDTOFromModel(equipment *models.Equipment) *EquipmentDTO {
	dto := &EquipmentDTO{
		ID:        string(equipment.ID),
		UpdatedAt: equipment.UpdatedAt,
		CreatedAt: equipment.CreatedAt,
		Name:      equipment.Name,
		Repeat:    equipment.Repeat,
		Weight:    equipment.Weight,
	}
	return dto
}
func (dto *EquipmentDTO) ToModel() *models.Equipment {
	model := &models.Equipment{
		ID:        models.UUID(dto.ID),
		UpdatedAt: dto.UpdatedAt,
		CreatedAt: dto.CreatedAt,
		Name:      dto.Name,
		Repeat:    dto.Repeat,
		Weight:    dto.Weight,
	}
	return model
}
