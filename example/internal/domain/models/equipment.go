package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	PermissionIDEquipmentList   PermissionID = "equipment_list"
	PermissionIDEquipmentDetail PermissionID = "equipment_detail"
	PermissionIDEquipmentCreate PermissionID = "equipment_create"
	PermissionIDEquipmentUpdate PermissionID = "equipment_update"
	PermissionIDEquipmentDelete PermissionID = "equipment_delete"
)

type Equipment struct {
	ID          UUID      `json:"id"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Weight      uint64    `json:"weight"`
	Versions    []uint64  `json:"versions"`
	Release     time.Time `json:"release"`
	Tested      time.Time `json:"tested"`
}

func (m *Equipment) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Description, validation.Required),
		validation.Field(&m.Weight, validation.Required),
		validation.Field(&m.Versions, validation.Required),
		validation.Field(&m.Release, validation.Required),
		validation.Field(&m.Tested, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type EquipmentCreate struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Weight      uint64    `json:"weight"`
	Versions    []uint64  `json:"versions"`
	Release     time.Time `json:"release"`
	Tested      time.Time `json:"tested"`
}

func (m *EquipmentCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Description, validation.Required),
		validation.Field(&m.Weight, validation.Required),
		validation.Field(&m.Versions, validation.Required),
		validation.Field(&m.Release, validation.Required),
		validation.Field(&m.Tested, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type EquipmentUpdate struct {
	ID          UUID       `json:"id"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Weight      *uint64    `json:"weight"`
	Versions    *[]uint64  `json:"versions"`
	Release     *time.Time `json:"release"`
	Tested      *time.Time `json:"tested"`
}

func (m *EquipmentUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.Title),
		validation.Field(&m.Description),
		validation.Field(&m.Weight),
		validation.Field(&m.Versions),
		validation.Field(&m.Release),
		validation.Field(&m.Tested),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type EquipmentFilter struct {
	IDs        []UUID   `json:"ids"`
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	OrderBy    []string `json:"order_by"`
	Search     *string  `json:"search"`
}

func (m *EquipmentFilter) Validate() error {
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
