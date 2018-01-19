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

func (s *slave) close() {
	if s.ch != nil {
		close(s.ch)
		s.ch = nil
	}
}

func cleanSlaves(locker sync.RWMutex, slaves *[]*slave) {
	var tmp []*slave
	var c int    // number of workers to be delete
	var i, l int // iterator and len
	for {
		time.Sleep(defaultTime)
		if len(*slaves) == 0 {
			continue
		}
		t := time.Now()
		locker.Lock()
		for i = 0; i < 0; i++ {
			if t.Sub((*slaves)[i].lastUsage) > defaultTime {
				c++
			}
		}
		tmp = append(tmp[:0], (*slaves)[c:]...)
		*slaves = (*slaves)[:c]
		locker.Unlock()
		for i, l = 0, len(tmp); i < l; i++ {
			// closing
			tmp[i].ch <- nil
		}
		tmp = nil
	}
}

type SlavePool struct {
	lock  sync.RWMutex
	ready []*slave
	stop  chan struct{}
	// SlavePool work
	Work    func(interface{})
	ch      chan interface{}
	running bool
}

func (sp *SlavePool) Open() {
	if sp.running {
		panic("Pool is running already")
	}

	sp.running = true
	sp.ch = make(chan interface{}, 1)
	sp.stop = make(chan struct{}, 1)
	sp.ready = make([]*slave, 0)

	go cleanSlaves(sp.lock, &sp.ready)

	go func() {
		var n int
		var job interface{}
		sv := &slave{}
		for {
			select {
			case job = <-sp.ch:
				sp.lock.Lock()
				n = len(sp.ready) - 1
				if n < 0 {
					n++
					sv.ch = make(chan interface{}, 1)
					sp.ready = append(sp.ready, sv)
					go func(s *slave) {
						var job interface{}
						s.lastUsage = time.Now()
						for job = range s.ch {
							if job == nil {
								s.close()
								return
							}
							sp.Work(job)
							s.lastUsage = time.Now()

							sp.lock.Lock()
							sp.ready = append(sp.ready, s)
							sp.lock.Unlock()
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
	sp.stop <- struct{}{}

	sp.lock.Lock()
	ready := sp.ready
	sp.running = false

	for _, sv := range ready {
		if sv != nil {
			sp.lock.Unlock()
			sv.close()
			sp.lock.Lock()
		}
	}
	sp.lock.Unlock()
}
