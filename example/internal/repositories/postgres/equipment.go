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

type EquipmentDTO struct {
	ID        string    `db:"id,omitempty"`
	Name      string    `db:"name"`
	Repeat    int       `db:"repeat"`
	Weight    int       `db:"weight"`
	UpdatedAt time.Time `db:"updated_at,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
}

func NewEquipmentDTOFromModel(equipment *models.Equipment) *EquipmentDTO {
	dto := &EquipmentDTO{
		ID:        string(equipment.ID),
		Name:      string(equipment.Name),
		Repeat:    int(equipment.Repeat),
		Weight:    int(equipment.Weight),
		UpdatedAt: equipment.UpdatedAt,
		CreatedAt: equipment.CreatedAt,
	}
	return dto
}

func (dto *EquipmentDTO) ToModel() *models.Equipment {
	model := &models.Equipment{
		ID:        models.UUID(dto.ID),
		Name:      string(dto.Name),
		Repeat:    int(dto.Repeat),
		Weight:    int(dto.Weight),
		UpdatedAt: dto.UpdatedAt,
		CreatedAt: dto.CreatedAt,
	}
	return model
}

type EquipmentListDTO []*EquipmentDTO

func (list EquipmentListDTO) ToModels() []*models.Equipment {
	listEquipment := make([]*models.Equipment, len(list))
	for i := range list {
		listEquipment[i] = list[i].ToModel()
	}
	return listEquipment
}

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
	dto := NewEquipmentDTOFromModel(equipment)
	q := sq.Insert("public.equipment").
		Columns(
			"name",
			"repeat",
			"weight",
			"updated_at",
			"created_at",
		).
		Values(
			dto.Name,
			dto.Repeat,
			dto.Weight,
			dto.UpdatedAt,
			dto.CreatedAt,
		).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(dto); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	equipment.ID = models.UUID(dto.ID)
	return nil
}

func (r *EquipmentRepository) Get(
	ctx context.Context,
	id models.UUID,
) (*models.Equipment, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &EquipmentDTO{}
	q := sq.Select(
		"equipment.id",
		"equipment.name",
		"equipment.repeat",
		"equipment.weight",
		"equipment.updated_at",
		"equipment.created_at",
	).
		From("public.equipment").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).
			WithParam("equipment_id", string(id))
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
	q := sq.Select(
		"equipment.id",
		"equipment.name",
		"equipment.repeat",
		"equipment.weight",
		"equipment.updated_at",
		"equipment.created_at",
	).
		From("public.equipment").
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

func (r *EquipmentRepository) Update(
	ctx context.Context,
	equipment *models.Equipment,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewEquipmentDTOFromModel(equipment)
	q := sq.Update("public.equipment").
		Where(sq.Eq{"id": equipment.ID}).
		Set("equipment.name", dto.Name).
		Set("equipment.repeat", dto.Repeat).
		Set("equipment.weight", dto.Weight).
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
	id models.UUID,
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
