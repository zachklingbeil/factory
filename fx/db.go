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

// DiskToMem converts tables into slices of structs.
func (d *Database) DiskToMem(table string, result any) error {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	rows, err := d.Query(query)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	colCount := len(cols)
	jsonRows := make([]map[string]any, 0, 10000)
	values := make([]any, colCount)
	valuePtrs := make([]any, colCount)

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		rowMap := make(map[string]any, colCount)

		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		for i, col := range cols {
			rowMap[col] = values[i]
		}
		jsonRows = append(jsonRows, rowMap)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iteration error: %w", err)
	}

	jsonData, err := json.Marshal(jsonRows)
	if err != nil {
		return fmt.Errorf("failed to marshal results to JSON: %w", err)
	}
	if err := json.Unmarshal(jsonData, result); err != nil {
		return fmt.Errorf("failed to unmarshal JSON into result: %w", err)
	}
	return nil
}

func (d *Database) ColumnToSlice(table string, column string, result any) error {
	query := fmt.Sprintf("SELECT %s FROM %s", column, table)
	rows, err := d.Query(query)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	slice := make([]any, 0, 250000) // Preallocate a slice with an initial capacity
	for rows.Next() {
		var value any
		if err := rows.Scan(&value); err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}
		slice = append(slice, value)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iteration error: %w", err)
	}

	// Marshal the slice into JSON and unmarshal it into the provided result
	jsonData, err := json.Marshal(slice)
	if err != nil {
		return fmt.Errorf("failed to marshal results to JSON: %w", err)
	}
	if err := json.Unmarshal(jsonData, result); err != nil {
		return fmt.Errorf("failed to unmarshal JSON into result: %w", err)
	}
	return nil
}
