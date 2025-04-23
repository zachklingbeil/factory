package fx

import (
	"math/big"

	"github.com/ethereum/go-ethereum/log"
)

type Math struct{}

func (m *Math) Int(value string) *big.Int {
	bigIntValue := new(big.Int)
	if _, ok := bigIntValue.SetString(value, 10); !ok {
		log.Error("Failed to convert string to big.Int: %s", value)
	}
	return bigIntValue
}
