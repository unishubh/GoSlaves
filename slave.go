package slaves

import "sync"

type work struct {
	work      func(interface{}) interface{}
	afterWork func(interface{})
}

type slave struct {
	readyChan chan struct{}
	jobChan   chan interface{}
	mx        sync.Mutex
	Owner     *SlavePool
	work      *work
	Type      []byte
}

// Open Starts the slave creating goroutine
// that waits job notification
func (s *slave) Open() error {
	if s.work == nil {
		return errworkIsNil
	}
	s.readyChan = make(chan struct{})
	s.jobChan = make(chan interface{})

	go func() {
		// Slave is ready to receive a job
		s.readyChan <- struct{}{}
		// Loop until jobChan is closed
		for data := range s.jobChan {
			ret := s.work.work(data)
			if s.work.afterWork != nil {
				s.work.afterWork(ret)
			}

			s.Owner.wg.Done()
			// notify slave is ready to work
			s.readyChan <- struct{}{}
		}

		close(s.readyChan)
	}()

	return nil
}

// SetWork sets new Work for slave.
// If toDo is nil, the parameter is ignored
// it's not the same with afterWork value, because this is not important
func (s *slave) SetWork(
	toDo func(interface{}) interface{},
	afterWork func(interface{}),
) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.work == nil {
		s.work = new(work)
	}

	if toDo != nil {
		s.work.work = toDo
	}
	s.work.afterWork = afterWork
}

// Close the slave waiting to finish his tasks
func (s *slave) Close() {
	close(s.jobChan)
}
