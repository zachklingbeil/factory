package peer

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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
	fmt.Println("Starting OutputPeersAsJSON...")

	query := `
        SELECT address, ens, loopring_ens, loopring_id FROM peers
    `
	fmt.Println("Executing query to fetch peers...")
	rows, err := p.Db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query peers table: %w", err)
	}
	defer rows.Close()
	fmt.Println("Query executed successfully.")

	// Create a slice to hold the JSON Peer objects
	var jsonPeers []struct {
		Address     string `json:"address"`
		ENS         string `json:"ens"`
		LoopringENS string `json:"loopring_ens"`
		LoopringID  int64  `json:"loopring_id"`
	}

	fmt.Println("Processing rows...")
	for rows.Next() {
		var peer Peer
		var loopringIDStr string
		if err := rows.Scan(&peer.Address, &peer.ENS, &peer.LoopringENS, &loopringIDStr); err != nil {
			return fmt.Errorf("failed to scan peer row: %w", err)
		}

		// Convert LoopringID from string to int64
		loopringID, err := strconv.ParseInt(loopringIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to convert LoopringID to int64: %w", err)
		}

		// Append to the JSON slice
		jsonPeers = append(jsonPeers, struct {
			Address     string `json:"address"`
			ENS         string `json:"ens"`
			LoopringENS string `json:"loopring_ens"`
			LoopringID  int64  `json:"loopring_id"`
		}{
			Address:     peer.Address,
			ENS:         peer.ENS,
			LoopringENS: peer.LoopringENS,
			LoopringID:  loopringID,
		})
	}
	fmt.Println("Rows processed successfully.")

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over peer rows: %w", err)
	}

	// Convert the slice of JSON Peer objects to JSON
	fmt.Println("Marshalling data to JSON...")
	jsonData, err := json.MarshalIndent(jsonPeers, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal peers to JSON: %w", err)
	}
	fmt.Println("Data marshalled to JSON successfully.")

	// Write JSON to a file or print to console
	fmt.Println("Creating JSON file...")
	file, err := os.Create("peers.json")
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()
	fmt.Println("JSON file created successfully.")

	fmt.Println("Writing data to JSON file...")
	if _, err := file.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}
	fmt.Println("Data written to JSON file successfully.")

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
