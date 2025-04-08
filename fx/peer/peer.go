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

func (p *Peers) HelloUniverse() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	peers := len(p.Addresses) // Use the actual length of the slice
	fmt.Printf("%d peers to process\n", peers)

	for _, address := range p.Addresses {
		peer := p.Map[address] // No need to check existence; LoadUnprocessedAddresses ensures validity

		// Populate missing fields for the peer
		p.GetENS(peer, peer.Address)
		p.GetLoopringENS(peer, peer.Address)
		p.GetLoopringID(peer, peer.Address)

		// Save the updated peer to the database
		if err := p.SavePeer(peer); err != nil {
			fmt.Printf("Error saving peer %s: %v\n", peer.Address, err)
		}

		// Update progress and print details
		fmt.Printf("%d | %s %s %s %s\n", peers, peer.ENS, peer.LoopringENS, peer.LoopringID, peer.Address)
		peers--
	}

	fmt.Println("Hello Universe")
}
