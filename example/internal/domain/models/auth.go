package models

import (
	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

type Token string

func (t Token) String() string {
	return string(t)
}

type TokenPair struct {
	Access  Token
	Refresh Token
}

func (c *TokenPair) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.Access),
		validation.Field(&c.Refresh),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type Login struct {
	Email    string
	Password string
}

func (c *Login) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.Email, is.Email),
		validation.Field(&c.Password),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

var Guest = &User{
	ID:        "",
	FirstName: "",
	LastName:  "",
	Password:  "",
	Email:     "",
	CreatedAt: time.Time{},
	UpdatedAt: time.Time{},
	GroupID:   GroupIDGuest,
}
