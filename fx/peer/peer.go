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

	if err := peers.LoadMap(); err != nil {
		fmt.Printf("Error loading map from database: %v\n", err)
	} else {
		fmt.Println("Map loaded successfully from the database.")
	}

	go peers.Checkpoint(20)
	return peers
}

func (p *Peers) HelloUniverse(value string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	formattedValue := p.Format(value)
	peer := p.getOrCreatePeer(formattedValue)

	if peer.Address == "" {
		p.GetLoopringAddress(peer, value)
	}
	if peer.ENS == "" {
		p.GetENS(peer, peer.Address)
	}
	if peer.LoopringENS == "" {
		p.GetLoopringENS(peer, peer.Address)
	}
	if peer.LoopringID == "" {
		p.GetLoopringID(peer, peer.Address)
	}
}

func (p *Peers) getOrCreatePeer(value string) *Peer {
	if peer, exists := p.Map[value]; exists {
		return peer
	}

	peer := &Peer{Address: value}
	p.Map[value] = peer
	return peer
}

func (p *Peers) BatchCreatePeers(addresses []string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	for _, address := range addresses {
		formattedAddress := p.Format(address)
		if _, exists := p.Map[formattedAddress]; !exists {
			p.Map[formattedAddress] = &Peer{Address: formattedAddress}
		}
	}
}
