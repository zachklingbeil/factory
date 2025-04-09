package fx

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Establish geth.ipc connection (http, websocket disabled)
func Node(ctx context.Context) (*rpc.Client, *ethclient.Client, error) {
	rpc, err := rpc.DialIPC(ctx, "/ethereum/.ethereum/geth.ipc")
	if err != nil {
		log.Printf("Failed to connect to the Ethereum client: %v", err)
		return nil, nil, err
	}
	// log.Println("Successfully connected to the Ethereum client.")
	eth := ethclient.NewClient(rpc)
	return rpc, eth, nil
}
