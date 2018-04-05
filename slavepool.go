package slaves

import (
	"sync"
)

var (
	// ChanSize is used in slave channel buffer size
	ChanSize = 20
	pool     = &sync.Pool{
		New: func() interface{} {
			return &slave{
				ch: make(chan interface{}, ChanSize),
			}
		},
	}
)

// This library follows the FILO scheme
type slave struct {
	ch chan interface{}
	sp *SlavePool
}

func (s *slave) work() {
	var job interface{}
	for job = range s.ch {
		if job == nil {
			return
		}

		s.sp.work(job)
	}
}

// SlavePool
type SlavePool struct {
	work func(interface{})
}

// NewPool creates SlavePool.
//
// returns nil if w is nil
func NewPool(w func(interface{})) *SlavePool {
	if w == nil {
		return nil
	}
	return &SlavePool{work: w}
}

// Serve executes job in w func
func (sp *SlavePool) Serve(job interface{}) {
	for {
		sv := pool.Get().(*slave)
		if sv.sp == nil {
			sv.sp = sp
			go sv.work()
		}
		select {
		case sv.ch <- job:
			pool.Put(sv)
			return
		default:
			pool.Put(sv)
		}
	}
}
