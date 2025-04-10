package peer

import (
	"fmt"
)

func (p *Peers) LoadPeers() error {
	query := `
        SELECT address, ens, loopringEns, loopringId FROM peers
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
			peer.LoopringID == -1 {
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

func (p *Peers) SavePeers(batch []*Peer) error {
	queryTemplate := `
    INSERT INTO peers (address, ens, loopringEns, loopringId)
    VALUES %s
    ON CONFLICT (address) DO UPDATE SET
        ens = EXCLUDED.ens,
        loopringEns = EXCLUDED.loopringEns,
        loopringId = EXCLUDED.loopringId
    `

	values := []any{}
	placeholders := ""

	for i, peer := range batch {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		values = append(values, peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
	}
	query := fmt.Sprintf(queryTemplate, placeholders)
	_, err := p.Db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to save peers batch: %w", err)
	}
	fmt.Printf("%d peers saved to the database\n", len(batch))
	return nil
}
