package slaves

import "time"

type Slave struct {
	ch        chan interface{}
	lastUsage time.Time
}

func (s *Slave) Open() {
	s.ch = make(chan interface{}, 1)
}

func (s *Slave) Close() {
	close(s.ch)
}
