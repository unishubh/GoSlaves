package slaves

import (
	"sync"
	"sync/atomic"
	"time"
)

// Queue allows programmer to stack tasks.
type Queue struct {
	closed  bool
	locker  sync.RWMutex
	jobs    []interface{}
	slaves  []*slave
	ch      chan struct{}
	max     int
	current int32
	wg      sync.WaitGroup
}

// max is the maximum goroutines to execute.
// 0 is the same as no limit
func DoQueue(max int, work func(obj interface{})) *Queue {
	ch := make(chan struct{}, 1)

	queue := &Queue{
		max:    max,
		ch:     ch,
		jobs:   make([]interface{}, 0),
		slaves: make([]*slave, 0),
	}

	go cleanSlaves(queue.locker, &queue.slaves)

	queue.wg.Add(1)
	go func() {
		defer queue.wg.Done()
		for !queue.closed {
			time.Sleep(time.Millisecond * 20)
			queue.locker.Lock()
			if len(queue.jobs) > 0 {
				ch <- struct{}{}
			}
			queue.locker.Unlock()
		}
	}()

	go func() {
		// selected slave
		var s *slave
		var c interface{}
		m := int32(max)
		for range ch {
			for _, c = range queue.jobs {
				queue.locker.Lock()
				// getting job to do
				if len(queue.jobs) > 1 {
					queue.jobs = queue.jobs[1:]
				} else {
					queue.jobs = queue.jobs[:0]
				}
				if len(queue.slaves) == 0 {
					if atomic.LoadInt32(&queue.current) >= m {
						queue.jobs = append(queue.jobs, c)
						continue
					} else {
						s = &slave{
							ch:        make(chan interface{}, 1),
							lastUsage: time.Now(),
						}
						go func(sv *slave) {
							atomic.AddInt32(&queue.current, 1)
							defer atomic.AddInt32(&queue.current, -1)
							var w interface{}
							for w = range sv.ch {
								if w == nil {
									sv.close()
									return
								}
								work(w)
								sv.lastUsage = time.Now()
							}
						}(s)
					}
				} else {
					s = queue.slaves[0]
					queue.slaves = queue.slaves[1:]
				}
				queue.locker.Unlock()
				// parsing job
				s.ch <- c
			}
		}
	}()

	return queue
}

func (queue *Queue) Serve(job interface{}) {
	queue.locker.Lock()
	queue.jobs = append(queue.jobs, job)
	queue.locker.Unlock()
}

func (queue *Queue) WaitClose() {
	queue.locker.Lock()
	queue.closed = true
	queue.locker.Unlock()

	queue.wg.Wait()
	queue.Close()
}

func (queue *Queue) Close() {
	queue.closed = true
	close(queue.ch)

	for _, s := range queue.slaves {
		s.close()
	}
}
