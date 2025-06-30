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
	ErrInvalidUserNameLength = errors.New("invalid username length")
	ErrUserTimeZero          = errors.New("the user time is zero")
)

type User struct {
	id        uuid.UUID
	name      string
	email     string
	google_id string
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(id uuid.UUID, name, email, google_id string, createdAt, updatedAt time.Time) (*User, error) {
	if len(name) < MIN_USER_NAME_LENGTH || len(name) > MAX_USER_NAME_LENGTH {
		return nil, ErrInvalidUserNameLength
	}
	if createdAt.IsZero() || updatedAt.IsZero() {
		return nil, ErrUserTimeZero
	}
	return &User{
		id:        id,
		name:      name,
		email:     email,
		google_id: google_id,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}
