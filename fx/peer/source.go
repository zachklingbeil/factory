package peer

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens/v3"
)

func (p *Peers) Format(address string) string {
	if strings.HasPrefix(address, "0x") {
		return "0x" + strings.ToLower(address[2:])
	}
	if strings.HasSuffix(address, ".eth") {
		return strings.ToLower(address)
	}
	return address
}

// ENS -> hex
func (p *Peers) GetAddress(peer *Peer, dotEth string) {
	address, err := ens.Resolve(p.Eth, dotEth)
	if err != nil {
		peer.Address = dotEth
		return
	}
	peer.Address = p.Format(address.Hex())
}

// hex -> ENS [.eth]
func (p *Peers) GetENS(peer *Peer, address string) {
	addr := common.HexToAddress(address)
	if ensName, err := ens.ReverseResolve(p.Eth, addr); err == nil {
		peer.ENS = p.Format(ensName)
	}
}

// hex -> LoopringENS [.loopring.eth]
func (p *Peers) GetLoopringENS(peer *Peer, address string) {
	url := fmt.Sprintf("https://api3.loopring.io/api/wallet/v3/resolveName?owner=%s", address)
	var response struct {
		Loopring string `json:"data"`
	}
	if data, err := p.Json.In(url, ""); err == nil && json.Unmarshal(data, &response) == nil {
		peer.LoopringENS = p.Format(response.Loopring)
	}
}

// hex -> LoopringId
func (p *Peers) GetLoopringID(peer *Peer, address string) {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/account?owner=%s", address)
	var response struct {
		ID int64 `json:"accountId"`
	}
	if data, err := p.Json.In(url, os.Getenv("LOOPRING_API_KEY")); err == nil && json.Unmarshal(data, &response) == nil {
		peer.LoopringID = strconv.FormatInt(response.ID, 10)
	}
}

// LoopringId -> hex
func (p *Peers) GetLoopringAddress(peer *Peer, id string) {
	accountID, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/account?accountId=%d", accountID)
	var response struct {
		Address string `json:"owner"`
	}
	if data, err := p.Json.In(url, os.Getenv("LOOPRING_API_KEY")); err == nil && json.Unmarshal(data, &response) == nil {
		peer.Address = p.Format(response.Address)
	}
}
