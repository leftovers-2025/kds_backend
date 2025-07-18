package email

import (
	"errors"
)

var (
	ErrEmailHostRequired     = errors.New("email host required")
	ErrEmailPortRequired     = errors.New("email port required")
	ErrEmailAddressRequired  = errors.New("email address required")
	ErrEmailPasswordRequired = errors.New("email password required")
)

type EmailAuth struct {
	host     string
	port     int
	address  string
	password string
}

func NewEmailAuth(host string, port int, address, password string) (*EmailAuth, error) {
	if host == "" {
		return nil, ErrEmailHostRequired
	}
	if port == 0 {
		return nil, ErrEmailPortRequired
	}
	if address == "" {
		return nil, ErrEmailAddressRequired
	}
	if password == "" {
		return nil, ErrEmailPasswordRequired
	}
	return &EmailAuth{
		host:     host,
		port:     port,
		address:  address,
		password: password,
	}, nil
}

func (e *EmailAuth) Host() string {
	return e.host
}

func (e *EmailAuth) Port() int {
	return e.port
}

func (e *EmailAuth) EmailAddress() string {
	return e.address
}

func (e *EmailAuth) Password() string {
	return e.password
}
