package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Arch struct {
	ID          UUID      `json:"id" form:"id"`
	Name        string    `json:"name" form:"name"`
	Tags        []string  `json:"tags" form:"tags"`
	Versions    []uint    `json:"versions" form:"versions"`
	OldVersions []uint64  `json:"old_versions" form:"old_versions"`
	Release     time.Time `json:"release" form:"release"`
	Tested      time.Time `json:"tested" form:"tested"`
	UpdatedAt   time.Time `json:"updated_at" form:"updated_at"`
	CreatedAt   time.Time `json:"created_at" form:"created_at,omitempty"`
}

func (c *Arch) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, is.UUID),
		validation.Field(&c.Name),
		validation.Field(&c.Tags),
		validation.Field(&c.Versions),
		validation.Field(&c.OldVersions),
		validation.Field(&c.Release),
		validation.Field(&c.Tested),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type ArchFilter struct {
	IDs        []UUID   `json:"ids" form:"ids"`
	PageSize   *uint64  `json:"page_size" form:"page_size"`
	PageNumber *uint64  `json:"page_number" form:"page_number"`
	OrderBy    []string `json:"order_by" form:"order_by"`
	Search     *string  `json:"search" form:"search"`
}

func (c *ArchFilter) Validate() error {
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

type ArchCreate struct {
	Name        string    `json:"name" form:"name"`
	Tags        []string  `json:"tags" form:"tags"`
	Versions    []uint    `json:"versions" form:"versions"`
	OldVersions []uint64  `json:"old_versions" form:"old_versions"`
	Release     time.Time `json:"release" form:"release"`
	Tested      time.Time `json:"tested" form:"tested"`
}

func (c *ArchCreate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Tags, validation.Required),
		validation.Field(&c.Versions, validation.Required),
		validation.Field(&c.OldVersions, validation.Required),
		validation.Field(&c.Release, validation.Required),
		validation.Field(&c.Tested, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type ArchUpdate struct {
	ID          UUID       `json:"id"`
	Name        *string    `json:"name" form:"name"`
	Tags        *[]string  `json:"tags" form:"tags"`
	Versions    *[]uint    `json:"versions" form:"versions"`
	OldVersions *[]uint64  `json:"old_versions" form:"old_versions"`
	Release     *time.Time `json:"release" form:"release"`
	Tested      *time.Time `json:"tested" form:"tested"`
}

func (c *ArchUpdate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, validation.Required, is.UUID),
		validation.Field(&c.Name),
		validation.Field(&c.Tags),
		validation.Field(&c.Versions),
		validation.Field(&c.OldVersions),
		validation.Field(&c.Release),
		validation.Field(&c.Tested),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

const (
	PermissionIDArchList   PermissionID = "arch_list"
	PermissionIDArchDetail PermissionID = "arch_detail"
	PermissionIDArchCreate PermissionID = "arch_create"
	PermissionIDArchUpdate PermissionID = "arch_update"
	PermissionIDArchDelete PermissionID = "arch_delete"
)
