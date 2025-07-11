package entity

import (
	"errors"
	"net/mail"
	"strings"
)

const (
	EMAIL_SUFIX = "@st.kobedenshi.ac.jp"
)

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrEmailInvalidDomain = errors.New("invalid email domain")
)

type Email string

func NewEmail(email string) (Email, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return "", ErrInvalidEmail
	}
	if !strings.HasSuffix(email, EMAIL_SUFIX) {
		return "", ErrEmailInvalidDomain
	}
	return Email(email), nil
}

func (e Email) String() string {
	return string(e)
}

func (e Email) IsTeacher() bool {
	return true
}
