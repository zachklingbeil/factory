package fx

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(ctx context.Context) (*redis.Client, error) {
	pw := os.Getenv("REDIS_PASSWORD")
	client := redis.NewClient(&redis.Options{
		Addr:          "redis:6379",
		Password:      pw,
		DB:            0,
		Protocol:      3,
		UnstableResp3: true,
	})
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	fmt.Println("redis [0]")
	return client, nil
}

// Loop loads all data from the specified Redis DB into value, then starts a goroutine to periodically save value back to Redis.
func (db *Database) Loop(distance time.Duration) error {
	value, err := db.Past()
	if err != nil {
		return fmt.Errorf("failed to load data from Redis: %w", err)
	}

	go db.Present(value, distance)
	return nil
}

// Present saves the given map to Redis every distance until db.Ctx is cancelled, using the specified Redis DB.
func (db *Database) Present(value map[any]any, distance time.Duration) {
	ticker := time.NewTicker(distance)
	defer ticker.Stop()
	ctx := *db.Ctx
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := db.Redis.Do(ctx, "SELECT", 0).Err(); err != nil {
				continue
			}
			db.Rw.RLock()
			for key, value := range value {
				keyJSON, err := json.Marshal(key)
				if err != nil {
					continue
				}
				valueJSON, err := json.Marshal(value)
				if err != nil {
					continue
				}
				redisKey := string(keyJSON)
				_ = db.Redis.Set(ctx, redisKey, valueJSON, 0).Err()
			}
			db.Rw.RUnlock()
		}
	}
}

// Past loads all keys from Redis and returns a map[any]any.
func (db *Database) Past() (map[any]any, error) {
	ctx := *db.Ctx
	if err := db.Redis.Do(ctx, "SELECT", 0).Err(); err != nil {
		return nil, err
	}
	result := make(map[any]any)
	db.Mu.Lock()
	defer db.Mu.Unlock()
	iter := db.Redis.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		redisKey := iter.Val()
		var key any
		if err := json.Unmarshal([]byte(redisKey), &key); err != nil {
			continue
		}
		valueJSON, err := db.Redis.Get(ctx, redisKey).Bytes()
		if err != nil {
			return nil, err
		}
		var value any
		if err := json.Unmarshal(valueJSON, &value); err != nil {
			return nil, err
		}
		result[key] = value
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
