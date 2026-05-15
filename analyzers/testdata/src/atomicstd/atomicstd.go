package atomicstd

import (
	"sync/atomic"
	atomicx "sync/atomic"
)

var counter int64
var unsignedCounter uint64

func bad() {
	atomic.AddInt64(&counter, 1)                        // want "use atomic.Int64 type"
	atomic.LoadInt64(&counter)                          // want "use atomic.Int64 type"
	atomic.StoreInt64(&counter, 10)                     // want "use atomic.Int64 type"
	atomic.CompareAndSwapUint64(&unsignedCounter, 1, 2) // want "use atomic.Uint64 type"
	atomicx.SwapUint64(&unsignedCounter, 3)             // want "use atomic.Uint64 type"
}

func good() {
	var ac atomic.Int64
	ac.Add(1)
	ac.Load()
	ac.Store(10)
}

type fakeAtomic struct{}

func (fakeAtomic) AddInt64(*int64, int64) int64 { return 0 }

func goodNonAtomicSelector() {
	var fake fakeAtomic
	fake.AddInt64(&counter, 1)
}
