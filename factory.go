package factory

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/redis/go-redis/v9"
	"github.com/zachklingbeil/factory/fx"
	"golang.org/x/oauth2/clientcredentials"
)

type Factory struct {
	Ctx  context.Context
	Mu   *sync.Mutex
	Rw   *sync.RWMutex
	When *sync.Cond
	Api  *fx.API
}

func Assemble() *Factory {
	ctx := context.Background()
	mu := &sync.Mutex{}
	rw := &sync.RWMutex{}
	when := sync.NewCond(mu)
	factory := &Factory{
		Ctx:  ctx,
		Mu:   mu,
		Rw:   rw,
		When: when,
	}
	return factory
}

// Establish geth.ipc connection
func (f *Factory) Node() (*rpc.Client, *ethclient.Client) {
	rpc, err := rpc.DialIPC(f.Ctx, "/ethereum/geth.ipc") // Updated path
	if err != nil {
		log.Printf("Failed to connect to the Ethereum client: %v", err)
		return nil, nil
	}
	// log.Println("Successfully connected to the Ethereum client.")
	eth := ethclient.NewClient(rpc)
	return rpc, eth
}

// Establish geth.ws connection using API key from environment variable
func (f *Factory) NodeWS(wsURL string) (*rpc.Client, *ethclient.Client, error) {
	fullURL := fmt.Sprintf("%s/%s", wsURL, os.Getenv("ETH_API_KEY"))
	rpcClient, err := rpc.DialContext(f.Ctx, fullURL)
	if err != nil {
		log.Printf("Failed to connect to Ethereum WebSocket: %v", err)
		return nil, nil, err
	}
	eth := ethclient.NewClient(rpcClient)
	return rpcClient, eth, nil
}

// Establish geth.http connection using API key from environment variable
func (f *Factory) NodeHTTP(httpURL string) (*rpc.Client, *ethclient.Client, error) {
	fullURL := fmt.Sprintf("%s/%s", httpURL, os.Getenv("ETH_API_KEY"))
	rpcClient, err := rpc.DialHTTP(fullURL)
	if err != nil {
		log.Printf("Failed to connect to Ethereum HTTP: %v", err)
		return nil, nil, err
	}
	eth := ethclient.NewClient(rpcClient)
	return rpcClient, eth, nil
}

func (f *Factory) ConnectRedis(dbNumber int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: os.Getenv("PASSWORD"),
		DB:       dbNumber,
	})

	if _, err := client.Ping(f.Ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

func (f *Factory) ConnectPostgres(dbName string) (*sql.DB, error) {
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

// NewOAuthClient returns an authenticated HTTP client (machine-to-machine, no user interaction)
func (f *Factory) NewOAuthClient(clientID, clientSecret, tokenURL string, scopes []string) (*http.Client, error) {
	ctx, cancel := context.WithTimeout(f.Ctx, 2*time.Minute)
	defer cancel()

	clientConfig := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       scopes,
	}

	// Get token and create HTTP client
	client := clientConfig.Client(ctx)
	if client == nil {
		return nil, fmt.Errorf("failed to create OAuth client")
	}

	// Test the client by making a token request to validate credentials
	token, err := clientConfig.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("OAuth client credentials flow failed: %w", err)
	}
	fmt.Printf("✓ OAuth client authenticated successfully (token expires: %v)\n", token.Expiry)
	return client, nil
}
