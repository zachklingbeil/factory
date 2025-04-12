package fx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Database struct {
	*sql.DB
	Redis *redis.Client
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

func ConnectRedis(password string, ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:          "redis:6379",
		Password:      password,
		DB:            0,
		Protocol:      3,
		UnstableResp3: true,
	})
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	fmt.Println("Connected to Redis")
	return client, nil
}

func Connect(dbName, redispassword string, ctx context.Context) (*Database, error) {
	db, err := ConnectPostgres(dbName)
	if err != nil {
		return nil, err
	}

	redis, err := ConnectRedis(redispassword, ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to PostgreSQL and Redis successfully")
	return &Database{
		DB:    db,
		Redis: redis,
	}, nil
}
