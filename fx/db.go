package fx

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Database struct {
	Pg  *sql.DB
	Rdb *redis.Client
	Mu  *sync.Mutex
	Rw  *sync.RWMutex
	Ctx context.Context
}

func Connect(dbName string, distance time.Duration, ctx context.Context, mu *sync.Mutex, rw *sync.RWMutex) (*Database, error) {
	pg, err := ConnectPostgres(dbName)
	if err != nil {
		return nil, err
	}

	rdb, err := ConnectRedis(ctx)
	if err != nil {
		return nil, err
	}

	return &Database{
		Pg:  pg,
		Rdb: rdb,
		Mu:  mu,
		Rw:  rw,
		Ctx: ctx,
	}, nil
}

func ConnectPostgres(dbName string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=postgres password=postgres dbname=%s host=postgres port=5432 sslmode=disable", dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database '%s': %w", dbName, err)
	}

	const maxRetries = 5
	for i := 1; i <= maxRetries; i++ {
		if err := db.Ping(); err == nil {
			// fmt.Printf("Connected to database '%s'\n", dbName)
			return db, nil
		}
		fmt.Printf("Retrying connection to database '%s' (%d/%d)...\n", dbName, i, maxRetries)
		time.Sleep(time.Second * time.Duration(i*2))
	}
	return nil, fmt.Errorf("failed to connect to database '%s' after %d retries", dbName, maxRetries)
}

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
	return client, nil
}
