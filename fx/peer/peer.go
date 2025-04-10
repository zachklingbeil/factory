package peer

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zachklingbeil/factory/fx"
)

//go:embed peers.json
var peerJSON []byte

type Peers struct {
	Json      *fx.JSON
	Eth       *ethclient.Client
	Db        *fx.Database
	Map       map[string]*Peer
	Addresses []string
	Mu        sync.RWMutex
	PeerChan  chan string
}

type Peer struct {
	Address     string `json:"address"`
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  int64  `json:"loopringId"`
}

func NewPeers(json *fx.JSON, eth *ethclient.Client, db *fx.Database) *Peers {
	peers := &Peers{
		Json:      json,
		Eth:       eth,
		Map:       make(map[string]*Peer),
		Addresses: nil,
		Db:        db,
	}

	// Initialize the Map with data from the embedded peer.json file
	if err := peers.InitializeFromEmbeddedJSON(); err != nil {
		fmt.Printf("Error initializing peers from embedded JSON: %v\n", err)
	}

	// Save the initialized peers to the database
	if err := peers.SavePeersToDB(); err != nil {
		fmt.Printf("Error saving peers to database: %v\n", err)
	}

	return peers
}

// func NewPeers(json *fx.JSON, eth *ethclient.Client, db *fx.Database) *Peers {
// 	peers := &Peers{
// 		Json:      json,
// 		Eth:       eth,
// 		Map:       make(map[string]*Peer),
// 		Addresses: nil,
// 		Db:        db,
// 	}

// 	if err := peers.LoadPeers(); err != nil {
// 		fmt.Printf("Error loading peers: %v\n", err)
// 	}

// 	peers.PeerChan = make(chan string, len(peers.Addresses))

// 	for _, address := range peers.Addresses {
// 		peers.PeerChan <- address
// 	}

// 	return peers
// }

func (p *Peers) HelloUniverse() {
	batchSize := 1000
	var batch []*Peer

	p.Mu.RLock()
	peers := len(p.Addresses)
	p.Mu.RUnlock()
	fmt.Printf("%d peers to process\n", peers)

	for {
		if len(batch) > 0 {
			if err := p.SavePeers(batch); err != nil {
				fmt.Printf("Error saving final batch: %v\n", err)
			}
			batch = batch[:0]
		}

		if peers == 0 && len(batch) == 0 {
			break
		}

		address := <-p.PeerChan

		p.Mu.Lock()
		peer := p.Map[address]
		p.Mu.Unlock()

		p.GetENS(peer, peer.Address)
		p.GetLoopringENS(peer, peer.Address)
		p.GetLoopringID(peer, peer.Address)

		batch = append(batch, peer)

		fmt.Printf("%d %s %s %d\n", peers, peer.ENS, peer.LoopringENS, peer.LoopringID)
		peers--

		if len(batch) >= batchSize {
			if err := p.SavePeers(batch); err != nil {
				fmt.Printf("Error saving batch: %v\n", err)
			}
			batch = batch[:0]
		}
	}
	fmt.Println("Hello Universe")
}

func (p *Peers) NewBlock(addresses []string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	peers := len(addresses)
	fmt.Printf("%d peers from new block\n", peers)

	for _, address := range addresses {
		if _, exists := p.Map[address]; !exists {
			p.Map[address] = &Peer{Address: address}
			p.Addresses = append(p.Addresses, address)
		}
		p.PeerChan <- address
	}
}

// InitializeFromEmbeddedJSON initializes the Map with data from the embedded peer.json file
func (p *Peers) InitializeFromEmbeddedJSON() error {
	var peerList []Peer

	// Unmarshal the embedded JSON into a slice of Peer structs
	if err := json.Unmarshal(peerJSON, &peerList); err != nil {
		return fmt.Errorf("failed to unmarshal embedded JSON: %w", err)
	}

	// Populate the Map and Addresses slice
	for _, peer := range peerList {
		p.Map[peer.Address] = &peer
		p.Addresses = append(p.Addresses, peer.Address)
	}

	fmt.Println("Peers initialized from embedded JSON successfully.")
	return nil
}

func (p *Peers) SavePeersToDB() error {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	// Convert the Map to a slice of Peer structs
	var peerData []Peer
	for _, peer := range p.Map {
		peerData = append(peerData, *peer)
	}

	// Insert the data into the database with the key "peer"
	// Convert the slice of Peer structs to JSON-compatible format
	var jsonData []map[string]any
	for _, peer := range peerData {
		jsonData = append(jsonData, map[string]any{
			"address":     peer.Address,
			"ens":         peer.ENS,
			"loopringEns": peer.LoopringENS,
			"loopringId":  peer.LoopringID,
		})
	}

	if err := p.Db.Insert("peer", jsonData); err != nil {
		return fmt.Errorf("failed to save peers to database: %w", err)
	}

	fmt.Println("Peers saved to database successfully.")
	return nil
}
