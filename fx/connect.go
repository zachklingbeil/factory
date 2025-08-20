package fx

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
func (f *Fx) Node() error {
	rpc, err := rpc.DialIPC(f.Ctx, "/.ethereum/geth.ipc") // Updated path
	if err != nil {
		log.Printf("Failed to connect to the Ethereum client: %v", err)
		return nil
	}
	// log.Println("Successfully connected to the Ethereum client.")
	eth := ethclient.NewClient(rpc)
	f.Rpc = rpc
	f.Eth = eth
	return nil
}

// GethHandler handles JSON-RPC requests via IPC
func (f *Fx) GethHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Parse the JSON-RPC request
	var req struct {
		Method string `json:"method"`
		Params []any  `json:"params"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON-RPC request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(f.Ctx, 30*time.Second)
	defer cancel()

	var result json.RawMessage
	if err := f.Rpc.CallContext(ctx, &result, req.Method, req.Params...); err != nil {
		http.Error(w, "RPC call failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(result)
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
