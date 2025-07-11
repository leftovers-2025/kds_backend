package event

type EventType string

const (
	EVENT_UNKNOWN      EventType = "event:unknown"
	EVENT_IMAGE_UPLOAD EventType = "event:images:upload"
)

type Event interface {
	Type() EventType
	Event() (any, error)
}

func (t EventType) String() string {
	return string(t)
}
