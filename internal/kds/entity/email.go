package entity

import (
	"errors"
	"net/mail"
)

type Email string

var (
	ErrInvalidEmail = errors.New("invalid email")
)

func NewEmail(email string) (Email, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return "", ErrInvalidEmail
	}
	return Email(email), nil
}

func (e Email) String() string {
	return string(e)
}
