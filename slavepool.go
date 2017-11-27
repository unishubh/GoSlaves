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
	lock  sync.Mutex
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
	sp.ch = make(chan interface{}, 2)
	sp.stop = make(chan struct{})
	if sp.timeout <= 0 {
		sp.timeout = defaultTime
	}

	go func() {
		var tmp []*slave
		for {
			time.Sleep(sp.timeout)
			sp.lock.Lock()
			i, ready := 0, sp.ready
			n := len(ready)
			for c := 0; c < n; c++ {
				if time.Since(ready[i].lastUsage) > sp.timeout {
					i++
				}
			}
			tmp = append(tmp[:0], ready[:i]...)
			if i > 0 {
				m := copy(ready, ready[i:])
				for i = m; i < n; i++ {
					ready = nil
				}
				sp.ready = ready
			}
			sp.lock.Unlock()

			for _, ch := range tmp {
				ch.ch <- nil
			}
		}
	}()

	go func() {
		for {
			select {
			case job := <-sp.ch:
				sv := &slave{}
				sp.lock.Lock()
				n := len(sp.ready) - 1
				if n <= 0 {
					sv.ch = make(chan interface{}, 1)
					sv.lastUsage = time.Now()
					go func(s *slave) {
						for job := range s.ch {
							if job == nil {
								close(s.ch)
								return
							}
							sp.Work(job)
							s.lastUsage = time.Now()

							sp.lock.Lock()
							sp.ready = append(sp.ready, s)
							sp.lock.Unlock()
						}
					}(sv)
				} else {
					sv = sp.ready[n]
					sp.ready = sp.ready[:n]
				}
				sp.lock.Unlock()
				sv.ch <- job
			case <-sp.stop:
				close(sp.stop)
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
		sv.ch <- nil
	}
	close(sp.ch)

	sp.stop <- struct{}{}
}

func (sp *SlavePool) SetTimeout(secs int) {
	sp.timeout = time.Duration(secs) * time.Second
}
