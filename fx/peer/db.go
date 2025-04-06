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

func (p *Peers) Checkpoint(intervalSeconds int) {
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	for range ticker.C {
		if err := p.SaveMap(); err != nil {
			fmt.Printf("Failed to save map to database: %v\n", err)
		} else {
			fmt.Println("Map saved to database successfully.")
		}
	}
}

func (p *Peers) SaveMap() error {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	query := `
	INSERT INTO peers (address, ens, loopring_ens, loopring_id)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (address) DO UPDATE
	SET ens = EXCLUDED.ens,
	loopring_ens = EXCLUDED.loopring_ens,
	loopring_id = EXCLUDED.loopring_id
	`
	for _, peer := range p.Map {
		if _, err := p.Db.Exec(query, peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID); err != nil {
			return fmt.Errorf("failed to save peer with address %s: %w", peer.Address, err)
		}
	}
	return nil
}

func (p *Peers) LoadMap() error {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	rows, err := p.Db.Query(`SELECT address, ens, loopring_ens, loopring_id FROM peers`)
	if err != nil {
		return fmt.Errorf("failed to load peers from database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var peer Peer
		if err := rows.Scan(&peer.Address, &peer.ENS, &peer.LoopringENS, &peer.LoopringID); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
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
