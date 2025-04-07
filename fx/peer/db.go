package peer

import (
	"fmt"
	"time"
)

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

func (p *Peers) LoadMap() error {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	rows, err := p.Db.Query(`SELECT address, ens, loopring_ens, loopring_id FROM peers`)
	if err != nil {
		return fmt.Errorf("failed to load peers from database: %w", err)
	}
	defer rows.Close()

	incompleteAddresses := make([]string, 0)

	for rows.Next() {
		var peer Peer
		if err := rows.Scan(&peer.Address, &peer.ENS, &peer.LoopringENS, &peer.LoopringID); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		// Add the peer to the map
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

		if peer.ENS == "" || peer.LoopringENS == "" || peer.LoopringID == "" {
			incompleteAddresses = append(incompleteAddresses, peer.Address)
		}
	}

	if len(incompleteAddresses) > 0 {
		fmt.Printf("Found %d peers with incomplete fields. Initializing...\n", len(incompleteAddresses))
		go p.HelloPeers(incompleteAddresses) // Run HelloPeers in a goroutine for non-blocking behavior
	}

	return nil
}

func (p *Peers) SaveMap() error {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	// Use a transaction for batch inserts/updates to improve performance
	tx, err := p.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `
    INSERT INTO peers (address, ens, loopring_ens, loopring_id)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (address) DO UPDATE
    SET ens = EXCLUDED.ens,
    loopring_ens = EXCLUDED.loopring_ens,
    loopring_id = EXCLUDED.loopring_id
    `

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, peer := range p.Map {
		if _, err := stmt.Exec(peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to save peer with address %s: %w", peer.Address, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
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
