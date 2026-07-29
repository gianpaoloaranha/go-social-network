package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/gianpaoloaranha/go-social-network/internal/infra/config"
)

func Connect(cfg config.Config) (*redis.Client, func(), error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		_ = client.Close()
		return nil, nil, fmt.Errorf("connect to redis: %w", err)
	}

	closeRedis := func() {
		_ = client.Close()
	}

	return client, closeRedis, nil
}
