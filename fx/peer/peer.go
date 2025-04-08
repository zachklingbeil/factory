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
	Addresses      []string
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
		Addresses:      nil,
		Db:             db,
	}

	// Load the entire map first
	if err := peers.LoadMap(); err != nil {
		fmt.Printf("Error loading map: %v\n", err)
	}

	// Then load unprocessed addresses
	if err := peers.LoadUnprocessedAddresses(); err != nil {
		fmt.Printf("Error loading unprocessed addresses: %v\n", err)
	}

	return peers
}

func (p *Peers) LoadMap() error {
	query := `
        SELECT address, ens, loopring_ens, loopring_id FROM peers
    `
	rows, err := p.Db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to load peers from database: %w", err)
	}
	defer rows.Close()
	p.Mu.Lock()
	defer p.Mu.Unlock()

	for rows.Next() {
		var peer Peer
		if err := rows.Scan(&peer.Address, &peer.ENS, &peer.LoopringENS, &peer.LoopringID); err != nil {
			return fmt.Errorf("failed to scan peer row: %w", err)
		}
		p.Map[peer.Address] = &peer
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over peer rows: %w", err)
	}

	fmt.Printf("Loaded %d peers into the map.\n", len(p.Map))
	return nil
}

func (p *Peers) LoadUnprocessedAddresses() error {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	var addresses []string
	for address, peer := range p.Map {
		if peer.ENS == "" || peer.ENS == "!" ||
			peer.LoopringENS == "" || peer.LoopringENS == "!" ||
			peer.LoopringID == "" || peer.LoopringID == "!" {
			addresses = append(addresses, address)
		}
	}

	p.Addresses = addresses // Set the slice with unprocessed addresses
	fmt.Printf("Loaded %d unprocessed addresses.\n", len(p.Addresses))
	return nil
}

func (p *Peers) SavePeer(peer *Peer) error {
	query := `
    INSERT INTO peers (address, ens, loopring_ens, loopring_id)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (address) DO UPDATE SET
        ens = EXCLUDED.ens,
        loopring_ens = EXCLUDED.loopring_ens,
        loopring_id = EXCLUDED.loopring_id
    `
	_, err := p.Db.Exec(query, peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
	if err != nil {
		return fmt.Errorf("failed to save peer %s: %w", peer.Address, err)
	}
	return nil
}
func (p *Peers) SavePeersBatch(peers []*Peer) error {
	query := `
    INSERT INTO peers (address, ens, loopring_ens, loopring_id)
    VALUES %s
    ON CONFLICT (address) DO UPDATE SET
        ens = EXCLUDED.ens,
        loopring_ens = EXCLUDED.loopring_ens,
        loopring_id = EXCLUDED.loopring_id
    `

	// Build the query with placeholders
	values := []interface{}{}
	placeholders := ""
	for i, peer := range peers {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		values = append(values, peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
	}

	// Format the query with the placeholders
	query = fmt.Sprintf(query, placeholders)

	// Execute the batch insert
	_, err := p.Db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to save peers batch: %w", err)
	}
	return nil
}
func (p *Peers) HelloUniverse() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	peers := len(p.Addresses) // Use the actual length of the slice
	fmt.Printf("%d peers to process\n", peers)

	batchSize := 1000 // Define the batch size
	var batch []*Peer

	for _, address := range p.Addresses {
		peer := p.Map[address] // No need to check existence; LoadUnprocessedAddresses ensures validity

		// Populate missing fields for the peer
		p.GetENS(peer, peer.Address)
		p.GetLoopringENS(peer, peer.Address)
		p.GetLoopringID(peer, peer.Address)

		// Add the peer to the batch
		batch = append(batch, peer)

		// If the batch size is reached, execute the batch insert
		if len(batch) >= batchSize {
			if err := p.SavePeersBatch(batch); err != nil {
				fmt.Printf("Error saving batch: %v\n", err)
			}
			batch = batch[:0] // Clear the batch
		}

		// Update progress and print details
		fmt.Printf("%d | %s %s %s\n", peers, peer.ENS, peer.LoopringENS, peer.LoopringID)
		peers--
	}

	// Insert any remaining peers in the batch
	if len(batch) > 0 {
		if err := p.SavePeersBatch(batch); err != nil {
			fmt.Printf("Error saving final batch: %v\n", err)
		}
	}

	fmt.Println("Hello Universe")
}
