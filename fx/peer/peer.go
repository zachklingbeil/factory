package peer

import (
	"fmt"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zachklingbeil/factory/fx"
)

type Peers struct {
	Json           *fx.JSON
	Eth            *ethclient.Client
	Db             *fx.Database
	LoopringApiKey string
	Map            map[string]*Peer
	Mu             sync.RWMutex
}

type Peer struct {
	Address     string
	ENS         string
	LoopringENS string
	LoopringID  string
}

func NewPeers(json *fx.JSON, eth *ethclient.Client, db *fx.Database) *Peers {
	peers := &Peers{
		Json:           json,
		Eth:            eth,
		LoopringApiKey: os.Getenv("LOOPRING_API_KEY"),
		Map:            make(map[string]*Peer),
		Db:             db,
	}

	if err := peers.CreateTable(); err != nil {
		fmt.Printf("Error ensuring peers table exists: %v\n", err)
	}

	if err := peers.LoadMap(); err != nil {
		fmt.Printf("Error loading map from database: %v\n", err)
	} else {
		fmt.Println("Map loaded successfully from the database.")
	}

	go peers.Checkpoint(60)
	return peers
}
