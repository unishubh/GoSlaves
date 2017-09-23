package slaves

import (
	"sync"
	"sync/atomic"
)

// SlavePool is the structure of the slave pool
type SlavePool struct {
	mx          sync.RWMutex
	wg          sync.WaitGroup
	running     uint32
	work        work
	jobs        Jobs
	Slaves      []*slave
	readySelect []int32
}

// Check if pool is running
func (sp *SlavePool) isRunning() bool {
	return (atomic.LoadUint32(&sp.running) == 1)
}

// Stablish the running parameter
func (sp *SlavePool) setRunning(set bool) {
	if set {
		atomic.StoreUint32(&sp.running, 1)
	} else {
		atomic.StoreUint32(&sp.running, 0)
	}
}

// MakePool creates a pool of slaves.
// work parameter is the function that will be executed when work is send. Cannot be nil.
// after function is the function that will be executed when work finish. Can be nil.
func MakePool(numSlaves int) (sp *SlavePool) {
	sp = &SlavePool{
		running:     0,
		Slaves:      make([]*slave, numSlaves),
		readySelect: make([]int32, numSlaves),
	}
	return
}

// return the number of slaves
func (sp *SlavePool) GetSlaves() int {
	return len(sp.Slaves)
}

// Redefine readySelect variable
func (sp *SlavePool) redefineSlaves() {
	sp.mx.Lock()
	defer sp.mx.Unlock()

	if len(sp.Slaves) != len(sp.readySelect) {
		sp.readySelect = make([]int32, len(sp.Slaves))
	}

	for i := range sp.Slaves {
		sp.readySelect[i] = sp.Slaves[i].ready
	}
}

// Delete slave from slave array
func (sp *SlavePool) deleteSlave(slave int) {
	sp.Slaves[slave].Close()

	sp.mx.Lock()
	sp.Slaves = append(sp.Slaves[:slave], sp.Slaves[slave+1:]...)
	sp.mx.Unlock()

	sp.redefineSlaves()
}

func (sp *SlavePool) DeleteSlave() {
	sp.deleteSlave(len(sp.Slaves) - 1)
}

func (sp *SlavePool) AddSlave() {
	new := &slave{
		work:  &sp.work,
		Owner: sp,
	}
	new.Open()

	sp.mx.Lock()
	sp.Slaves = append(sp.Slaves, new)
	sp.mx.Unlock()

	sp.redefineSlaves()
}

func (sp *SlavePool) prepareEnv() {
	// caught the slaves in range
	for i := range sp.Slaves {
		sp.Slaves[i] = &slave{
			work:  &sp.work,
			Owner: sp,
		}
		sp.Slaves[i].Open()
	}

	sp.redefineSlaves()
}

// Open the slave pool initialising all slaves
// With specified values. toDo cannot be nil.
// If any slave have been created, the library makes 4 by default
func (sp *SlavePool) Open(
	toDo func(interface{}) interface{},
	after func(interface{}),
) error {
	if sp.isRunning() {
		return errAlreadyRunning
	}
	if toDo == nil {
		return errFuncNil
	}
	if len(sp.Slaves) == 0 {
		sp.Slaves = make([]*slave, 4)
	}
	if len(sp.readySelect) == 0 {
		sp.readySelect = make([]int32, 4)
	}

	// assign work to do
	sp.work = work{
		work:      toDo,
		afterWork: after,
	}
	sp.prepareEnv()

	go func() {
		for {
			switch {
			case !sp.isRunning():
				return // exit if not running
			default:
				for chosen := range sp.readySelect {
					// get the slave who is ready
					if atomic.LoadInt32(&sp.Slaves[chosen].ready) == 1 {
						job := sp.jobs.Get()
						if job == nil {
							break
						}

						sp.Slaves[chosen].jobChan <- job
					}
				}
			}
		}
	}()

	sp.setRunning(true)
	return nil
}

// SendWork receives the work and select
// one unemployed slave in goroutine
func (sp *SlavePool) SendWork(job interface{}) {
	if sp.isRunning() {
		sp.wg.Add(1)
		sp.jobs.Put(job)
	}
}

// Close the pool waiting all slaves and his tasks
func (sp *SlavePool) Close() {
	sp.wg.Wait()

	for i := range sp.Slaves {
		sp.Slaves[i].Close()
	}

	sp.setRunning(false)
}
