package slaves

import (
	"github.com/themester/GoSlaves/slave"
	"math"
)

// SlavePool object
type SlavePool struct {
	slaves []*slave.Slave
}

// MakePool creates a pool and initialise slaves
// num is the number of slaves
// work is the function to execute in the slaves
// after the function to execute after execution of work
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

// Open all slaves
func (sp *SlavePool) Open() {
	for _, s := range sp.slaves {
		if s != nil {
			s.Open()
		}
	}
}

// Get the length of the slave array
func (sp *SlavePool) Slaves() int {
	return len(sp.slaves)
}

// Send work to the pool.
// This function get the slave with less number
// of works and send him the job
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

// Close the pool waiting the end
// of all jobs
func (sp *SlavePool) Close() {
	for _, s := range sp.slaves {
		s.Close()
	}
}
