package event

import "github.com/redis/go-redis/v9"

type EventWorker struct {
	client *redis.Client
}

func NewEventWorker(client *redis.Client) *EventWorker {
	if client == nil {
		panic("nil Redis Client")
	}
	return &EventWorker{
		client: client,
	}
}

func (w *EventWorker) processImageUpload() {
}

func (w *EventWorker) grantRoot() {
}
