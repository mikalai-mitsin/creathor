package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Day struct {
	ID          UUID      `json:"id" form:"id"`
	Name        string    `json:"name" form:"name"`
	Repeat      int       `json:"repeat" form:"repeat"`
	EquipmentID string    `json:"equipment_id" form:"equipment_id"`
	UpdatedAt   time.Time `json:"updated_at" form:"updated_at"`
	CreatedAt   time.Time `json:"created_at" form:"created_at,omitempty"`
}

func (c *Day) Validate() error {
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

type DayFilter struct {
	IDs        []UUID   `json:"ids" form:"ids"`
	PageSize   *uint64  `json:"page_size" form:"page_size"`
	PageNumber *uint64  `json:"page_number" form:"page_number"`
	OrderBy    []string `json:"order_by" form:"order_by"`
	Search     *string  `json:"search" form:"search"`
}

func (c *DayFilter) Validate() error {
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

type DayCreate struct {
	Name        string `json:"name" form:"name"`
	Repeat      int    `json:"repeat" form:"repeat"`
	EquipmentID string `json:"equipment_id" form:"equipment_id"`
}

func (c *DayCreate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Repeat, validation.Required),
		validation.Field(&c.EquipmentID, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type DayUpdate struct {
	ID          UUID    `json:"id"`
	Name        *string `json:"name" form:"name"`
	Repeat      *int    `json:"repeat" form:"repeat"`
	EquipmentID *string `json:"equipment_id" form:"equipment_id"`
}

func (c *DayUpdate) Validate() error {
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
	PermissionIDDayList   PermissionID = "day_list"
	PermissionIDDayDetail PermissionID = "day_detail"
	PermissionIDDayCreate PermissionID = "day_create"
	PermissionIDDayUpdate PermissionID = "day_update"
	PermissionIDDayDelete PermissionID = "day_delete"
)
