package slaves

import (
	"sync/atomic"
)

type work struct {
	work      func(interface{}) interface{}
	afterWork func(interface{})
}

type slave struct {
	ready   int32
	jobChan chan interface{}
	Owner   *SlavePool
	work    *work
	Type    []byte
}

// Open Starts the slave creating goroutine
// that waits job notification
func (s *slave) Open() error {
	if s.work == nil {
		return errworkIsNil
	}
	s.ready = 1
	s.jobChan = make(chan interface{})

	go func() {
		// Loop until jobChan is closed
		for data := range s.jobChan {
			atomic.StoreInt32(&s.ready, 0)

			ret := s.work.work(data)
			if s.work.afterWork != nil {
				s.work.afterWork(ret)
			}
			s.Owner.wg.Add(-1)

			// notify slave is ready to work
			atomic.StoreInt32(&s.ready, 1)
		}
	}()

	return nil
}

// Close the slave waiting to finish his tasks
func (s *slave) Close() {
	close(s.jobChan)
}
