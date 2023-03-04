package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	PermissionIDArchList   PermissionID = "arch_list"
	PermissionIDArchDetail PermissionID = "arch_detail"
	PermissionIDArchCreate PermissionID = "arch_create"
	PermissionIDArchUpdate PermissionID = "arch_update"
	PermissionIDArchDelete PermissionID = "arch_delete"
)

type Arch struct {
	ID          UUID      `json:"id"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Versions    []uint64  `json:"versions"`
	Release     time.Time `json:"release"`
	Tested      time.Time `json:"tested"`
}

func (m *Arch) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Description, validation.Required),
		validation.Field(&m.Tags, validation.Required),
		validation.Field(&m.Versions, validation.Required),
		validation.Field(&m.Release, validation.Required),
		validation.Field(&m.Tested, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type ArchCreate struct {
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Versions    []uint64  `json:"versions"`
	Release     time.Time `json:"release"`
	Tested      time.Time `json:"tested"`
}

func (m *ArchCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Description, validation.Required),
		validation.Field(&m.Tags, validation.Required),
		validation.Field(&m.Versions, validation.Required),
		validation.Field(&m.Release, validation.Required),
		validation.Field(&m.Tested, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type ArchUpdate struct {
	ID          UUID       `json:"id"`
	Name        *string    `json:"name"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Tags        *[]string  `json:"tags"`
	Versions    *[]uint64  `json:"versions"`
	Release     *time.Time `json:"release"`
	Tested      *time.Time `json:"tested"`
}

func (m *ArchUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.Name),
		validation.Field(&m.Title),
		validation.Field(&m.Description),
		validation.Field(&m.Tags),
		validation.Field(&m.Versions),
		validation.Field(&m.Release),
		validation.Field(&m.Tested),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type ArchFilter struct {
	IDs        []UUID   `json:"ids"`
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	OrderBy    []string `json:"order_by"`
	Search     *string  `json:"search"`
}

func (m *ArchFilter) Validate() error {
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
