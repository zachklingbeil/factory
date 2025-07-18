package zero

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Zero struct {
	One *atomic.Bool
}

func (z *Zero) Knowledge(n int) {
	for range n {
		z.One.Store(true)
	}
}

func (z *Zero) Proof(n int) {
	start := time.Now().UnixMicro()
	for range n {
		z.One.Store(true)
	}
	elapsed := time.Now().UnixMicro() - start
	nsPerOp := float64(elapsed*1000) / float64(n)
	opsPerNs := uint8(1 / nsPerOp)
	fmt.Printf("%d per nano\n", opsPerNs)
}

func (z *Zero) Run(n int) {
	for {
		z.Proof(n)
	}
}
