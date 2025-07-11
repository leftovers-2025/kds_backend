package entity

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTagIdRequired   = errors.New("tag id is required")
	ErrTagNameRequired = errors.New("tag name is required")
)

type Tag struct {
	id   uuid.UUID
	name string
}

func NewTag(id uuid.UUID, name string) (*Tag, error) {
	if id == uuid.Nil {
		return nil, ErrTagIdRequired
	}
	if name == "" {
		return nil, ErrTagNameRequired
	}
	return &Tag{
		id:   id,
		name: name,
	}, nil
}

func (e *Tag) Id() uuid.UUID {
	return e.id
}

func (e *Tag) Name() string {
	return e.name
}
