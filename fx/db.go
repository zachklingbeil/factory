package fx

import (
	"context"
	"database/sql"
	"fmt"
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
