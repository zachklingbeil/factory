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
		peers.HelloUniverse() // Directly handle incomplete addresses
	}

	go peers.Checkpoint(60)
	return peers
}

func (p *Peers) HelloUniverse() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	// Collect all addresses with incomplete fields
	incompleteAddresses := []string{}
	for _, peer := range p.Map {
		if peer.ENS == "" || peer.LoopringENS == "" || peer.LoopringID == "" {
			incompleteAddresses = append(incompleteAddresses, peer.Address)
		}
	}

	work := len(incompleteAddresses) // Track the number of incomplete addresses
	fmt.Printf("%d incomplete peers\n", work)

	// Process each incomplete address
	for _, address := range incompleteAddresses {
		value := p.Format(address)

		peer, exists := p.Map[value]
		if !exists {
			peer = &Peer{Address: value}
			p.Map[value] = peer
		}

		// Populate missing fields for the peer
		if peer.ENS == "" {
			p.GetENS(peer, peer.Address)
			if peer.ENS == "" {
				peer.ENS = "."
			}
		}
		if peer.LoopringENS == "" {
			p.GetLoopringENS(peer, peer.Address)
			if peer.LoopringENS == "" {
				peer.LoopringENS = "."
			}
		}
		if peer.LoopringID == "" {
			p.GetLoopringID(peer, peer.Address)
			if peer.LoopringID == "" {
				peer.LoopringID = "."
			}
		}

		work--
		fmt.Printf("%d\n", work)
	}
}

func (p *Peers) LoadMap() error {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	// Load all peers from the database into a slice
	var peers []Peer
	if err := p.Db.DiskToMem("peers", &peers); err != nil {
		return fmt.Errorf("failed to load peers from database: %w", err)
	}

	// Populate the map with the loaded peers
	for _, peer := range peers {
		p.Map[peer.Address] = &peer
		if peer.ENS != "" {
			p.Map[peer.ENS] = &peer
		}
		if peer.LoopringENS != "" {
			p.Map[peer.LoopringENS] = &peer
		}
		if peer.LoopringID != "" {
			p.Map[peer.LoopringID] = &peer
		}
	}

	return nil
}
