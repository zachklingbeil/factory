package peer

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	_ "github.com/lib/pq"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zachklingbeil/factory/fx"
)

//go:embed peers.json
var embeddedPeersJSON string

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

	// Create the peers table if it doesn't exist
	if err := peers.CreatePeersTable(); err != nil {
		fmt.Printf("Error creating peers table: %v\n", err)
	}

	if err := peers.LoadAndSavePeersFromJSON(); err != nil {
		fmt.Printf("Error loading peers: %v\n", err)
	}

	peers.PeerChan = make(chan string, len(peers.Addresses))

	for _, address := range peers.Addresses {
		peers.PeerChan <- address
	}

	return peers
}

func (p *Peers) LoadPeersFromJSON() error {
	// Use the embedded JSON data
	data := []byte(embeddedPeersJSON)

	// Unmarshal the JSON data into a slice of Peer objects
	var peers []Peer
	if err := json.Unmarshal(data, &peers); err != nil {
		return fmt.Errorf("failed to unmarshal embedded JSON data: %w", err)
	}

	// Populate the Map and Addresses fields
	p.Mu.Lock()
	defer p.Mu.Unlock()

	for _, peer := range peers {
		// Add the peer to the Map
		p.Map[peer.Address] = &peer

		// Add the peer's address to the Addresses slice if fields are invalid
		if peer.ENS == "." || peer.LoopringENS == "." || peer.LoopringID == -1 {
			p.Addresses = append(p.Addresses, peer.Address)
		}
	}

	fmt.Printf("%d peers loaded from embedded JSON\n", len(peers))
	return nil
}

func (p *Peers) LoadAndSavePeersFromJSON() error {
	// Load peers from the embedded JSON
	if err := p.LoadPeersFromJSON(); err != nil {
		return fmt.Errorf("failed to load peers from embedded JSON: %w", err)
	}

	// Save all peers to the database
	if err := p.SavePeers(); err != nil {
		return fmt.Errorf("failed to save peers to the database: %w", err)
	}

	return nil
}

func (p *Peers) CreatePeersTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS peers (
        address TEXT PRIMARY KEY,
        ens TEXT NOT NULL,
        loopringEns TEXT NOT NULL,
        loopringId BIGINT NOT NULL
    );
    `
	_, err := p.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create peers table: %w", err)
	}
	fmt.Println("Peers table created or already exists.")
	return nil
}
func (p *Peers) SavePeers() error {
	const batchSize = 1000 // Number of peers per batch
	queryTemplate := `
    INSERT INTO peers (address, ens, loopring_ens, loopring_id)
    VALUES %s
    ON CONFLICT (address) DO UPDATE SET
        ens = EXCLUDED.ens,
        loopring_ens = EXCLUDED.loopring_ens,
        loopring_id = EXCLUDED.loopring_id
    `

	p.Mu.RLock()
	defer p.Mu.RUnlock()

	peers := make([]*Peer, 0, len(p.Map))
	for _, peer := range p.Map {
		peers = append(peers, peer)
	}

	for i := 0; i < len(peers); i += batchSize {
		end := i + batchSize
		if end > len(peers) {
			end = len(peers)
		}

		batch := peers[i:end]
		values := []any{}
		placeholders := ""

		for j, peer := range batch {
			if j > 0 {
				placeholders += ", "
			}
			placeholders += fmt.Sprintf("($%d, $%d, $%d, $%d)", j*4+1, j*4+2, j*4+3, j*4+4)
			values = append(values, peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
		}

		query := fmt.Sprintf(queryTemplate, placeholders)

		if _, err := p.Db.Exec(query, values...); err != nil {
			return fmt.Errorf("failed to save peers batch: %w", err)
		}
	}

	fmt.Printf("%d peers saved to the database\n", len(peers))
	return nil
}

// func (p *Peers) SavePeers() error {
// 	query := `
//     INSERT INTO peers (address, ens, loopringEns, loopringId)
//     VALUES %s
//     ON CONFLICT (address) DO UPDATE SET
//         ens = EXCLUDED.ens,
//         loopringEns = EXCLUDED.loopringEns,
//         loopringId = EXCLUDED.loopringId
//     `

// 	values := []any{}
// 	placeholders := ""
// 	i := 0

// 	p.Mu.RLock()
// 	defer p.Mu.RUnlock()

// 	for _, peer := range p.Map {
// 		if i > 0 {
// 			placeholders += ", "
// 		}
// 		placeholders += fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
// 		values = append(values, peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
// 		i++
// 	}

// 	query = fmt.Sprintf(query, placeholders)

// 	_, err := p.Db.Exec(query, values...)
// 	if err != nil {
// 		return fmt.Errorf("failed to save peers batch: %w", err)
// 	}
// 	fmt.Printf("%d peers saved to the database\n", len(p.Map))
// 	return nil
// }

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

// func (p *Peers) HelloUniverse() {
// 	batchSize := 1000
// 	var batch []*Peer

// 	p.Mu.RLock()
// 	peers := len(p.Addresses)
// 	p.Mu.RUnlock()
// 	fmt.Printf("%d peers to process\n", peers)

// 	for {
// 		if len(batch) > 0 {
// 			if err := p.SavePeers(batch); err != nil {
// 				fmt.Printf("Error saving final batch: %v\n", err)
// 			}
// 			batch = batch[:0]
// 		}

// 		if peers == 0 && len(batch) == 0 {
// 			break
// 		}

// 		address := <-p.PeerChan

// 		p.Mu.Lock()
// 		peer := p.Map[address]
// 		p.Mu.Unlock()

// 		p.GetENS(peer, peer.Address)
// 		p.GetLoopringENS(peer, peer.Address)
// 		p.GetLoopringID(peer, peer.Address)

// 		batch = append(batch, peer)

// 		fmt.Printf("%d %s %s %d\n", peers, peer.ENS, peer.LoopringENS, peer.LoopringID)
// 		peers--

// 		if len(batch) >= batchSize {
// 			if err := p.SavePeers(batch); err != nil {
// 				fmt.Printf("Error saving batch: %v\n", err)
// 			}
// 			batch = batch[:0]
// 		}
// 	}
// 	fmt.Println("Hello Universe")
// }

// func (p *Peers) SavePeers(peers []*Peer) error {
// 	query := `
//     INSERT INTO peers (address, ens, loopringEns, loopringId)
//     VALUES %s
//     ON CONFLICT (address) DO UPDATE SET
//         ens = EXCLUDED.ens,
//         loopringEns = EXCLUDED.loopringEns,
//         loopringId = EXCLUDED.loopringId
//     `

// 	values := []any{}
// 	placeholders := ""
// 	for i, peer := range peers {
// 		if i > 0 {
// 			placeholders += ", "
// 		}
// 		placeholders += fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
// 		values = append(values, peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
// 	}

// 	query = fmt.Sprintf(query, placeholders)

//		_, err := p.Db.Exec(query, values...)
//		if err != nil {
//			return fmt.Errorf("failed to save peers batch: %w", err)
//		}
//		return nil
//	}
