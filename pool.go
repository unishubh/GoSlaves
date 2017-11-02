package slaves

import (
	"math"
	"sync"
)

// SlavePool object
type SlavePool struct {
	mx     sync.Mutex
	Slaves []*Slave
}

// MakePool creates a pool and initialise Slaves
// num is the number of Slaves
// work is the function to execute in the Slaves
// after the function to execute after execution of work
func MakePool(num uint, work func(interface{}) interface{},
	after func(interface{})) SlavePool {

	sp := SlavePool{
		Slaves: make([]*Slave, num),
	}
	for i := range sp.Slaves {
		sp.Slaves[i] = NewSlave("", work, after)
	}
	return sp
}

// Len Gets the length of the slave array
func (sp *SlavePool) Len() int {
	return len(sp.Slaves)
}

// Add slave to the pool
func (sp *SlavePool) Add(s Slave) {
	slave := NewSlave(s.Name, s.Work, s.After)
	sp.mx.Lock()
	sp.Slaves = append(sp.Slaves, slave)
	sp.mx.Unlock()
}

// Delete the last slave
func (sp *SlavePool) Del() {
	sp.mx.Lock()
	sp.Slaves = sp.Slaves[:len(sp.Slaves)-1]
	sp.mx.Unlock()
}

// SendWork Send work to the pool.
// This function get the slave with less number
// of works and send him the job
func (sp *SlavePool) SendWork(job interface{}) {
	v := math.MaxInt32
	sel := 0
	for i, s := range sp.Slaves {
		if len := s.ToDo(); len < v {
			v, sel = len, i
		}
	}
	sp.Slaves[sel].SendWork(job)
}

func (sp *SlavePool) SendWorkTo(name string, job interface{}) {
	v := math.MaxInt32
	sel := 0
	for i, s := range sp.Slaves {
		if len := s.ToDo(); len < v && name == s.Name {
			v, sel = len, i
		}
	}
	sp.Slaves[sel].SendWork(job)
}

// Close closes the pool waiting
// the end of all jobs
func (sp *SlavePool) Close() {
	for _, s := range sp.Slaves {
		s.Close()
	}
}
