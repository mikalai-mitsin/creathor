package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Plan struct {
	ID          UUID      `json:"id" db:"id,omitempty" form:"id"`
	Name        string    `json:"name" db:"name" form:"name"`
	Repeat      uint64    `json:"repeat" db:"repeat" form:"repeat"`
	EquipmentID string    `json:"equipment_id" db:"equipment_id" form:"equipment_id"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at,omitempty" form:"updated_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at,omitempty" form:"created_at,omitempty"`
}

func (c *Plan) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, is.UUID),
		validation.Field(&c.Name),
		validation.Field(&c.Repeat),
		validation.Field(&c.EquipmentID),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type PlanFilter struct {
	IDs        []UUID   `json:"ids" form:"ids"`
	PageSize   *uint64  `json:"page_size" form:"page_size"`
	PageNumber *uint64  `json:"page_number" form:"page_number"`
	OrderBy    []string `json:"order_by" form:"order_by"`
	Search     *string  `json:"search" form:"search"`
}

func (c *PlanFilter) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.IDs),
		validation.Field(&c.PageSize),
		validation.Field(&c.PageNumber),
		validation.Field(&c.OrderBy),
		validation.Field(&c.Search),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type PlanCreate struct {
	Name        string `json:"name" form:"name"`
	Repeat      uint64 `json:"repeat" form:"repeat"`
	EquipmentID string `json:"equipment_id" form:"equipment_id"`
}

func (c *PlanCreate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.Name),
		validation.Field(&c.Repeat),
		validation.Field(&c.EquipmentID),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type PlanUpdate struct {
	ID          UUID    `json:"id"`
	Name        *string `json:"name" form:"name"`
	Repeat      *uint64 `json:"repeat" form:"repeat"`
	EquipmentID *string `json:"equipment_id" form:"equipment_id"`
}

func (c *PlanUpdate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, validation.Required, is.UUID),
		validation.Field(&c.Name),
		validation.Field(&c.Repeat),
		validation.Field(&c.EquipmentID),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

const (
	PermissionIDPlanList   PermissionID = "plan_list"
	PermissionIDPlanDetail PermissionID = "plan_detail"
	PermissionIDPlanCreate PermissionID = "plan_create"
	PermissionIDPlanUpdate PermissionID = "plan_update"
	PermissionIDPlanDelete PermissionID = "plan_delete"
)
