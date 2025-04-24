package fx

import (
	"fmt"
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

func (m *Math) Up(high, max int64) error {
	for i := high; i <= max; i++ {
		fmt.Println(i)
	}
	return nil
}

func (m *Math) Down(low int64) error {
	for i := low; i >= 1; i-- {
		fmt.Println(i)
	}
	return nil
}
