package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Session struct {
	ID          string    `json:"id" db:"id,omitempty" form:"id"`
	Title       string    `json:"title" db:"title" form:"title"`
	Description string    `json:"description" db:"description" form:"description"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at,omitempty" form:"updated_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at,omitempty" form:"created_at,omitempty"`
}

func (c *Session) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, is.UUID),
		validation.Field(&c.Title),
		validation.Field(&c.Description),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type SessionFilter struct {
	IDs        []string `json:"ids" form:"ids"`
	PageSize   *uint64  `json:"page_size" form:"page_size"`
	PageNumber *uint64  `json:"page_number" form:"page_number"`
	OrderBy    []string `json:"order_by" form:"order_by"`
}

func (c *SessionFilter) Validate() error {
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

type SessionCreate struct {
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
}

func (c *SessionCreate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.Title),
		validation.Field(&c.Description),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type SessionUpdate struct {
	ID          string  `json:"id"`
	Title       *string `json:"title" form:"title"`
	Description *string `json:"description" form:"description"`
}

func (c *SessionUpdate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, validation.Required, is.UUID),
		validation.Field(&c.Title),
		validation.Field(&c.Description),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

const (
	PermissionIDSessionList   PermissionID = "session_list"
	PermissionIDSessionDetail PermissionID = "session_detail"
	PermissionIDSessionCreate PermissionID = "session_create"
	PermissionIDSessionUpdate PermissionID = "session_update"
	PermissionIDSessionDelete PermissionID = "session_delete"
)
