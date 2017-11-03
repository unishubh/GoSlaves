package slaves

import (
	"errors"
	"math"
	"sync"
)

var (
	ErrStacked = errors.New("error stacked: All slaves are busy")
)

// SlavePool object
type SlavePool struct {
	mx     sync.Mutex
	work   func(interface{}) interface{}
	after  func(interface{})
	wg     sync.WaitGroup
	Slaves []*Slave
	// NotStack default is false
	NotStack bool
}

// MakePool creates a pool and initialise Slaves
// num is the number of Slaves
// work is the function to execute in the Slaves
// after the function to execute after execution of work
func MakePool(num uint, work func(interface{}) interface{},
	after func(interface{})) SlavePool {

	if work == nil {
		work = defaultWork
	}
	if after == nil {
		after = defaultAfter
	}
	sp := SlavePool{
		Slaves: make([]*Slave, num),
		work:   work,
		after:  after,
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

// Working returns the number of slaves
// who are working
func (sp *SlavePool) Working() int {
	var c int = 0
	for _, s := range sp.Slaves {
		if s.ToDo() > 0 {
			c++
		}
	}
	return c
}

// Add adds slave to the pool
func (sp *SlavePool) Add(s Slave) *Slave {
	slave := NewSlave(s.Name, s.Work, s.After)
	sp.mx.Lock()
	sp.Slaves = append(sp.Slaves, slave)
	sp.mx.Unlock()
	return slave
}

// Del deletes the last slave
func (sp *SlavePool) Del() {
	sp.mx.Lock()
	sp.Slaves = sp.Slaves[:len(sp.Slaves)-1]
	sp.mx.Unlock()
}

// AddUSend creates new slave and adds into
// the new queue sends job
func (sp *SlavePool) AddUSend(s Slave, job interface{}) {
	sp.Add(s).SendWork(job)
}

// findNonStacked search non-busy slaves
func (sp *SlavePool) findNonStacked() int {
	for i, s := range sp.Slaves {
		if s.ToDo() == 0 {
			return i
		}
	}
	return -1
}

// SendWork Send work to the pool.
// This function get the slave with less number
// of works and send him the job
func (sp *SlavePool) SendWork(job interface{}) error {
	if sp.NotStack {
		if s := sp.findNonStacked(); s > 0 {
			sp.Slaves[s].SendWork(job)
			return nil
		}
		return ErrStacked
	}

	v := math.MaxInt32
	sel := 0
	for i, s := range sp.Slaves {
		if s != nil {
			if len := s.ToDo(); len < v {
				v, sel = len, i
			}
		}
	}
	sp.Slaves[sel].SendWork(job)
	return nil
}

// SendWorkTo send work to specified worker
func (sp *SlavePool) SendWorkTo(name string, job interface{}) {
	v := math.MaxInt32
	sel := 0
	for i, s := range sp.Slaves {
		if s != nil {
			if len := s.ToDo(); len < v && name == s.Name {
				v, sel = len, i
			}
		}
	}
	sp.Slaves[sel].SendWork(job)
}

// Send executes work in new goroutine
// without adding into slave pool
func (sp *SlavePool) Send(job interface{}) {
	sp.wg.Add(1)
	go func() {
		sp.after(sp.work(job))
		sp.wg.Done()
	}()
}

// Close closes the pool waiting
// the end of all jobs
func (sp *SlavePool) Close() {
	sp.wg.Wait()
	for _, s := range sp.Slaves {
		s.Close()
	}
}

// Force close closes all pool slaves
// without waiting slave job end
func (sp *SlavePool) ForceClose() {
	for i := sp.Len() - 1; i > 0; i = sp.Len() - 1 {
		sp.Slaves[i].Close()
		sp.Del()
	}
}
