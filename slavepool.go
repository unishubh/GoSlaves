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
			default:
				time.Sleep(defaultTime)
			}
			var s *Slave
			sp.pool.ck.Lock()
			for i := 0; i < len(sp.pool.slaves); i++ {
				s = sp.pool.slaves[i]
				if s == nil {
					sp.pool.slaves = sp.pool.slaves[:i+
						copy(sp.pool.slaves[i:], sp.pool.slaves[i+1:])]
					i--
				} else {
					if time.Since(s.lastUsage) > defaultTime {
						s.Close()
						sp.pool.slaves = sp.pool.slaves[:i+
							copy(sp.pool.slaves[i:], sp.pool.slaves[i+1:])]
						i--
					}
				}
				s = nil
			}
			sp.pool.ck.Unlock()
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
