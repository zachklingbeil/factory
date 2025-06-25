package fx

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(dbNumber int, ctx context.Context) (*redis.Client, error) {
	pw := os.Getenv("PASSWORD")
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: pw,
		DB:       dbNumber,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
