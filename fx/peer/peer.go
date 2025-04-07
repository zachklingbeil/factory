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
	Addresses      []string // Add a slice to store addresses
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

	if err := peers.LoadAddresses(); err != nil {
		fmt.Printf("Error loading addresses: %v\n", err)
	}

	return peers
}

func (p *Peers) InitPeers() error {
	// Step 1: Fetch all addresses from the old table
	var addresses []string
	err := p.Db.ColumnToSlice("peers2", "address", &addresses) // Replace "peers2" with the old table name
	if err != nil {
		return fmt.Errorf("failed to load addresses from old table: %w", err)
	}

	// Step 2: Create the new peers table
	query := `
    CREATE TABLE IF NOT EXISTS peers (
        address TEXT PRIMARY KEY,
        ens TEXT,
        loopring_ens TEXT,
        loopring_id TEXT
    )`
	if _, err := p.Db.Exec(query); err != nil {
		return fmt.Errorf("failed to create new peers table: %w", err)
	}

	// Step 3: Insert all peers with only the address field in a single transaction
	tx, err := p.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := tx.Prepare(`
    INSERT INTO peers (address, ens, loopring_ens, loopring_id)
    VALUES ($1, NULL, NULL, NULL)
    ON CONFLICT (address) DO NOTHING
    `)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, address := range addresses {
		if _, err := stmt.Exec(address); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert peer with address %s: %w", address, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("Successfully initialized %d peers with addresses only.\n", len(addresses))
	return nil
}

func (p *Peers) CreateTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS peers (
        address TEXT PRIMARY KEY,
        ens TEXT,
        loopring_ens TEXT,
        loopring_id TEXT
    )`
	_, err := p.Db.Exec(query)
	return err
}

// LoadAddresses fetches all addresses from the peers table and stores them in the Peers struct
func (p *Peers) LoadAddresses() error {
	var addresses []string
	err := p.Db.ColumnToSlice("peers2", "address", &addresses)
	if err != nil {
		return fmt.Errorf("failed to load addresses: %w", err)
	}

	p.Mu.Lock()
	defer p.Mu.Unlock()
	p.Addresses = addresses
	return nil
}

func (p *Peers) HelloUniverse() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	for _, address := range p.Addresses {
		value := p.Format(address)

		peer, exists := p.Map[value]
		if !exists {
			peer = &Peer{Address: value}
			p.Map[value] = peer
		}

		// Save the peer with just the address
		if err := p.SavePeer(peer); err != nil {
			fmt.Printf("Error saving peer %s: %v\n", peer.Address, err)
		}
	}
	fmt.Printf("Finished creating peers with addresses.\n")

	// Step 2: Process peers for ENS, LoopringENS, and LoopringID
	fmt.Printf("Processing peers for ENS, LoopringENS, and LoopringID...\n")
	work := len(p.Addresses) // Track the number of addresses
	fmt.Printf("%d peers to process\n", work)

	for _, address := range p.Addresses {
		value := p.Format(address)

		peer, exists := p.Map[value]
		if !exists {
			fmt.Printf("Peer not found in map for address: %s\n", address)
			continue
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

		// Save the updated peer
		if err := p.SavePeer(peer); err != nil {
			fmt.Printf("Error saving peer %s: %v\n", peer.Address, err)
		}

		work--
		fmt.Printf("%d peers remaining\n", work)
	}
	fmt.Printf("Finished processing peers.\n")
}

// func (p *Peers) LoadUnprocessedAddresses() error {
// 	var addresses []string
// 	query := `
//         SELECT address FROM peers
//         WHERE ens IS NULL OR loopring_ens IS NULL OR loopring_id IS NULL
//     `
// 	err := p.Db.ColumnToSlice(query, "address", &addresses)
// 	if err != nil {
// 		return fmt.Errorf("failed to load unprocessed addresses: %w", err)
// 	}

// 	p.Mu.Lock()
// 	defer p.Mu.Unlock()
// 	p.Addresses = addresses
// 	return nil
// }

// func (p *Peers) HelloUniverse() {
// 	p.Mu.Lock()
// 	defer p.Mu.Unlock()

// 	// Iterate over all addresses in p.Addresses
// 	work := len(p.Addresses) // Track the number of addresses
// 	fmt.Printf("%d\n", work)

// 	for _, address := range p.Addresses {
// 		value := p.Format(address)

// 		peer, exists := p.Map[value]
// 		if !exists {
// 			peer = &Peer{Address: value}
// 			p.Map[value] = peer
// 		}

// 		// Populate missing fields for the peer
// 		if peer.ENS == "" {
// 			p.GetENS(peer, peer.Address)
// 			if peer.ENS == "" {
// 				peer.ENS = "."
// 			}
// 		}
// 		if peer.LoopringENS == "" {
// 			p.GetLoopringENS(peer, peer.Address)
// 			if peer.LoopringENS == "" {
// 				peer.LoopringENS = "."
// 			}
// 		}
// 		if peer.LoopringID == "" {
// 			p.GetLoopringID(peer, peer.Address)
// 			if peer.LoopringID == "" {
// 				peer.LoopringID = "."
// 			}
// 		}
// 		if err := p.SavePeer(peer); err != nil {
// 			fmt.Printf("Error saving peer %s: %v\n", peer.Address, err)
// 		}

// 		work--
// 		fmt.Printf("%d\n", work)
// 	}
// }

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
