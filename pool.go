package slaves

import (
	"sync"
)

type SlavePool struct {
	lock    sync.Mutex
	slaves  []*Slave
	ready   []*Slave
	works   int
	stop    chan struct{}
	Workers int
	Work    func(interface{}) interface{}
}

func (sp *SlavePool) Open() {
	if sp.Workers <= 0 {
		sp.Workers = 256
	}

	sp.slaves = make([]*Slave, sp.Workers)
	for i := range sp.slaves {
		sp.slaves[i] = &Slave{
			Owner: sp,
		}
		sp.slaves[i].Open()
		sp.ready = append(sp.ready, sp.slaves[i])
	}

	sp.stop = make(chan struct{})
	go func() {
		for {
			select {
			case <-sp.stop:
				return
			}
			n := len(sp.slaves)
			if n > sp.Workers {
				sp.lock.Lock()
				i := n - 1
				for ; i > sp.Workers; i-- {
					sp.slaves[i].Close()
					sp.slaves[i] = nil
				}
				sp.slaves = sp.slaves[:i]
				sp.lock.Unlock()
			}
		}
	}()
}

func (sp *SlavePool) SendWork(job interface{}) bool {
	sp.lock.Lock()
	defer sp.lock.Unlock()

	n := len(sp.ready) - 1
	if n < 0 {
		return false
	}

	s := sp.ready[n]
	s.ch <- job
	sp.ready = sp.ready[:n]
	return true
}

func (sp *SlavePool) Close() {
	sp.stop <- struct{}{}
	for _, c := range sp.slaves {
		c.Close()
	}
}
