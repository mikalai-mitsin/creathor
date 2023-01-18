package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Equipment struct {
	ID        string    `json:"id" db:"id,omitempty" form:"id"`
	Name      string    `json:"name" db:"name" form:"name"`
	Repeat    int       `json:"repeat" db:"repeat" form:"repeat"`
	Weight    int       `json:"weight" db:"weight" form:"weight"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at,omitempty" form:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at,omitempty" form:"created_at,omitempty"`
}

func (c *Equipment) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, is.UUID),
		validation.Field(&c.Name),
		validation.Field(&c.Repeat),
		validation.Field(&c.Weight),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type EquipmentFilter struct {
	IDs        []string `json:"ids" form:"ids"`
	PageSize   *uint64  `json:"page_size" form:"page_size"`
	PageNumber *uint64  `json:"page_number" form:"page_number"`
	OrderBy    []string `json:"order_by" form:"order_by"`
}

func (c *EquipmentFilter) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.IDs),
		validation.Field(&c.PageSize),
		validation.Field(&c.PageNumber),
		validation.Field(&c.OrderBy),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type EquipmentCreate struct {
	Name   string `json:"name" form:"name"`
	Repeat int    `json:"repeat" form:"repeat"`
	Weight int    `json:"weight" form:"weight"`
}

func (c *EquipmentCreate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.Name),
		validation.Field(&c.Repeat),
		validation.Field(&c.Weight),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type EquipmentUpdate struct {
	ID     string  `json:"id"`
	Name   *string `json:"name" form:"name"`
	Repeat *int    `json:"repeat" form:"repeat"`
	Weight *int    `json:"weight" form:"weight"`
}

func (c *EquipmentUpdate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, validation.Required, is.UUID),
		validation.Field(&c.Name),
		validation.Field(&c.Repeat),
		validation.Field(&c.Weight),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

const (
	PermissionIDEquipmentList   PermissionID = "equipment_list"
	PermissionIDEquipmentDetail PermissionID = "equipment_detail"
	PermissionIDEquipmentCreate PermissionID = "equipment_create"
	PermissionIDEquipmentUpdate PermissionID = "equipment_update"
	PermissionIDEquipmentDelete PermissionID = "equipment_delete"
)
