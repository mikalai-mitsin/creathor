package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	PermissionIDDayList   PermissionID = "day_list"
	PermissionIDDayDetail PermissionID = "day_detail"
	PermissionIDDayCreate PermissionID = "day_create"
	PermissionIDDayUpdate PermissionID = "day_update"
	PermissionIDDayDelete PermissionID = "day_delete"
)

type Day struct {
	ID          UUID      `json:"id"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Repeat      int       `json:"repeat"`
	EquipmentID string    `json:"equipment_id"`
}

func (m *Day) Validate() error {
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

type DayCreate struct {
	Name        string `json:"name"`
	Repeat      int    `json:"repeat"`
	EquipmentID string `json:"equipment_id"`
}

func (m *DayCreate) Validate() error {
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

type DayUpdate struct {
	ID          UUID    `json:"id"`
	Name        *string `json:"name"`
	Repeat      *int    `json:"repeat"`
	EquipmentID *string `json:"equipment_id"`
}

func (m *DayUpdate) Validate() error {
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

type DayFilter struct {
	IDs        []UUID   `json:"ids"`
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	OrderBy    []string `json:"order_by"`
	Search     *string  `json:"search"`
}

func (m *DayFilter) Validate() error {
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
