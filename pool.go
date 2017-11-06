package slaves

import (
	"sync"
	"time"
)

type Pool struct {
	lock   sync.RWMutex
	slaves []*Slave
}

func (p *Pool) Get() *Slave {
	p.lock.Lock()
	defer p.lock.Unlock()

	n := len(p.slaves) - 1
	if n < 0 {
		return nil
	}

	s := p.slaves[n]
	if s != nil {
		p.slaves = p.slaves[:n]
	}

	return s
}

func (p *Pool) Put(s *Slave) {
	if p.slaves == nil {
		p.slaves = make([]*Slave, 0)
	}
	p.lock.Lock()
	p.slaves = append(p.slaves, s)
	p.lock.Unlock()
}

func (p *Pool) Make(f func(interface{})) *Slave {
	s := &Slave{
		Work:      f,
		lastUsage: time.Now(),
	}
	s.Open()

	go func() {
		var r interface{}
		for r = range s.ch {
			s.Work(r)
			s.lastUsage = time.Now()
			p.Put(s)
			r = nil
		}
	}()
	return s
}

func (p *Pool) Len() int {
	return len(p.slaves)
}
