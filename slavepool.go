package slaves

import (
	"runtime"
)

// This library uses a queue system. See Serve function.
type slave struct {
	ch chan interface{}
}

func newSlave(w func(interface{})) *slave {
	s := &slave{
		ch: make(chan interface{}, 1),
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
	sv []*slave
	n  int
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
		n:  n,
		sv: make([]*slave, n, n),
	}

	for i := 0; i < n; i++ {
		sp.sv[i] = newSlave(w)
	}

	return sp
}

// Serve sends work to goroutine pool
//
// If all slaves are busy this function will stop until any of this ends a task.
func (sp *SlavePool) Serve(w interface{}) {
	i := 0
	for {
		select {
		case sp.sv[i].ch <- w:
			return
		default: // channel is busy
			i++
			if i == sp.n {
				i = 0
			}
		}
	}
}

// Close closes the SlavePool
func (sp *SlavePool) Close() {
	for i := 0; i < sp.n; i++ {
		sp.sv[i].close()
	}
}
