package slaves

import (
	"sync"
)

type Jobs struct {
	mx   sync.RWMutex
	jobs []interface{}
}

// Put work in jobs stack
func (jobs *Jobs) Put(job interface{}) {
	jobs.mx.Lock()
	jobs.jobs = append(jobs.jobs, job)
	jobs.mx.Unlock()
}

// Get the first job and deletes from stack
func (jobs *Jobs) Get() interface{} {
	jobs.mx.Lock()
	defer jobs.mx.Unlock()

	if len(jobs.jobs) > 0 {
		ret := jobs.jobs[0]
		jobs.jobs = jobs.jobs[1:]

		return ret
	}

	return nil
}
