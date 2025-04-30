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

func (m *Math) StringToInt(value string) *big.Int {
	i := new(big.Int)
	_, ok := i.SetString(value, 10)
	if !ok {
		log.Error("Failed to convert string to big.Int: %s", value)
		return nil
	}
	return i
}

// func (m *Math) Up(from int64, callback func(int64), stop *bool) {
// 	for i := from; ; i++ {
// 		m.When.L.Lock()
// 		if *stop {
// 			m.When.L.Unlock()
// 			break
// 		}
// 		m.When.L.Unlock()

// 		fmt.Println(i)
// 		callback(i)
// 	}
// }

// func (m *Math) Down(from int64, callback func(int64)) {
// 	for i := from; i >= 1; i-- {
// 		fmt.Println(i)
// 		callback(i)
// 	}
// }
