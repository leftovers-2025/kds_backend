package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	MIN_USER_NAME_LENGTH = 1
	MAX_USER_NAME_LENGTH = 64
)

var (
	ErrUserNameRequired      = errors.New("User 'name' is required")
	ErrUserGoogleIdRequired  = errors.New("User 'googleId' is required")
	ErrInvalidUserNameLength = errors.New("invalid User 'name' length")
	ErrUserTimeZero          = errors.New("the time related User is zero")
)

type User struct {
	id        uuid.UUID
	name      string
	email     Email
	google_id string
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(id uuid.UUID, name, email, googleId string, createdAt, updatedAt time.Time) (*User, error) {
	if name == "" {
		return nil, ErrUserNameRequired
	}
	if googleId == "" {
		return nil, ErrUserGoogleIdRequired
	}
	mail, err := NewEmail(email)
	if err != nil {
		return nil, err
	}
	if len(name) < MIN_USER_NAME_LENGTH || len(name) > MAX_USER_NAME_LENGTH {
		return nil, ErrInvalidUserNameLength
	}
	if createdAt.IsZero() || updatedAt.IsZero() {
		return nil, ErrUserTimeZero
	}
	return &User{
		id:        id,
		name:      name,
		email:     mail,
		google_id: googleId,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (u *User) Id() uuid.UUID {
	return u.id
}

func (u *User) GoogleId() string {
	return u.google_id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
