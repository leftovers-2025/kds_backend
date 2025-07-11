package entity

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrLocationIdRequired   = errors.New("locaiton id is required")
	ErrLocationNameRequired = errors.New("locaiton name is required")
)

type Location struct {
	id   uuid.UUID
	name string
}

var NilLocation Location

func NewLocation(id uuid.UUID, name string) (*Location, error) {
	if id == uuid.Nil {
		return &NilLocation, ErrLocationIdRequired
	}
	if name == "" {
		return &NilLocation, ErrLocationNameRequired
	}
	return &Location{
		id:   id,
		name: name,
	}, nil
}

func (e Location) Id() uuid.UUID {
	return e.id
}

func (e Location) Name() string {
	return e.name
}
