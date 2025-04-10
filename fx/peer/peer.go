package peer

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zachklingbeil/factory/fx"
)

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
	Address       string
	ENS           string
	LoopringENS   string
	LoopringID    string
	LoopringIDINT int64
}

func NewPeers(json *fx.JSON, eth *ethclient.Client, db *fx.Database) *Peers {
	peers := &Peers{
		Json:      json,
		Eth:       eth,
		Map:       make(map[string]*Peer),
		Addresses: nil,
		Db:        db,
	}

	if err := peers.LoadPeers(); err != nil {
		fmt.Printf("Error loading peers: %v\n", err)
	}
	if err := peers.UpdateLoopringIDInt(); err != nil {
		fmt.Printf("Error updating LoopringIDInt: %v\n", err)
	}
	// Save the peers to a JSON file
	if err := peers.SavePeersToJSON("peers.json"); err != nil {
		fmt.Printf("Error saving peers to JSON: %v\n", err)
	}

	peers.PeerChan = make(chan string, len(peers.Addresses))

	for _, address := range peers.Addresses {
		peers.PeerChan <- address
	}

	return peers
}

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

		fmt.Printf("%d %s %s %s\n", peers, peer.ENS, peer.LoopringENS, peer.LoopringID)
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

// SavePeersToJSON saves the Peers.Map to a JSON file as a slice of Peer objects and logs the count.
func (p *Peers) SavePeersToJSON(filename string) error {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	// Create a slice to hold the Peer objects
	var peersSlice []*Peer
	for _, peer := range p.Map {
		peersSlice = append(peersSlice, peer)
	}

	// Log the number of peers
	fmt.Printf("Number of peers to save: %d\n", len(peersSlice))

	// Marshal the slice to JSON
	data, err := json.MarshalIndent(peersSlice, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal peers to JSON: %w", err)
	}

	// Write the JSON data to a file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write JSON data to file: %w", err)
	}

	fmt.Printf("Peers saved to JSON file: %s\n", filename)
	return nil
}

// UpdateLoopringIDInt updates the LoopringIDInt column in the database by converting LoopringID strings to integers.
func (p *Peers) UpdateLoopringIDInt() error {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	for _, peer := range p.Map {
		// Convert LoopringID string to int64
		loopringIDInt, err := strconv.ParseInt(peer.LoopringID, 10, 64)
		if err != nil {
			fmt.Printf("Failed to convert LoopringID '%s' to int: %v\n", peer.LoopringID, err)
			continue
		}

		// Update the database with the new value
		query := "UPDATE peers SET loopring_id_int = ? WHERE address = ?"
		_, err = p.Db.Exec(query, loopringIDInt, peer.Address)
		if err != nil {
			fmt.Printf("Failed to update LoopringIDInt for address '%s': %v\n", peer.Address, err)
			continue
		}

		// Update the in-memory map
		peer.LoopringIDINT = loopringIDInt
	}

	fmt.Println("LoopringIDInt column updated successfully.")
	return nil
}
