package fx

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Database struct {
	*sql.DB
}

func NewDatabase(dbName string) (*Database, error) {
	connStr := fmt.Sprintf("user=postgres password=postgres dbname=%s host=postgres port=5432 sslmode=disable", dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database '%s': %w", dbName, err)
	}

	const maxRetries = 5
	for i := 1; i <= maxRetries; i++ {
		if err := db.Ping(); err == nil {
			// fmt.Printf("Connected to database '%s'\n", dbName)
			return &Database{DB: db}, nil
		}
		fmt.Printf("Retrying connection to database '%s' (%d/%d)...\n", dbName, i, maxRetries)
		time.Sleep(time.Second * time.Duration(i*2))
	}
	return nil, fmt.Errorf("failed to connect to database '%s' after %d retries", dbName, maxRetries)
}

func (db *Database) Consolidate() error {
	newDB, err := NewDatabase("timefactory")
	if err != nil {
		return fmt.Errorf("failed to create or connect to database 'timefactory': %w", err)
	}
	defer newDB.Close()

	// Create the table with columns key, value (as INTEGER), and jsonb
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS data (
            key TEXT PRIMARY KEY,
            value INTEGER NOT NULL,
            jsonb_data JSONB
        );
    `
	_, err = newDB.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create table in database 'timefactory': %w", err)
	}

	fmt.Println("Database 'timefactory' and table 'data' created successfully.")
	return nil
}

func (db *Database) Insert(key string, data []map[string]any) error {
	// Check if the key already exists and get the current value
	selectQuery := `SELECT value FROM data WHERE key = $1;`
	var currentValue int
	err := db.QueryRow(selectQuery, key).Scan(&currentValue)
	if err != nil {
		if err != sql.ErrNoRows {
			return fmt.Errorf("failed to query current value for key '%s': %w", key, err)
		}
		// If no rows exist, start with value 0
		currentValue = 0
	}

	// Increment the value
	newValue := currentValue + 1

	// Marshal the slice of structs into JSON
	jsonbData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal jsonb data for key '%s': %w", key, err)
	}

	// Insert or update the row in the database
	insertQuery := `
        INSERT INTO data (key, value, jsonb_data)
        VALUES ($1, $2, $3)
        ON CONFLICT (key) DO UPDATE
        SET value = $2, jsonb_data = EXCLUDED.jsonb_data;
    `
	_, err = db.Exec(insertQuery, key, newValue, jsonbData)
	if err != nil {
		return fmt.Errorf("failed to insert or update data for key '%s': %w", key, err)
	}

	fmt.Printf("Data for key '%s' inserted/updated successfully with value '%d'.\n", key, newValue)
	return nil
}

func (db *Database) GetData(key string) ([]map[string]any, int, error) {
	selectQuery := `SELECT value, jsonb_data FROM data WHERE key = $1;`

	row := db.QueryRow(selectQuery, key)

	var value int
	var jsonbData []byte

	if err := row.Scan(&value, &jsonbData); err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, fmt.Errorf("no data found for key '%s'", key)
		}
		return nil, 0, fmt.Errorf("failed to scan row for key '%s': %w", key, err)
	}

	var data []map[string]any
	if err := json.Unmarshal(jsonbData, &data); err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal jsonb data for key '%s': %w", key, err)
	}

	return data, value, nil
}
