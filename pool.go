package slaves

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// SlavePool is the structure of the slave pool
type SlavePool struct {
	running     uint32
	mx          sync.Mutex
	work        work
	wg          sync.WaitGroup
	jobs        []interface{}
	Slaves      []*slave
	readySelect []reflect.SelectCase
}

// Check if pool is running
func (sp *SlavePool) isRunning() bool {
	return (atomic.LoadUint32(&sp.running) == 1)
}

// Stablish the running parameter
func (sp *SlavePool) setRunning(set bool) {
	sp.mx.Lock()
	defer sp.mx.Unlock()

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
		jobs:        make([]interface{}, 0),
		Slaves:      make([]*slave, numSlaves),
		readySelect: make([]reflect.SelectCase, numSlaves),
	}
	return
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
		sp.readySelect = make([]reflect.SelectCase, 4)
	}

	sp.work = work{
		work:      toDo,
		afterWork: after,
	}

	// caught the slaves in range
	for i := range sp.Slaves {
		sp.Slaves[i] = &slave{
			work:  &sp.work,
			Owner: sp,
		}
		sp.Slaves[i].Open()

		sp.readySelect[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(sp.Slaves[i].readyChan),
		}
	}

	sp.setRunning(true)
	return nil
}

// SendWork receives the work and select
// one unemployed slave in goroutine
func (sp *SlavePool) SendWork(job interface{}) {
	sp.wg.Add(1)
	sp.jobs = append(sp.jobs, job)
}

func (sp *SlavePool) Work() error {
	jobs := sp.jobs
	go func() {
		for _, j := range jobs {
			chosen, _, ok := reflect.Select(sp.readySelect)
			if chosen >= 0 && ok {
				sp.Slaves[chosen].jobChan <- j
			}
		}
	}()

	return nil
}

// Close the pool waiting all slaves and his tasks
func (sp *SlavePool) Close() {
	sp.wg.Wait()

	for _, s := range sp.Slaves {
		s.Close()
	}
	sp.setRunning(false)
}
