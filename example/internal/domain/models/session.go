package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Session struct {
	ID        string    `json:"id" db:"id,omitempty" form:"id"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at,omitempty" form:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at,omitempty" form:"created_at,omitempty"`
}

func (c *Session) Validate() *errs.Error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, is.UUID),
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

func (c *SessionFilter) Validate() *errs.Error {
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
}

func (c *SessionCreate) Validate() *errs.Error {
	err := validation.ValidateStruct(
		c,
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type SessionUpdate struct {
	ID string `json:"id"`
}

func (c *SessionUpdate) Validate() *errs.Error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, validation.Required, is.UUID),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}
