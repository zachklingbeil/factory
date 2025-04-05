package fx

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wealdtech/go-ens/v3"
)

type Peers struct {
	Json           *JSON
	Eth            *ethclient.Client
	Db             *Database
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

func NewPeers(json *JSON, eth *ethclient.Client, db *Database) *Peers {
	peers := &Peers{
		Json:           json,
		Eth:            eth,
		LoopringApiKey: os.Getenv("LOOPRING_API_KEY"),
		Map:            make(map[string]*Peer),
		Db:             db,
	}

	// Ensure the peers table exists
	if err := peers.CreateTable(); err != nil {
		fmt.Printf("Error ensuring peers table exists: %v\n", err)
	}

	// Try to load the map from the database
	if err := peers.LoadMap(); err != nil {
		fmt.Printf("Error loading map from database: %v\n", err)
	} else {
		fmt.Println("Map loaded successfully from the database.")
	}

	// Start periodic checkpointing to save the map to the database
	peers.Checkpoint(20) // Save every 60 seconds

	return peers
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

func (p *Peers) LoadMap() error {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	// Query the map data from the database
	query := `SELECT data FROM peers ORDER BY updated_at DESC LIMIT 1`
	var jsonData []byte
	err := p.Db.QueryRow(query).Scan(&jsonData)
	if err != nil {
		if err == sql.ErrNoRows {
			// No data in the database, initialize an empty map
			p.Map = make(map[string]*Peer)
			return nil
		}
		return fmt.Errorf("failed to load map from database: %w", err)
	}

	// Deserialize the JSON data into the map
	err = json.Unmarshal(jsonData, &p.Map)
	if err != nil {
		return fmt.Errorf("failed to deserialize map: %w", err)
	}

	return nil
}

func (p *Peers) SaveMap() error {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	jsonData, err := json.Marshal(p.Map)
	if err != nil {
		return fmt.Errorf("failed to serialize map: %w", err)
	}

	// Insert or update the map in the database
	query := `
        INSERT INTO peers (data, updated_at)
        VALUES ($1, NOW())
        ON CONFLICT (id) DO UPDATE
        SET data = EXCLUDED.data, updated_at = EXCLUDED.updated_at
    `
	_, err = p.Db.Exec(query, jsonData)
	if err != nil {
		return fmt.Errorf("failed to save map to database: %w", err)
	}

	return nil
}

func (p *Peers) Checkpoint(intervalSeconds int) {
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	go func() {
		for range ticker.C {
			err := p.SaveMap()
			if err != nil {
				fmt.Printf("Failed to save map to database: %v\n", err)
			} else {
				fmt.Println("Map saved to database successfully.")
			}
		}
	}()
}

func (p *Peers) CreateTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS peers (
        id SERIAL PRIMARY KEY,
        data JSONB NOT NULL,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`
	_, err := p.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create peers table: %w", err)
	}
	return nil
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
