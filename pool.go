package slaves

import (
	"sync"
	"time"
)

var (
	defaultTime = time.Millisecond * 200
)

type SlavePool struct {
	pool sync.Pool
	stop chan struct{}
	Work func(interface{})
}

func (sp *SlavePool) Open() {
	for i := 0; i < 5; i++ {
		s := &Slave{
			Owner: sp,
		}
		s.Open()
		sp.pool.Put(s)
		s = nil
	}

	sp.stop = make(chan struct{})
	go func() {
		for {
			select {
			case <-sp.stop:
				return
			}
			st := sp.pool.Get()
			if st != nil {
				s := st.(*Slave)
				if time.Since(s.lastUsage) > defaultTime {
					s.Close()
				} else {
					sp.pool.Put(s)
				}
			}
		}
	}()
}

func (sp *SlavePool) SendWork(job interface{}) bool {
	s := sp.pool.Get()
	if s == nil {
		for i := 0; i < 5; i++ {
			t := &Slave{
				Owner: sp,
			}
			t.Open()
			sp.pool.Put(t)
			t = nil
		}
		s = sp.pool.Get()
	}
	if s == nil {
		return false
	}
	s.(*Slave).ch <- job

	return true
}

func (sp *SlavePool) Close() {
	sp.stop <- struct{}{}
	for {
		s := sp.pool.Get()
		if s == nil {
			break
		}
		s.(*Slave).Close()
	}
}
