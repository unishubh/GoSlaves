package slaves

import (
	"runtime"
	"sync"
)

var (
	// ChanSize is used in slave channel buffer size
	ChanSize = 20
)

// This library follows the FIFO scheme
type slave struct {
	ch chan interface{}
}

func newSlave(w func(interface{})) *slave {
	s := &slave{
		ch: make(chan interface{}, ChanSize),
	}
	go func() {
		var job interface{}
		for job = range s.ch {
			w(job)
		}
	}()
	return s
}

func (s *slave) close() {
	close(s.ch)
}

// SlavePool
type SlavePool struct {
	sv   []*slave
	i, n int

	locker sync.Mutex
}

// NewPool creates SlavePool.
//
// returns nil if w is nil
func NewPool(w func(interface{})) *SlavePool {
	if w == nil {
		return nil
	}
	n := runtime.GOMAXPROCS(0)

	sp := &SlavePool{
		i:  1,
		n:  n,
		sv: make([]*slave, n, n),
	}

	for i := 0; i < n; i++ {
		sp.sv[i] = newSlave(w)
	}

	return sp
}

// Serve sends work to goroutine pool
func (sp *SlavePool) Serve(w interface{}) {
	sp.locker.Lock()
	s := sp.sv[0]
	sp.sv[0], sp.sv[sp.i] = sp.sv[sp.i], sp.sv[0]
	sp.i++
	if sp.i == sp.n {
		sp.i = 1
	}
	sp.locker.Unlock()
	s.ch <- w
}

// Close closes the SlavePool
func (sp *SlavePool) Close() {
	for i := 0; i < sp.n; i++ {
		sp.sv[i].close()
	}
}
