// Factory provides a common context for sourcing and distrubting data.
// Includes an Ethereum, HTTP, RPC client, a database connection, and json i/o logic.
package factory

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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

func NewFactory(dbName string) (*Factory, error) {
	ctx := context.Background()
	http := &http.Client{}

	rpc, eth, err := cmd.Node(ctx)
	if err != nil {
		return nil, err
	}

	db, err := NewDatabase(dbName)
	if err != nil {
		return nil, err
	}

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

// NewDatabase initializes a new Database instance.
func NewDatabase(dbName string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=postgres password=postgres dbname=%s host=postgres port=5432 sslmode=disable", dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database '%s': %w", dbName, err)
	}

	const maxRetries = 5
	for i := 1; i <= maxRetries; i++ {
		if err := db.Ping(); err == nil {
			log.Println("Connected to database")
			return &sql.DB{}, nil
		}

		log.Printf("Connection attempt %d/%d failed. Retrying in %ds...",
			i, maxRetries, i*2)
		time.Sleep(time.Duration(i*2) * time.Second)
	}
	return nil, fmt.Errorf("failed to connect to database '%s' after %d attempts", dbName, maxRetries)
}
