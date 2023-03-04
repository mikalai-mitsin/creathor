package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	PermissionIDPlanList   PermissionID = "plan_list"
	PermissionIDPlanDetail PermissionID = "plan_detail"
	PermissionIDPlanCreate PermissionID = "plan_create"
	PermissionIDPlanUpdate PermissionID = "plan_update"
	PermissionIDPlanDelete PermissionID = "plan_delete"
)

type Plan struct {
	ID          UUID      `json:"id"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Repeat      uint64    `json:"repeat"`
	EquipmentID string    `json:"equipment_id"`
}

func (m *Plan) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Repeat, validation.Required),
		validation.Field(&m.EquipmentID, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type PlanCreate struct {
	Name        string `json:"name"`
	Repeat      uint64 `json:"repeat"`
	EquipmentID string `json:"equipment_id"`
}

func (m *PlanCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Repeat, validation.Required),
		validation.Field(&m.EquipmentID, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type PlanUpdate struct {
	ID          UUID    `json:"id"`
	Name        *string `json:"name"`
	Repeat      *uint64 `json:"repeat"`
	EquipmentID *string `json:"equipment_id"`
}

func (m *PlanUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.Name),
		validation.Field(&m.Repeat),
		validation.Field(&m.EquipmentID),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type PlanFilter struct {
	IDs        []UUID   `json:"ids"`
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	OrderBy    []string `json:"order_by"`
	Search     *string  `json:"search"`
}

func (m *PlanFilter) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.IDs),
		validation.Field(&m.PageNumber),
		validation.Field(&m.PageSize),
		validation.Field(&m.OrderBy),
		validation.Field(&m.Search),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}
