package datasource

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     getEndpoint(),
		Password: getPassword(),
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic("redis connection failed")
	}
	return client
}

func getEndpoint() string {
	address, ok := os.LookupEnv("REDIS_ADDRESS")
	if !ok {
		panic("\"REDIS_ADDRESS\" is not set")
	}
	return address
}

func getPassword() string {
	password, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		panic("\"REDIS_PASSWORD\" is not set")
	}
	return password
}
