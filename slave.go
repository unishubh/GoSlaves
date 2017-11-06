package slaves

import "time"

type Slave struct {
	ch        chan interface{}
	lastUsage time.Time
	Owner     *SlavePool
}

func (s *Slave) Open() {
	s.ch = make(chan interface{}, 1)
	go func() {
		for w := range s.ch {
			s.Owner.Work(w)
			s.lastUsage = time.Now()
			s.Owner.pool.Put(s)
		}
	}()
}

func (s *Slave) Close() {
	close(s.ch)
}
