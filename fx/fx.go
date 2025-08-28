package fx

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/zachklingbeil/factory/fx/io"
	"goauthentik.io/api/v3"
)

type Fx struct {
	Ctx      context.Context
	Rpc      *rpc.Client
	Eth      *ethclient.Client
	postgres *sql.DB
	redis    *redis.Client
	Auth     *api.APIClient // Authentik API client for management
	*io.IO
	apiKey string
}

func InitFx() *Fx {
	ctx := context.Background()
	fx := &Fx{
		Ctx:    ctx,
		IO:     io.NewIO(ctx),
		apiKey: os.Getenv("API_KEY"),
	}
	fx.Node()
	return fx
}

// NewAuth creates a new Authentik API client with the given baseURL and apikey.
func (f *Fx) Authentik(baseURL, apiKey string) {
	cfg := api.NewConfiguration()
	cfg.Host = baseURL
	cfg.Scheme = "https"
	cfg.DefaultHeader = map[string]string{
		"Authorization": "Bearer " + apiKey,
	}
	client := api.NewAPIClient(cfg)
	f.Auth = client
}

// TestConnection fetches and prints the current user info to test the connection.
func (f *Fx) WhoAmIAuthentik() error {
	user, _, err := f.Auth.CoreApi.CoreUsersMeRetrieve(context.TODO()).Execute()
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
