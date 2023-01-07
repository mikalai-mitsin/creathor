package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserSession struct {
	ID        string    `json:"id" db:"id,omitempty"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at,omitempty"`
}

func (c *UserSession) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, is.UUID),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type UserSessionFilter struct {
	IDs        []string `json:"ids"`
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	OrderBy    []string `json:"order_by"`
}

func (c *UserSessionFilter) Validate() error {
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

type UserSessionCreate struct {
}

func (c *UserSessionCreate) Validate() error {
	err := validation.ValidateStruct(
		c,
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type UserSessionUpdate struct {
	ID string `json:"id"`
}

func (c *UserSessionUpdate) Validate() error {
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
	PermissionIDUserSessionList   PermissionID = "user_session_list"
	PermissionIDUserSessionDetail PermissionID = "user_session_detail"
	PermissionIDUserSessionCreate PermissionID = "user_session_create"
	PermissionIDUserSessionUpdate PermissionID = "user_session_update"
	PermissionIDUserSessionDelete PermissionID = "user_session_delete"
)
