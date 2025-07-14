package redis

import (
	"context"

	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
	"github.com/redis/go-redis/v9"
)

const (
	REDIS_TOKEN_WHITELIST_NAME = "whitelist:token:list"
)

type RedisTokenRepository struct {
	client *redis.Client
}

func NewRedisTokenRepository(client *redis.Client) port.TokenRepository {
	return &RedisTokenRepository{
		client: client,
	}
}

// ホワイトリストに追加
func (r *RedisTokenRepository) AddWhitelist(refreshToken *entity.RefreshToken) error {
	result := r.client.SAdd(context.Background(), REDIS_TOKEN_WHITELIST_NAME, refreshToken.Id().String())
	return result.Err()
}

// ホワイトリストから削除
func (r *RedisTokenRepository) RemoveWhitelist(tokenId uuid.UUID) error {
	result := r.client.SRem(context.Background(), REDIS_TOKEN_WHITELIST_NAME, tokenId.String())
	return result.Err()
}

// ホワイトリストに存在するか確認
func (r *RedisTokenRepository) InWhitelist(tokenId uuid.UUID) (bool, error) {
	result, err := r.client.SIsMember(context.Background(), REDIS_TOKEN_WHITELIST_NAME, tokenId.String()).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}
