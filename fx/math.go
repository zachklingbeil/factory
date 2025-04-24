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

// Up counts up from `high` to `max`
func (m *Math) Up(high, max int64, callback func(int64)) {
	for i := high; i <= max; i++ {
		fmt.Println(i)
		callback(i)
	}
}

// Down counts down from `low` to 1
func (m *Math) Down(low int64, callback func(int64)) {
	for i := low; i >= 1; i-- {
		fmt.Println(i)
		callback(i)
	}
}
