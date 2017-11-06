package slaves

import (
	"time"
)

var (
	defaultTime = time.Second
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

	sp.pool = &Pool{
		f: sp.Work,
	}
	for i := 0; i < 5; i++ {
		go sp.pool.Make()
	}

	sp.stop = make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-sp.stop:
				return
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
		s = sp.pool.Make()
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
