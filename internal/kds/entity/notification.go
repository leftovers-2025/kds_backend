package entity

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNotificationUserIdRequired = errors.New("notification userId required")
)

type Notification struct {
	userId    uuid.UUID
	enabled   bool
	locations []Location
	tags      []Tag
}

func NewNotificaton(userId uuid.UUID, enabled bool, locations []Location, tags []Tag) (*Notification, error) {
	if userId == uuid.Nil {
		return nil, ErrNotificationUserIdRequired
	}
	return &Notification{
		userId:    userId,
		enabled:   enabled,
		locations: locations,
		tags:      tags,
	}, nil
}

func (n *Notification) UserId() uuid.UUID {
	return n.userId
}

func (n *Notification) IsEnabled() bool {
	return n.enabled
}

func (n *Notification) Locations() []Location {
	return n.locations
}

func (n *Notification) Tags() []Tag {
	return n.tags
}
