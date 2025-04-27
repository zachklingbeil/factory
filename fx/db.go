package fx

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Data struct {
	Pg  *sql.DB
	RB  *redis.Client
	Ctx context.Context
}

func Source(dbName string, ctx context.Context) (*Data, error) {
	db := &Data{
		Ctx: ctx,
	}
	db.ConnectPostgres(dbName)
	db.ConnectRedis(0, ctx)
	return db, nil
}

func (d *Data) ConnectRedis(dbNumber int, ctx context.Context) error {
	pw := os.Getenv("PASSWORD")
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: pw,
		DB:       dbNumber,
	})
	if _, err := client.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	d.RB = client
	return nil
}

func (d *Data) ConnectPostgres(dbName string) error {
	connStr := fmt.Sprintf("user=postgres password=postgres dbname=%s host=postgres port=5432 sslmode=disable", dbName)
	pg, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open connection to Data '%s': %w", dbName, err)
	}

	const maxRetries = 5
	for i := 1; i <= maxRetries; i++ {
		if err := pg.Ping(); err == nil {
			return nil
		}
		fmt.Printf("Retrying connection to Data '%s' (%d/%d)...\n", dbName, i, maxRetries)
		time.Sleep(time.Second * time.Duration(i*2))
	}
	d.Pg = pg
	return fmt.Errorf("failed to connect to Data '%s' after %d retries", dbName, maxRetries)
}

func (d *Data) Save(key string, source []any) error {
	err := StoreSlice(key, source, d.RB, d.Ctx)
	if err != nil {
		return fmt.Errorf("failed to store slice using Redis source: %w", err)
	}
	return nil
}

func (d *Data) Source(key string) []any {
	var result []any
	items := SourceSlice[any](d.Ctx, d.RB, key)
	result = append(result, items...)
	return result
}

func SourceSlice[T any](ctx context.Context, rb *redis.Client, key string) []T {
	var items []T
	source, err := rb.SMembers(ctx, key).Result()
	if err != nil {
		log.Fatalf("Failed to fetch items from Redis set '%s': %v", key, err)
	}

	for _, s := range source {
		var item T
		if err := json.Unmarshal([]byte(s), &item); err != nil {
			log.Printf("Skipping invalid item: %v (data: %s)", err, s)
			continue
		}
		items = append(items, item)
	}
	return items
}

// StoreRedisSlice stores a slice of items into Redis as a set.
func StoreSlice[T any](key string, source []T, redis *redis.Client, ctx context.Context) error {
	pipe := redis.Pipeline()
	for _, item := range source {
		data, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("failed to marshal item for key '%s': %w", key, err)
		}
		pipe.SAdd(ctx, key, data)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to store slice in Redis for key '%s': %w", key, err)
	}
	return nil
}
