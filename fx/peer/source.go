package peer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens/v3"
)

func (p *Peers) Format(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
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

// hex -> ENS [.eth], "." for no ENS, "!" for errors
func (p *Peers) GetENS(peer *Peer, address string) {
	addr := common.HexToAddress(address)
	ensName, err := ens.ReverseResolve(p.Eth, addr)
	if err != nil {
		peer.ENS = "!"
		return
	}
	if ensName == "" {
		peer.ENS = "."
		return
	}
	peer.ENS = p.Format(ensName)
}

// hex -> LoopringENS [.loopring.eth], "." for no LoopringENS, "!" for errors
func (p *Peers) GetLoopringENS(peer *Peer, address string) {
	url := fmt.Sprintf("https://api3.loopring.io/api/wallet/v3/resolveName?owner=%s", address)
	var response struct {
		Loopring string `json:"data"`
	}
	data, err := p.Json.In(url, "")
	if err != nil {
		peer.LoopringENS = "!"
		return
	}
	if json.Unmarshal(data, &response) != nil || response.Loopring == "" {
		peer.LoopringENS = "."
		return
	}
	peer.LoopringENS = p.Format(response.Loopring)
}

// hex -> LoopringId, -1 for no Loopring ID and errors
func (p *Peers) GetLoopringID(peer *Peer, address string) {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/account?owner=%s", address)
	var response struct {
		ID int64 `json:"accountId"`
	}
	data, err := p.Json.In(url, "")
	if err != nil || json.Unmarshal(data, &response) != nil || response.ID == 0 {
		peer.LoopringID = -1
		return
	}
	peer.LoopringID = response.ID
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
	if data, err := p.Json.In(url, ""); err == nil && json.Unmarshal(data, &response) == nil {
		peer.Address = p.Format(response.Address)
	} else {
		peer.Address = "!"
	}
}
