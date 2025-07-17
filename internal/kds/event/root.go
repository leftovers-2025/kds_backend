package event

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrGrantRootUserIdRequired = errors.New("user id is required")
)

type GrantRootEvent struct {
	userId uuid.UUID
}

type GrantRoot struct {
	UserId string `json:"userId"`
}

func NewGrantRootEvent(userId uuid.UUID) (*GrantRootEvent, error) {
	if userId == uuid.Nil {
		return nil, ErrGrantRootUserIdRequired
	}
	return &GrantRootEvent{
		userId: userId,
	}, nil
}

func (e *GrantRootEvent) Event() (any, error) {
	return GrantRoot{
		UserId: e.userId.String(),
	}, nil
}

func (e *GrantRootEvent) Type() EventType {
	return EVENT_GRANT_ROOT
}
