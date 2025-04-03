package factory

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	_ "github.com/lib/pq" // PostgreSQL driver
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
