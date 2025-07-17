package event

type EventType string

const (
	EVENT_UNKNOWN    EventType = "event:unknown"
	EVENT_GRANT_ROOT EventType = "event:grant:root"
)

type Event interface {
	Type() EventType
	Event() (any, error)
}

func (t EventType) String() string {
	return string(t)
}
