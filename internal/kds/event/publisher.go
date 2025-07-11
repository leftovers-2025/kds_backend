package event

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
)

var (
	ErrEventTypeUnknown = errors.New("event type is unknown")
)

type EventPublisher interface {
	Publish(Event) error
}

type RedisEventPublisher struct {
	client *redis.Client
}

func NewRedisEventPublisher(client *redis.Client) EventPublisher {
	return &RedisEventPublisher{
		client: client,
	}
}

// Redisにイベントを追加
func (p *RedisEventPublisher) Publish(event Event) error {
	// UNKNOWNの場合は処理をしない
	if event.Type() == EVENT_UNKNOWN {
		return ErrEventTypeUnknown
	}
	// イベントのデータを取得
	eventData, err := event.Event()
	if err != nil {
		return err
	}
	jsonBody := []byte{}
	// イベントデータをJSON化
	err = json.Unmarshal(jsonBody, eventData)
	if err != nil {
		return err
	}
	jsonString := string(jsonBody)
	// Redisにアップロード
	result := p.client.SAdd(context.Background(), event.Type().String(), jsonString)
	return result.Err()
}
