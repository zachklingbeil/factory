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

func (m *Math) Up(from int64, callback func(int64)) {
	for i := from; ; i++ {
		fmt.Println(i)
		callback(i)
	}
}

func (m *Math) Down(from int64, callback func(int64)) {
	for i := from; i >= 1; i-- {
		fmt.Println(i)
		callback(i)
	}
}
