package slaves

import (
	"time"
)

var (
	defaultTime = time.Second * 10
)

type SlavePool struct {
	pool *Pool
	stop chan struct{}
	// SlavePool work
	Work func(interface{})
	// Limit of slaves
	Limit int
	// Minimum number of slaves
	// waiting for tasks
	MinSlaves int
	// Time to reassign slaves
	Timeout time.Duration
}

func (sp *SlavePool) Open() {
	if sp.pool != nil {
		panic("pool already running")
	}
	if sp.Timeout <= 0 {
		sp.Timeout = defaultTime
	}
	if sp.MinSlaves <= 0 {
		sp.MinSlaves = 5
	}

	sp.pool = &Pool{
		f: sp.Work,
	}
	for i := 0; i < sp.MinSlaves; i++ {
		go sp.pool.Make()
	}

	sp.stop = make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-sp.stop:
				return
			default:
				for {
					if sp.pool.StackLen() == 0 {
						break
					}
					job := sp.pool.GetStack()
					if job != nil {
						sp.Serve(job)
					}
				}
				time.Sleep(sp.Timeout)
			}
			var s *Slave
			for i := 0; i < sp.pool.Len(); i++ {
				s = sp.pool.slaves[i]
				if s == nil {
					sp.pool.slaves = sp.pool.slaves[:i+
						copy(sp.pool.slaves[i:], sp.pool.slaves[i+1:])]
					i--
				} else {
					if sp.pool.Len() > sp.MinSlaves && time.Since(s.lastUsage) > sp.Timeout {
						sp.pool.ck.Lock()
						s.Close()
						sp.pool.slaves = sp.pool.slaves[:i+
							copy(sp.pool.slaves[i:], sp.pool.slaves[i+1:])]
						i--
						sp.pool.ck.Unlock()
					}
				}
				s = nil
			}
		}
	}()
}

func (sp *SlavePool) Serve(job interface{}) bool {
	s := sp.pool.Get()
	if s == nil {
		l := sp.pool.Len()
		if l >= sp.Limit {
			return false
		}
		s = sp.pool.Make()
	}
	if s == nil {
		return false
	}
	s.ch <- job

	return true
}

func (sp *SlavePool) Enqueue(job interface{}) {
	sp.pool.AddStack(job)
}

func (sp *SlavePool) Close() {
	sp.stop <- struct{}{}
	for {
		s := sp.pool.Get()
		if s == nil {
			break
		}
		s.Close()
	}
}
