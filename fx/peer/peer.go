package peer

import (
	"fmt"
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
	Address     string
	ENS         string
	LoopringENS string
	LoopringID  string
}

func NewPeers(json *fx.JSON, eth *ethclient.Client, db *fx.Database) *Peers {
	peers := &Peers{
		Json:      json,
		Eth:       eth,
		Map:       make(map[string]*Peer),
		Addresses: nil,
		Db:        db,
		PeerChan:  make(chan string, 100),
	}

	if err := peers.LoadPeers(); err != nil {
		fmt.Printf("Error loading peers: %v\n", err)
	}

	return peers
}

func (p *Peers) HelloUniverse() {
	batchSize := 1000
	var batch []*Peer

	for {
		if len(batch) > 0 {
			if err := p.SavePeers(batch); err != nil {
				fmt.Printf("Error saving final batch: %v\n", err)
			}
			batch = batch[:0]
		}
		address := <-p.PeerChan

		p.Mu.Lock()
		if _, exists := p.Map[address]; !exists {
			p.Map[address] = &Peer{Address: address}
			p.Addresses = append(p.Addresses, address)
		}
		peer := p.Map[address]
		p.Mu.Unlock()

		p.GetENS(peer, peer.Address)
		p.GetLoopringENS(peer, peer.Address)
		p.GetLoopringID(peer, peer.Address)

		batch = append(batch, peer)

		if len(batch) >= batchSize {
			if err := p.SavePeers(batch); err != nil {
				fmt.Printf("Error saving batch: %v\n", err)
			}
		}

		if len(batch) > 0 {
			if err := p.SavePeers(batch); err != nil {
				fmt.Printf("Error saving final batch: %v\n", err)
			}
		}
	}
}

// NewBlock sends addresses through the channel, updates p.Addresses for new peers.
func (p *Peers) NewBlock(addresses []string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	for _, address := range addresses {
		if _, exists := p.Map[address]; !exists {
			p.Map[address] = &Peer{Address: address}
			p.Addresses = append(p.Addresses, address)
			p.PeerChan <- address
		}
	}
}
