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

func ConnectPostgres(dbName string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=postgres password=postgres dbname=%s host=postgres port=5432 sslmode=disable", dbName)
	pg, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database '%s': %w", dbName, err)
	}

	const maxRetries = 5
	for i := 1; i <= maxRetries; i++ {
		if err := pg.Ping(); err == nil {
			return pg, nil
		}
		fmt.Printf("Retrying connection to database '%s' (%d/%d)...\n", dbName, i, maxRetries)
		time.Sleep(time.Second * time.Duration(i*2))
	}
	return nil, fmt.Errorf("failed to connect to database '%s' after %d retries", dbName, maxRetries)
}

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

// func Connect(dbName string, dbNum int, ctx context.Context, mu *sync.Mutex, rw *sync.RWMutex) (*Database, error) {
// 	db := &Database{
// 		Mu:  mu,
// 		Rw:  rw,
// 		Ctx: ctx,
// 	}

// 	// Connect to Postgres
// 	if err := db.ConnectPostgres(dbName); err != nil {
// 		return nil, err
// 	}

// 	// Connect to Redis
// 	client, err := db.ConnectRedis(dbNum)
// 	if err != nil {
// 		return nil, err
// 	}
// 	db.Rdb = client // Assign the Redis client to the Database struct

// 	return db, nil
// }
