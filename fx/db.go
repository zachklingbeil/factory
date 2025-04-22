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
	pw := os.Getenv("PASSWORD")
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

// // Loop loads all data from Redis into a map, starts Present as a goroutine, and returns the map.
// func (db *Database) Loop(distance time.Duration) error {
// 	if err := db.Past(); err != nil {
// 		return fmt.Errorf("failed to load data from Redis: %w", err)
// 	}
// 	go db.Present(distance)
// 	return nil
// }

// // Present saves the Circuit's One map to Redis every distance until the context is cancelled.
// func (db *Database) Present(distance time.Duration) {
// 	ticker := time.NewTicker(distance)
// 	defer ticker.Stop()
// 	ctx := db.Ctx
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-ticker.C:
// 			if err := db.Rdb.Do(ctx, "SELECT", 0).Err(); err != nil {
// 				continue
// 			}
// 			db.Rw.RLock()
// 			for key, value := range db.One {
// 				keyJSON, err := json.Marshal(key)
// 				if err != nil {
// 					continue
// 				}
// 				valueJSON, err := json.Marshal(value)
// 				if err != nil {
// 					continue
// 				}
// 				redisKey := string(keyJSON)
// 				_ = db.Rdb.Set(ctx, redisKey, valueJSON, 0).Err()
// 			}
// 			db.Rw.RUnlock()
// 		}
// 	}
// }

// // Past loads all keys from Redis into the Circuit's One map.
// func (db *Database) Past() error {
// 	if err := db.Rdb.Do(db.Ctx, "SELECT", 0).Err(); err != nil {
// 		return err
// 	}
// 	db.Mu.Lock()
// 	defer db.Mu.Unlock()
// 	iter := db.Rdb.Scan(db.Ctx, 0, "*", 0).Iterator()
// 	for iter.Next(db.Ctx) {
// 		redisKey := iter.Val()
// 		var key Zero
// 		if err := json.Unmarshal([]byte(redisKey), &key); err != nil {
// 			continue
// 		}
// 		valueJSON, err := db.Rdb.Get(db.Ctx, redisKey).Bytes()
// 		if err != nil {
// 			return err
// 		}
// 		var value any
// 		if err := json.Unmarshal(valueJSON, &value); err != nil {
// 			return err
// 		}
// 		db.One[key] = value
// 	}
// 	if err := iter.Err(); err != nil {
// 		return err
// 	}
// 	return nil
// }
