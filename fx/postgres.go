package fx

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func ConnectPostgres(dbName string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=postgres password=postgres dbname=%s host=postgres port=5432 sslmode=disable", dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database '%s': %w", dbName, err)
	}

	const maxRetries = 5
	for i := 1; i <= maxRetries; i++ {
		if err := db.Ping(); err == nil {
			return db, nil
		}
		fmt.Printf("Retrying connection to database '%s' (%d/%d)...\n", dbName, i, maxRetries)
		time.Sleep(time.Second * time.Duration(i*2))
	}

	db.Close()
	return nil, fmt.Errorf("failed to connect to database '%s' after %d retries", dbName, maxRetries)
}
