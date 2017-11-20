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
		for !sp.stop {
			time.Sleep(sp.Timeout)
			for {
				if atomic.LoadInt32(&sp.working) <= sp.MinSlaves {
					break
				}
				s := sp.ready.Get()
				if s == nil {
					break
				}
				s.(chan interface{}) <- nil
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
						ch := make(chan interface{}, 1)
						sp.ready.Put(ch)
						go func() {
							defer atomic.AddInt32(&sp.working, -1)
							for t := range ch {
								if t == nil {
									close(ch)
									return
								}
								sp.Work(t)
								sp.ready.Put(ch)
							}
						}()
						ch <- task
					} else {
						// re-add to task queue
						sp.pool.Put(task)
					}
				} else {
					s.(chan interface{}) <- task
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
