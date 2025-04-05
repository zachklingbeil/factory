package fx

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wealdtech/go-ens/v3"
)

type Peers struct {
	Json           *JSON
	Eth            *ethclient.Client
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

func NewPeers(json *JSON, eth *ethclient.Client) *Peers {
	return &Peers{
		Json:           json,
		Eth:            eth,
		LoopringApiKey: os.Getenv("LOOPRING_API_KEY"),
		Map:            make(map[string]*Peer),
	}
}

func (p *Peers) HelloUniverse(value string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	// Check if the peer exists
	peer, exists := p.Map[value]
	if !exists {
		peer = &Peer{}
		p.Map[value] = peer
	}

	// Resolve address for LoopringID's without an address
	// Update ENS|LoopringENS|LoopringId when the value hasn't been set
	if peer.Address == "" {
		p.GetLoopringAddress(peer, value)
	}

	if peer.ENS == "" && peer.Address != "" {
		p.GetENS(peer, peer.Address)
	}

	if peer.LoopringENS == "" && peer.Address != "" {
		p.GetLoopringENS(peer, peer.Address)
	}

	if peer.LoopringID == "" && peer.Address != "" {
		p.GetLoopringID(peer, peer.Address)
	}
}

// Uniform
func (p *Peers) Format(address string) string {
	// Format Ethereum addresses
	if strings.HasPrefix(address, "0x") {
		return "0x" + strings.ToUpper(address[2:])
	}

	//  Format ENS names to lowercase
	if strings.HasSuffix(address, ".eth") {
		return strings.ToLower(address)
	}
	return address
}

// ENS -> hex
func (p *Peers) GetAddress(name string) string {
	address, err := ens.Resolve(p.Eth, name)
	if err != nil {
		return name
	}
	return p.Format(address.Hex())
}

// hex -> ENS [.eth]
func (p *Peers) GetENS(peer *Peer, address string) {
	addr := common.HexToAddress(address)
	dotEth, err := ens.ReverseResolve(p.Eth, addr)
	if err == nil {
		peer.ENS = p.Format(dotEth)
		p.Map[peer.ENS] = peer
	}
}

// hex -> LoopringENS [.loopring.eth]
func (p *Peers) GetLoopringENS(peer *Peer, address string) {
	url := fmt.Sprintf("https://api3.loopring.io/api/wallet/v3/resolveName?owner=%s", address)
	var dot struct {
		Loopring string `json:"data"`
	}

	response, err := p.Json.In(url, "")
	if err == nil && json.Unmarshal(response, &dot) == nil {
		peer.LoopringENS = p.Format(dot.Loopring)
		p.Map[peer.LoopringENS] = peer
	}
}

// hex -> LoopringId
func (p *Peers) GetLoopringID(peer *Peer, address string) {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/account?owner=%s", address)
	var res struct {
		ID int64 `json:"accountId"`
	}

	response, err := p.Json.In(url, os.Getenv("LOOPRING_API_KEY"))
	if err == nil && json.Unmarshal(response, &res) == nil {
		peer.LoopringID = fmt.Sprintf("%d", res.ID)
		p.Map[peer.LoopringID] = peer
	}
}

// LoopringId -> hex
func (p *Peers) GetLoopringAddress(peer *Peer, id string) {
	accountID, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/account?accountId=%d", accountID)
	var loopring struct {
		Address string `json:"owner"`
	}

	response, err := p.Json.In(url, os.Getenv("LOOPRING_API_KEY"))
	if err == nil && json.Unmarshal(response, &loopring) == nil {
		formattedAddress := p.Format(loopring.Address)
		peer.Address = formattedAddress
		p.Map[formattedAddress] = peer
	}
}
