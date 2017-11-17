package slaves

import (
	"sync"
	"time"
)

type Pool struct {
	ck     sync.Mutex
	stack  []interface{}
	slaves []*Slave
	f      func(interface{})
}

func (p *Pool) AddStack(job interface{}) {
	p.ck.Lock()
	if p.stack == nil {
		p.stack = make([]interface{}, 0)
	}
	p.stack = append(p.stack, job)
	p.ck.Unlock()
}

func (p *Pool) GetStack() interface{} {
	p.ck.Lock()
	defer p.ck.Unlock()
	if p.stack == nil {
		return nil
	}
	n := p.StackLen()
	r := p.stack[n]
	p.stack = p.stack[:n]
	return r
}

func (p *Pool) StackLen() int {
	p.ck.Lock()
	defer p.ck.Unlock()
	return len(p.stack)
}

func (p *Pool) Get() *Slave {
	p.ck.Lock()
	defer p.ck.Unlock()

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

	p.ck.Lock()
	p.slaves = append(p.slaves, s)
	p.ck.Unlock()
}

func (p *Pool) Make() *Slave {
	s := &Slave{
		lastUsage: time.Now(),
	}
	s.Open()

	go func() {
		var r interface{}
		for {
			select {
			case r = <-s.ch:
				p.f(r)
				s.lastUsage = time.Now()
				p.Put(s)
				r = nil
			}
		}
	}()

	return s
}

func (p *Pool) Len() int {
	p.ck.Lock()
	defer p.ck.Unlock()
	return len(p.slaves)
}
