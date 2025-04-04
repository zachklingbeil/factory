package factory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	_ "github.com/lib/pq"
	"github.com/zachklingbeil/factory/cmd"
)

type Factory struct {
	Ctx  context.Context
	Db   *sql.DB
	Eth  *ethclient.Client
	Http *http.Client
	Rpc  *rpc.Client
	Json *cmd.JSON
	Mu   sync.Mutex
}

// NewFactory initializes the Factory with all required components, including the database connection.
func NewFactory(dbName string) (*Factory, error) {
	ctx := context.Background()
	http := &http.Client{}

	// Initialize Ethereum RPC and Eth client
	rpc, eth, err := cmd.Node(ctx)
	if err != nil {
		return nil, err
	}

	// Initialize the database connection
	db, err := database(dbName)
	if err != nil {
		return nil, err
	}

	// Initialize JSON helper
	json := cmd.Json(*http, ctx)

	return &Factory{
		Rpc:  rpc,
		Eth:  eth,
		Http: http,
		Json: json,
		Ctx:  ctx,
		Db:   db,
	}, nil
}

func database(dbName string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=postgres password=postgres dbname=%s host=postgres port=5432 sslmode=disable", dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database '%s': %w", dbName, err)
	}

	const maxRetries = 5
	for i := 1; i <= maxRetries; i++ {
		if err := db.Ping(); err == nil {
			fmt.Printf("Connected to database '%s'\n", dbName)
			return db, nil
		}
		fmt.Printf("Retrying connection to database '%s' (%d/%d)...\n", dbName, i, maxRetries)
		time.Sleep(time.Second * time.Duration(i*2))
	}
	return nil, fmt.Errorf("failed to connect to database '%s' after %d retries", dbName, maxRetries)
}

// DiskToMem converts tables into slices of structs.
func (f *Factory) DiskToMem(table string, result any) error {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	rows, err := f.Db.Query(query)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	// Prepare a slice to hold JSON objects
	var jsonRows []map[string]any

	for rows.Next() {
		// Create a map to hold the row data
		rowMap := make(map[string]any, len(cols))
		values := make([]any, len(cols))
		valuePtrs := make([]any, len(cols))

		// Assign pointers to the values
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		// Map the column names to their corresponding values
		for i, col := range cols {
			rowMap[col] = values[i]
		}

		// Append the row map to the JSON rows slice
		jsonRows = append(jsonRows, rowMap)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iteration error: %w", err)
	}

	// Marshal the JSON rows slice into a JSON array
	jsonData, err := json.Marshal(jsonRows)
	if err != nil {
		return fmt.Errorf("failed to marshal results to JSON: %w", err)
	}

	// Unmarshal the JSON data into the provided result slice
	if err := json.Unmarshal(jsonData, result); err != nil {
		return fmt.Errorf("failed to unmarshal JSON into result: %w", err)
	}
	return nil
}
