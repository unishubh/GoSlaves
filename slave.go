package slaves

type Slave struct {
	ch    chan interface{}
	Owner *SlavePool
}

func (s *Slave) Open() {
	s.ch = make(chan interface{})
	go s.do()
}

func (s *Slave) do() {
	for w := range s.ch {
		s.Owner.Work(w)
		s.Owner.ready = append(
			s.Owner.ready, s,
		)
	}
}

func (s *Slave) Close() {
	close(s.ch)
}
