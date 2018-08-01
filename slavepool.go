package slaves

import (
	"runtime"
)

// This library uses a queue system. See Serve function.
type slave struct {
	ch chan interface{}
}

func newSlave(w func(interface{})) (s slave) {
	s.ch = make(chan interface{}, 1)
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

// Pool slaves
type Pool struct {
	sv []slave
	n  int
}

// NewPool creates SlavePool.
//
// if workers is 0 default workers will be created
// use workers var if you know what you are doing
func NewPool(workers int, w func(interface{})) (p Pool) {
	if w == nil {
		return
	}
	if workers <= 0 {
		workers = runtime.GOMAXPROCS(0)
	}

	p.n = workers
	p.sv = make([]slave, p.n, p.n)
	for i := 0; i < p.n; i++ {
		p.sv[i] = newSlave(w)
	}
	return
}

// Serve sends work to goroutine pool
//
// If all slaves are busy this function will stop until any of this ends a task.
func (p *Pool) Serve(w interface{}) {
	i := 0
	for {
		select {
		case p.sv[i].ch <- w:
			return
		default: // channel is busy
			i++
			if i == p.n {
				i = 0
			}
		}
	}
}

// ServeNonStop returns true if work have been sended to the goroutine pool.
//
// This function returns a state and does not block the workflow.
func (p *Pool) ServeNonStop(w interface{}) bool {
	i := 0
	for i < p.n {
		select {
		case p.sv[i].ch <- w:
			return true
		default:
			i++
		}
	}
	return false
}

// Close closes the SlavePool
func (p *Pool) Close() {
	for i := 0; i < p.n; i++ {
		p.sv[i].close()
	}
}
