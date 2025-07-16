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
	ErrUserIdRequired        = errors.New("User 'id' is required")
	ErrUserNameRequired      = errors.New("User 'name' is required")
	ErrUserGoogleIdRequired  = errors.New("User 'googleId' is required")
	ErrUserRoleUnknown       = errors.New("User 'role' is unknown")
	ErrInvalidUserNameLength = errors.New("invalid User 'name' length")
	ErrUserTimeZero          = errors.New("the time related User is zero")
	ErrUserInvalidPermission = errors.New("user invalid permission")
)

type User struct {
	id        uuid.UUID
	name      string
	email     Email
	role      Role
	google_id string
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(id uuid.UUID, name, googleId string, email Email, role Role, createdAt, updatedAt time.Time) (*User, error) {
	if id == uuid.Nil {
		return nil, ErrUserIdRequired
	}
	if name == "" {
		return nil, ErrUserNameRequired
	}
	if googleId == "" {
		return nil, ErrUserGoogleIdRequired
	}
	if role == ROLE_UNKNOWN {
		return nil, ErrUserRoleUnknown
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
		email:     email,
		role:      role,
		google_id: googleId,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (u *User) Id() uuid.UUID {
	return u.id
}

func (u *User) Role() Role {
	return u.role
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

func (u *User) UpdateRole(role Role) error {
	if role == ROLE_UNKNOWN {
		return ErrUserRoleUnknown
	}
	u.role = role
	u.updatedAt = time.Now()
	return nil
}

func (u *User) CanSeeUsers() error {
	if u.role == ROLE_TEACHER || u.role == ROLE_ROOT {
		return nil
	}
	return ErrUserInvalidPermission
}
