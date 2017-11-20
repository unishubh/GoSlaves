package slaves

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	defaultTime = time.Second * 5
)

type SlavePool struct {
	pool    sync.Pool
	ready   sync.Pool
	stop    bool
	working int32
	// SlavePool work
	Work func(interface{})
	// Limit of slaves
	Limit uint
	// Minimum number of slaves
	// waiting for tasks
	MinSlaves int
	// Time to reassign slaves
	Timeout time.Duration
}

func (sp *SlavePool) Open() {
	sp.working = 0
	if sp.Timeout <= 0 {
		sp.Timeout = defaultTime
	}
	if sp.MinSlaves <= 0 {
		sp.MinSlaves = 5
	}
	if sp.Limit <= 0 {
		sp.Limit = 1024 * 256
	}

	go func() {
		var p sync.Pool
		for {
			time.Sleep(sp.Timeout)
			for {
				sl := sp.ready.Get()
				if sl == nil {
					break
				}
				s := sl.(*Slave)
				if time.Since(s.lastUsage) > sp.Timeout {
					s.ch <- nil
					close(s.ch)
				} else {
					p.Put(s)
				}
			}
			for {
				sl := p.Get()
				if sl == nil {
					break
				}
				s := sl.(*Slave)
				sp.ready.Put(s)
			}
		}
	}()
	go func() {
		for !sp.stop {
			task := sp.pool.Get()
			if task == nil {
				time.Sleep(time.Millisecond * 20)
			} else {
				// get one slave from pool
				s := sp.ready.Get()
				if s == nil {
					if int32(sp.Limit) > atomic.LoadInt32(&sp.working) {
						// create new slave
						atomic.AddInt32(&sp.working, 1)
						go func() {
							defer atomic.AddInt32(&sp.working, -1)
							s := &Slave{
								ch:        make(chan interface{}),
								lastUsage: time.Now(),
							}
							sp.ready.Put(s)
							for t := range s.ch {
								if t == nil {
									return
								}
								sp.Work(t)
								s.lastUsage = time.Now()
								sp.ready.Put(s)
							}
						}()
						sp.pool.Put(task)
					} else {
						// re-add to task queue
						sp.pool.Put(task)
					}
				} else {
					sl := s.(*Slave)
					sl.ch <- task
				}
			}
		}
	}()
}

func (sp *SlavePool) Serve(job interface{}) {
	sp.pool.Put(job)
}

func (sp *SlavePool) Close() {
	sp.stop = true
}
