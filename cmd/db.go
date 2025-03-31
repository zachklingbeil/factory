// Simplified logic for managing multiple database connections with retry support.
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	Conn map[string]*sql.DB
	Mu   *sync.Mutex
}

// NewDatabase initializes a new Database instance.
func NewDatabase() *Database {
	return &Database{
		Conn: make(map[string]*sql.DB),
	}
}

// Connect establishes or retrieves a connection to a PostgreSQL database with retry logic.
func (d *Database) Connect(dbName string) (*sql.DB, error) {
	d.Mu.Lock()
	defer d.Mu.Unlock()

	// Return existing connection if it exists
	if db, exists := d.Conn[dbName]; exists {
		log.Printf("Reusing existing connection to database '%s'.", dbName)
		return db, nil
	}

	connStr := fmt.Sprintf("user=postgres password=postgres dbname=%s host=postgres port=5432 sslmode=disable", dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database '%s': %w", dbName, err)
	}

	const maxRetries = 5
	for i := 1; i <= maxRetries; i++ {
		if err := db.Ping(); err == nil {
			log.Printf("Connection to database '%s' established successfully.", dbName)
			d.Conn[dbName] = db
			return db, nil
		}
		log.Printf("Retrying connection to database '%s' (%d/%d)...", dbName, i, maxRetries)
		time.Sleep(time.Second * time.Duration(i*2))
	}
	return nil, fmt.Errorf("failed to connect to database '%s' after %d retries", dbName, maxRetries)
}

// CloseAll closes all database Conn.
func (d *Database) CloseAll() {
	d.Mu.Lock()
	defer d.Mu.Unlock()

	for dbName, db := range d.Conn {
		if err := db.Close(); err != nil {
			log.Printf("Error closing connection to database '%s': %v", dbName, err)
		} else {
			log.Printf("Connection to database '%s' closed.", dbName)
		}
		delete(d.Conn, dbName)
	}
}
