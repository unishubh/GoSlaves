package slaves

import (
	"time"
)

type work struct {
	work      func(interface{}) interface{}
	afterWork func(interface{})
}

type slave struct {
	ready int32
	exit  chan struct{}
	work  *work
	jobs  Jobs
	Owner *SlavePool
	Type  string
}

// Open Starts the slave creating goroutine
// that waits job notification
func (s *slave) Open() error {
	if s.work == nil {
		return errworkIsNil
	}
	s.ready = 1
	s.exit = make(chan struct{})

	go func() {
		// Loop until jobChan is closed
		for {
			select {
			case <-s.exit:
				return
			default:
				if data := s.jobs.get(); data != nil {
					ret := s.work.work(data)
					if s.work.afterWork != nil {
						s.work.afterWork(ret)
					}
					s.Owner.wg.Add(-1)
				}
			}

			time.Sleep(time.Millisecond * 10)
		}
	}()

	return nil
}

// Close the slave waiting to finish his tasks
func (s *slave) Close() {
	s.exit <- struct{}{}
}

func (s *slave) GetJobs() int {
	return len(s.jobs.jobs)
}
