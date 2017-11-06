package slaves

import (
	"time"
)

var (
	defaultTime = time.Millisecond * 200
)

type SlavePool struct {
	pool *Pool
	stop chan struct{}
	Work func(interface{})
}

func (sp *SlavePool) Open() {
	if sp.pool != nil {
		panic("pool already running")
	}

	sp.pool = new(Pool)
	for i := 0; i < 5; i++ {
		go sp.pool.Make(sp.Work)
	}

	sp.stop = make(chan struct{})
	go func() {
		for {
			select {
			case <-sp.stop:
				return
			default:
				time.Sleep(defaultTime)
			}
			s := sp.pool.Get()
			if s != nil {
				if time.Since(s.lastUsage) > defaultTime {
					s.Close()
					s = nil
				} else {
					sp.pool.Put(s)
				}
			}
		}
	}()
}

func (sp *SlavePool) Serve(job interface{}) bool {
	s := sp.pool.Get()
	if s == nil {
		for i := 0; i < 5; i++ {
			s = sp.pool.Make(sp.Work)
		}
	}
	if s == nil {
		return false
	}
	s.ch <- job

	return true
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
