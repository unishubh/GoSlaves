package slaves

import (
	"runtime"
)

var (
	// ChanSize is used in slave channel buffer size
	ChanSize = 20
)

// This library follows the FILO scheme
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
	sv []*slave

	// W is working channel
	W chan interface{}
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
		W:  make(chan interface{}, ChanSize*n),
		sv: make([]*slave, 0, n),
	}

	for i := 0; i < n; i++ {
		sp.sv = append(sp.sv, newSlave(w))
	}

	go func() {
		var s *slave
		i := 1
		for w := range sp.W {
			if i == n {
				i = 1
			}
			s = sp.sv[0]
			sp.sv[0], sp.sv[i] = sp.sv[i], sp.sv[0]
			s.ch <- w
			i++
		}
	}()

	return sp
}

func (s *SlavePool) Close() {
	n := len(s.sv)
	for i := 0; i < n; i++ {
		s.sv[i].close()
	}
}
