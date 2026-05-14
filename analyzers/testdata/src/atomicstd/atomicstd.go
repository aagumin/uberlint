package atomicstd

import "sync/atomic"

var counter int64

func bad() {
	atomic.AddInt64(&counter, 1)    // want "use atomic.Int64 type"
	atomic.LoadInt64(&counter)      // want "use atomic.Int64 type"
	atomic.StoreInt64(&counter, 10) // want "use atomic.Int64 type"
}

func good() {
	var ac atomic.Int64
	ac.Add(1)
	ac.Load()
	ac.Store(10)
}
