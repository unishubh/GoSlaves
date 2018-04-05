package slaves

import (
	"sync"
)

var (
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

		s.sp.Work(job)
	}
}

type SlavePool struct {
	// SlavePool work
	Work func(interface{})

	running bool
}

func (sp *SlavePool) Open() {
	if sp.running {
		panic("Pool is running already")
	}

	sp.running = true
}

func (sp *SlavePool) Serve(job interface{}) {
	if sp.running {
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
			}
			pool.Put(sv)
		}
	}
}

func (sp *SlavePool) Close() {
	sp.running = false
}
