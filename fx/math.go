package fx

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/log"
)

type Math struct {
	When *sync.Cond
}

func NewMath(when *sync.Cond) *Math {
	return &Math{
		When: when,
	}
}

func (m *Math) Int(value string) *big.Int {
	bigIntValue := new(big.Int)
	if _, ok := bigIntValue.SetString(value, 10); !ok {
		log.Error("Failed to convert string to big.Int: %s", value)
	}
	return bigIntValue
}

func (m *Math) Up(from int64, callback func(int64), stop *bool) {
	for i := from; ; i++ {
		m.When.L.Lock()
		if *stop {
			m.When.L.Unlock()
			break
		}
		m.When.L.Unlock()

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

// panic: runtime error: invalid memory address or nil pointer dereference
// [signal SIGSEGV: segmentation violation code=0x1 addr=0x20 pc=0x7b220a]

// goroutine 93 [running]:
// github.com/zachklingbeil/block/loopring.(*Loopring).Coordinates(0xc003ad1410, 0x0)
//         /block/loopring/fx.go:63 +0x2a
// github.com/zachklingbeil/block/loopring.(*Loopring).BlockByBlock(0xc003ad1410, 0xe89b)
//         /block/loopring/fx.go:37 +0x32
// github.com/zachklingbeil/block/loopring.(*Loopring).Loop.func1(0xe89b)
//         /block/loopring/fx.go:26 +0x2b
// github.com/zachklingbeil/factory/fx.(*Math).Up(0x0?, 0x9298a1?, 0xc0057c7fa8)
//         /go/pkg/mod/github.com/zachklingbeil/factory@v1.1.37/fx/math.go:23 +0x74
// github.com/zachklingbeil/block/loopring.(*Loopring).Loop(0xc003ad1410)
//         /block/loopring/fx.go:25 +0x11a
// created by github.com/zachklingbeil/block/loopring.Connect in goroutine 1
//         /block/loopring/main.go:19 +0xdb
