// Logic for connecting an application to a db within a Docker network (--driver bridge).
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Attempts to connect to the db using a hardcoded connection string:
// "user=postgres password=postgres dbname=postgres host=postgres port=5432 sslmode=disable"
// Returns:
//
//	*sql.DB: A database connection.
//
// or
// error: An error if the connection cannot be established.
func Database(dbName string) *sql.DB {
	// Format the connection string with the provided database name
	conn := fmt.Sprintf("user=postgres password=postgres dbname=%s host=postgres port=5432 sslmode=disable", dbName)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	maxRetries := 5
	for i := 1; i <= maxRetries; i++ {
		if err := db.Ping(); err == nil {
			log.Println("Database connection established successfully.")
			return db
		}
		log.Printf("Retrying database connection (%d/%d)...", i, maxRetries)
		time.Sleep(time.Second * time.Duration(i*2))
	}
	log.Fatalf("Error connecting to the database after %d retries.", maxRetries)
	return nil
}

func Close(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing the database connection: %v", err)
		} else {
			log.Println("Database connection closed successfully.")
		}
	}
}
