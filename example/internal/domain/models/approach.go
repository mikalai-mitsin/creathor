package models

import (
	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

type Approach struct {
	ID        string    `json:"id" db:"id,omitempty"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at,omitempty"`
}

func (c *Approach) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, is.UUID),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type ApproachFilter struct {
	IDs        []string `json:"ids"`
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	OrderBy    []string `json:"order_by"`
}

func (c *ApproachFilter) Validate() error {
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

type ApproachCreate struct {
}

func (c *ApproachCreate) Validate() error {
	err := validation.ValidateStruct(
		c,
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type ApproachUpdate struct {
	ID string `json:"id"`
}

func (c *ApproachUpdate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, validation.Required, is.UUID),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

const (
	PermissionIDApproachList   PermissionID = "approach_list"
	PermissionIDApproachDetail PermissionID = "approach_detail"
	PermissionIDApproachCreate PermissionID = "approach_create"
	PermissionIDApproachUpdate PermissionID = "approach_update"
	PermissionIDApproachDelete PermissionID = "approach_delete"
)
