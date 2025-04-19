package fx

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Database struct {
	*sql.DB
	Redis *redis.Client
	Mu    *sync.Mutex
	Rw    *sync.RWMutex
	Ctx   *context.Context
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

// Connect creates a Database struct with a Redis client for the specified DB number.
func Connect(ctx context.Context, dbName string) (*Database, error) {
	db, err := ConnectPostgres(dbName)
	if err != nil {
		return nil, err
	}

	redis, err := ConnectRedis(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Connected to PostgreSQL and Redis (DB) successfully\n")
	return &Database{
		DB:    db,
		Redis: redis,
	}, nil
}
