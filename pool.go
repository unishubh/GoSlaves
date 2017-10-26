package slaves

import (
	"github.com/themester/GoSlaves/slave"
	"math"
)

type SlavePool struct {
	slaves []*slave.Slave
}

func MakePool(num uint, work func(interface{}) interface{},
	after func(interface{})) SlavePool {

	sp := SlavePool{
		slaves: make([]*slave.Slave, num),
	}
	for i := range sp.slaves {
		sp.slaves[i] = &slave.Slave{
			Work:  work,
			After: after,
		}
	}
	return sp
}

func (sp *SlavePool) Open() {
	for _, s := range sp.slaves {
		if s != nil {
			s.Open()
		}
	}
}

func (sp *SlavePool) Slaves() int {
	return len(sp.slaves)
}

func (sp *SlavePool) SendWork(job interface{}) {
	v := math.MaxInt32
	sel := 0
	for i, s := range sp.slaves {
		if len := s.ToDo(); len < v {
			v, sel = len, i
		}
	}
	sp.slaves[sel].SendWork(job)
}

func (sp *SlavePool) Close() {
	for _, s := range sp.slaves {
		s.Close()
	}
}
