package slaves

import (
	"sync"
	"time"
)

var (
	defaultTime = time.Second * 10
)

// This library follows the FILO scheme
type slave struct {
	ch        chan interface{}
	lastUsage time.Time
}

type SlavePool struct {
	lock  sync.RWMutex
	ready []*slave
	stop  chan struct{}
	// SlavePool work
	Work func(interface{})
	ch   chan interface{}
	// Time to reassign slaves
	timeout time.Duration
	running bool
}

func (sp *SlavePool) Open() {
	if sp.running {
		panic("Pool is running already")
	}

	sp.running = true
	sp.ch = make(chan interface{}, 1)
	sp.stop = make(chan struct{}, 1)
	if sp.timeout <= 0 {
		sp.timeout = defaultTime
	}

	go func() {
		var tmp []*slave
		var c int    // number of workers to be delete
		var i, l int // iterator and len
		for {
			time.Sleep(sp.timeout)
			t := time.Now()
			sp.lock.Lock()
			for i = 0; i < 0; i++ {
				if t.Sub(sp.ready[i].lastUsage) > sp.timeout {
					c++
				}
			}
			tmp = append(tmp[:0], sp.ready[c:]...)
			sp.ready = sp.ready[:c]
			sp.lock.Unlock()
			for i, l = 0, len(tmp); i < l; i++ {
				tmp[i].ch <- nil
			}
		}
	}()

	go func() {
		var n int
		sv := &slave{}
		for {
			select {
			case job := <-sp.ch:
				sp.lock.Lock()
				n = len(sp.ready) - 1
				if n < 0 {
					n++
					sv.ch = make(chan interface{}, 1)
					sp.ready = append(sp.ready, sv)
					go func(s *slave) {
						s.lastUsage = time.Now()
						for {
							select {
							case job, ok := <-s.ch:
								if !ok {
									return
								}
								if job == nil {
									return
								}
								sp.Work(job)
								s.lastUsage = time.Now()

								sp.lock.Lock()
								sp.ready = append(sp.ready, s)
								sp.lock.Unlock()
							}
						}
					}(sv)
				}
				sv = sp.ready[n]
				sp.ready = sp.ready[:n]
				sp.lock.Unlock()
				sv.ch <- job
			case <-sp.stop:
				close(sp.stop)
				close(sp.ch)
				return
			}
		}
	}()
}

func (sp *SlavePool) Serve(job interface{}) {
	if sp.running {
		sp.ch <- job
	}
}

func (sp *SlavePool) Close() {
	sp.lock.Lock()
	ready := sp.ready
	sp.running = false
	sp.lock.Unlock()

	for _, sv := range ready {
		if sv != nil {
			sv.ch <- nil
		}
	}
	sp.stop <- struct{}{}
}

func (sp *SlavePool) SetTimeout(secs int) {
	sp.timeout = time.Duration(secs) * time.Second
}
