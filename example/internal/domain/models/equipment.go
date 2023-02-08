package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Equipment struct {
	ID        UUID      `json:"id" form:"id"`
	Name      string    `json:"name" form:"name"`
	Repeat    int       `json:"repeat" form:"repeat"`
	Weight    int       `json:"weight" form:"weight"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
	CreatedAt time.Time `json:"created_at" form:"created_at,omitempty"`
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
	IDs        []UUID   `json:"ids" form:"ids"`
	PageSize   *uint64  `json:"page_size" form:"page_size"`
	PageNumber *uint64  `json:"page_number" form:"page_number"`
	OrderBy    []string `json:"order_by" form:"order_by"`
	Search     *string  `json:"search" form:"search"`
}

func (c *EquipmentFilter) Validate() error {
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

type EquipmentCreate struct {
	Name   string `json:"name" form:"name"`
	Repeat int    `json:"repeat" form:"repeat"`
	Weight int    `json:"weight" form:"weight"`
}

func (c *EquipmentCreate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Repeat, validation.Required),
		validation.Field(&c.Weight, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type EquipmentUpdate struct {
	ID     UUID    `json:"id"`
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
