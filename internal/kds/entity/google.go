package entity

import "errors"

var (
	ErrGoogleUserIdRequired   = errors.New("GoogleUser 'id' is required")
	ErrGoogleUserNameRequired = errors.New("GoogleUser 'name' is required")
)

type GoogleUser struct {
	id    string
	name  string
	email Email
}

func NewGoogleUser(id, name, email string) (*GoogleUser, error) {
	if id == "" {
		return nil, ErrGoogleUserIdRequired
	}
	if name == "" {
		return nil, ErrGoogleUserNameRequired
	}
	mail, err := NewEmail(email)
	if err != nil {
		return nil, err
	}
	return &GoogleUser{
		id:    id,
		name:  name,
		email: mail,
	}, nil
}

func (u *GoogleUser) Id() string {
	return u.id
}

func (u *GoogleUser) Name() string {
	return u.name
}

func (u *GoogleUser) Email() Email {
	return u.email
}
