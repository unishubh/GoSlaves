package slave

import (
	"errors"
	"github.com/themester/GoSlaves/jobs"
	"sync"
)

var (
	defaultAfter   = func(obj interface{}) {}
	ErrSlaveOpened = errors.New("slave is already opened")
)

type Slave struct {
	jobs *jobs.Jobs
	// Name of slave
	Name string
	// Work of slave
	Work func(interface{}) interface{}
	// Function that will be execute when
	// Work finishes. The return value of
	// Work() will be parse to After()
	After func(interface{})
	wg    sync.WaitGroup
}

func NewSlave(name string,
	work func(interface{}) interface{},
	after func(interface{})) *Slave {

	return &Slave{
		Work:  work,
		After: after,
	}
}

func (s *Slave) Open() error {
	if s.jobs != nil {
		return ErrSlaveOpened
	}
	s.jobs = new(jobs.Jobs)
	s.jobs.Open()

	if s.After == nil {
		s.After = defaultAfter
	}

	go func() {
		for {
			job, err := s.jobs.Get()
			if err != nil {
				return
			}
			s.After(s.Work(job))
			s.wg.Done()
		}
	}()

	return nil
}

func (s *Slave) SendWork(job interface{}) {
	s.wg.Add(1)
	s.jobs.Put(job)
}

func (s *Slave) ToDo() int {
	return s.jobs.Len()
}

func (s *Slave) Close() {
	s.wg.Wait()
	s.jobs.Close()
}
