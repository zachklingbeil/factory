package peer

import (
	"encoding/json"
	"fmt"
	"os"
)

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
	var addresses []string

	for rows.Next() {
		var peer Peer
		if err := rows.Scan(&peer.Address, &peer.ENS, &peer.LoopringENS, &peer.LoopringID); err != nil {
			return fmt.Errorf("failed to scan peer row: %w", err)
		}
		p.Map[peer.Address] = &peer

		if peer.ENS == "" || peer.ENS == "!" ||
			peer.LoopringENS == "" || peer.LoopringENS == "!" ||
			peer.LoopringID == "" || peer.LoopringID == "!" {
			addresses = append(addresses, peer.Address)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over peer rows: %w", err)
	}

	p.Addresses = addresses
	fmt.Printf("%d peers\n", len(p.Map))
	return nil
}
func (p *Peers) OutputPeersAsJSON() error {
	query := `
        SELECT address, ens, loopring_ens, loopring_id FROM peers
    `
	rows, err := p.Db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query peers table: %w", err)
	}
	defer rows.Close()

	// Create a slice to hold the Peer objects
	var peers []Peer

	for rows.Next() {
		var peer Peer
		if err := rows.Scan(&peer.Address, &peer.ENS, &peer.LoopringENS, &peer.LoopringID); err != nil {
			return fmt.Errorf("failed to scan peer row: %w", err)
		}
		peers = append(peers, peer)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over peer rows: %w", err)
	}

	// Convert the slice of Peer objects to JSON
	jsonData, err := json.MarshalIndent(peers, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal peers to JSON: %w", err)
	}

	// Write JSON to a file or print to console
	file, err := os.Create("peers.json")
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	fmt.Println("Peers table exported to peers.json")
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

	values := []any{}
	placeholders := ""
	for i, peer := range peers {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		values = append(values, peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
	}

	query = fmt.Sprintf(query, placeholders)

	_, err := p.Db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to save peers batch: %w", err)
	}
	return nil
}
