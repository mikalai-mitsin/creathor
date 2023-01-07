package models

import (
	"fmt"
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/crypto/bcrypt"
)

const (
	PermissionIDUserList   PermissionID = "user_list"
	PermissionIDUserDetail PermissionID = "user_detail"
	PermissionIDUserCreate PermissionID = "user_create"
	PermissionIDUserUpdate PermissionID = "user_update"
	PermissionIDUserDelete PermissionID = "user_delete"
)

type User struct {
	ID        string    `db:"id,omitempty" json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Password  string    `db:"password" json:"password"`
	Email     string    `db:"email" json:"email"`
	GroupID   GroupID   `db:"group_id" json:"group_id"`
	CreatedAt time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (u *User) Validate() error {
	err := validation.ValidateStruct(
		u,
		validation.Field(&u.ID, is.UUID),
		validation.Field(&u.FirstName),
		validation.Field(&u.LastName),
		validation.Field(&u.Password),
		validation.Field(&u.Email, is.EmailFormat),
		validation.Field(&u.CreatedAt),
		validation.Field(&u.UpdatedAt),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

func (u *User) SetPassword(password string) {
	fromPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u.Password = string(fromPassword)
}

func (u *User) CheckPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return errs.NewInvalidParameter("email or password")
	}
	return nil
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

type UserFilter struct {
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	Search     *string  `json:"search"`
	OrderBy    []string `json:"order_by"`
}

func (c *UserFilter) Validate() error {
	err := validation.ValidateStruct(
		c,
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

type UserCreate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserCreate) Validate() error {
	err := validation.ValidateStruct(
		u,
		validation.Field(&u.Email, is.Email, validation.Required),
		validation.Field(&u.Password, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type UserUpdate struct {
	ID        string  `json:"id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Password  *string `json:"password"`
	Email     *string `json:"email"`
}

func (u *UserUpdate) Validate() error {
	err := validation.ValidateStruct(
		u,
		validation.Field(&u.ID, validation.Required, is.UUID),
		validation.Field(&u.FirstName),
		validation.Field(&u.LastName),
		validation.Field(&u.Password, validation.Length(8, 100)),
		validation.Field(&u.Email, is.EmailFormat),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type SetPassword struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

func (u *SetPassword) Validate() error {
	err := validation.ValidateStruct(
		u,
		validation.Field(&u.UserID, validation.Required, is.UUID),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 100)),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}
