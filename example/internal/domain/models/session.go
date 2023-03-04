package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	PermissionIDSessionList   PermissionID = "session_list"
	PermissionIDSessionDetail PermissionID = "session_detail"
	PermissionIDSessionCreate PermissionID = "session_create"
	PermissionIDSessionUpdate PermissionID = "session_update"
	PermissionIDSessionDelete PermissionID = "session_delete"
)

type Session struct {
	ID          UUID      `json:"id"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

func (m *Session) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Description, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type SessionCreate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (m *SessionCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Description, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type SessionUpdate struct {
	ID          UUID    `json:"id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (m *SessionUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.Title),
		validation.Field(&m.Description),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type SessionFilter struct {
	IDs        []UUID   `json:"ids"`
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	OrderBy    []string `json:"order_by"`
	Search     *string  `json:"search"`
}

func (m *SessionFilter) Validate() error {
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
