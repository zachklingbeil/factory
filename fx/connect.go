package fx

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Establish geth.ipc connection
func (f *Fx) Node() (*rpc.Client, *ethclient.Client) {
	rpc, err := rpc.DialIPC(f.Ctx, "/ethereum/geth.ipc") // Updated path
	if err != nil {
		log.Printf("Failed to connect to the Ethereum client: %v", err)
		return nil, nil
	}
	// log.Println("Successfully connected to the Ethereum client.")
	eth := ethclient.NewClient(rpc)
	f.rpc = rpc
	f.eth = eth
	return rpc, eth
}

// Establish geth.ws connection using API key from environment variable
func (f *Fx) NodeWS(wsURL, apikey string) (*rpc.Client, *ethclient.Client, error) {
	fullURL := fmt.Sprintf("%s/%s", wsURL, apikey)
	rpcClient, err := rpc.DialContext(f.Ctx, fullURL)
	if err != nil {
		log.Printf("Failed to connect to Ethereum WebSocket: %v", err)
		return nil, nil, err
	}
	eth := ethclient.NewClient(rpcClient)
	f.rpc = rpcClient
	f.eth = eth
	return rpcClient, eth, nil
}

// Establish geth.http connection using API key from environment variable
func (f *Fx) NodeHTTP(httpURL, apikey string) (*rpc.Client, *ethclient.Client, error) {
	fullURL := fmt.Sprintf("%s/%s", httpURL, apikey)
	rpcClient, err := rpc.DialHTTP(fullURL)
	if err != nil {
		log.Printf("Failed to connect to Ethereum HTTP: %v", err)
		return nil, nil, err
	}
	eth := ethclient.NewClient(rpcClient)
	f.rpc = rpcClient
	f.eth = eth
	return rpcClient, eth, nil
}

func (f *Fx) ConnectRedis(dbNumber int, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: password,
		DB:       dbNumber,
	})

	if _, err := client.Ping(f.Ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	f.redis = client
	return client, nil
}

func (f *Fx) ConnectPostgres(dbName string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=postgres password=postgres dbname=%s host=postgres port=5432 sslmode=disable", dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database '%s': %w", dbName, err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database '%s': %w", dbName, err)
	}
	f.postgres = db
	return db, nil
}

// NewOAuthClient returns an authenticated HTTP client using OAuth2 client credentials flow
// with automatic token refreshing and all requested scopes.
func (f *Fx) NewOAuthClient(clientID, clientSecret, tokenURL string, scopes []string) (*http.Client, error) {
	ctx, cancel := context.WithTimeout(f.Ctx, 2*time.Minute)
	defer cancel()

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       scopes,
	}

	// Create a TokenSource that automatically refreshes tokens
	tokenSource := config.TokenSource(ctx)

	// Get initial token to verify credentials
	token, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("OAuth client credentials flow failed: %w", err)
	}
	fmt.Printf("✓ OAuth client authenticated successfully (token expires: %v)\n", token.Expiry)

	// Create an HTTP client that uses the TokenSource for automatic refreshing
	client := oauth2.NewClient(ctx, tokenSource)
	return client, nil
}
