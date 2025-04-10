package peer

import (
	"fmt"
	"sync"

	_ "github.com/lib/pq"

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
	Address     string `json:"address"`
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  int64  `json:"loopringId"`
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

	peers.PeerChan = make(chan string, len(peers.Addresses))
	for _, address := range peers.Addresses {
		peers.PeerChan <- address
	}

	return peers
}

func (p *Peers) NewBlock(addresses []string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	peers := len(addresses)
	fmt.Printf("%d new peers\n", peers)

	for _, address := range addresses {
		if _, exists := p.Map[address]; !exists {
			p.Map[address] = &Peer{Address: address}
			p.Addresses = append(p.Addresses, address)
		}
		p.PeerChan <- address
	}
}

func (p *Peers) HelloUniverse() {
	batchSize := 1000
	var batch []*Peer

	p.Mu.RLock()
	peers := len(p.Addresses)
	p.Mu.RUnlock()
	fmt.Printf("%d peers to process\n", peers)

	for {
		if peers == 0 && len(batch) == 0 {
			break
		}

		if len(batch) >= batchSize || (peers == 0 && len(batch) > 0) {
			if err := p.SavePeers(batch); err != nil {
				fmt.Printf("Error saving batch: %v\n", err)
			}
			batch = batch[:0]
		}

		if peers > 0 {
			address := <-p.PeerChan

			p.Mu.Lock()
			peer := p.Map[address]
			p.Mu.Unlock()

			p.GetENS(peer, peer.Address)
			p.GetLoopringENS(peer, peer.Address)
			p.GetLoopringID(peer, peer.Address)
			batch = append(batch, peer)

			fmt.Printf("%d %s %s %d\n", peers, peer.ENS, peer.LoopringENS, peer.LoopringID)
			peers--
		}
	}
	fmt.Println("Hello Universe")
}
