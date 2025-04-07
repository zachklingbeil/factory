package peer

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

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

		// Start Checkpoint in a separate goroutine
		go peers.Checkpoint(30) // Periodic saves start immediately

		// Start HelloUniverse in a separate goroutine
		go peers.HelloUniverse()
	}

	return peers
}

func (p *Peers) Checkpoint(intervalSeconds int) {
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop() // Ensure the ticker is stopped when the function exits

	for range ticker.C {
		start := time.Now()
		if err := p.SaveMap(); err != nil {
			fmt.Printf("Failed to save map to database: %v\n", err)
		} else {
			fmt.Printf("Map saved to database successfully in %v.\n", time.Since(start))
		}
	}
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

func (p *Peers) CreateTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS json (
        timestamp TIMESTAMP PRIMARY KEY,
    	  data JSONB
    )`
	_, err := p.Db.Exec(query)
	return err
}

func (p *Peers) SaveMap() error {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	// Serialize the map into JSON
	peersSlice := make([]Peer, 0, len(p.Map))
	for _, peer := range p.Map {
		peersSlice = append(peersSlice, *peer)
	}

	jsonData, err := json.Marshal(peersSlice)
	if err != nil {
		return fmt.Errorf("failed to serialize peers map: %w", err)
	}

	// Insert the JSON data into the database
	query := `
    INSERT INTO json (timestamp, data)
    VALUES ($1, $2)
    ON CONFLICT (timestamp) DO UPDATE
    SET data = EXCLUDED.data
    `
	_, err = p.Db.Exec(query, time.Now(), jsonData)
	if err != nil {
		return fmt.Errorf("failed to save peers to database: %w", err)
	}

	fmt.Println("Checkpoint: Map saved to database.")
	return nil
}
