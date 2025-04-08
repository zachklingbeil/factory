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

	if err := peers.LoadPeers(); err != nil {
		fmt.Printf("Error loading peers: %v\n", err)
	}

	return peers
}

func (p *Peers) LoadPeers() error {
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
	// Initialize the map
	var addresses []string // Temporary slice for unprocessed addresses

	for rows.Next() {
		var peer Peer
		if err := rows.Scan(&peer.Address, &peer.ENS, &peer.LoopringENS, &peer.LoopringID); err != nil {
			return fmt.Errorf("failed to scan peer row: %w", err)
		}
		p.Map[peer.Address] = &peer

		// Check if the peer is unprocessed
		if peer.ENS == "" ||
			peer.LoopringENS == "" || peer.LoopringENS == "!" ||
			peer.LoopringID == "" || peer.LoopringID == "!" {
			addresses = append(addresses, peer.Address)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over peer rows: %w", err)
	}

	p.Addresses = addresses // Assign unprocessed addresses
	fmt.Printf("%d peers, %d peers to process\n", len(p.Map), len(p.Addresses))
	return nil
}

func (p *Peers) SavePeers(peers []*Peer) error {
	query := `
    INSERT INTO peers (address, ens, loopring_ens, loopring_id)
    VALUES %s
    ON CONFLICT (address) DO UPDATE SET
        ens = EXCLUDED.ens,
        loopring_ens = EXCLUDED.loopring_ens,
        loopring_id = EXCLUDED.loopring_id
    `

	// Build the query with placeholders
	values := []any{}
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

	peers := len(p.Addresses)
	fmt.Printf("%d peers to process\n", peers)

	batchSize := 1000
	var batch []*Peer

	for _, address := range p.Addresses {
		peer := p.Map[address]

		p.GetENS(peer, peer.Address)
		p.GetLoopringENS(peer, peer.Address)
		p.GetLoopringID(peer, peer.Address)

		batch = append(batch, peer)

		if len(batch) >= batchSize {
			if err := p.SavePeers(batch); err != nil {
				fmt.Printf("Error saving batch: %v\n", err)
			} else {
				fmt.Printf("Successfully updated a batch of %d peers.\n", len(batch))
			}
			batch = batch[:0]
		}

		fmt.Printf("%d %s %s %s\n", peers, peer.ENS, peer.LoopringENS, peer.LoopringID)
		peers--
	}

	if len(batch) > 0 {
		if err := p.SavePeers(batch); err != nil {
			fmt.Printf("Error saving final batch: %v\n", err)
		}
	}
	fmt.Println("Hello Universe")
}
