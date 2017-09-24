package slaves

import (
	"sync"
)

type iwork struct {
	typed []byte
	job   interface{}
}

type Jobs struct {
	mx   sync.RWMutex
	jobs []iwork
}

// Put work in jobs stack
func (jobs *Jobs) put(job iwork) {
	jobs.mx.Lock()
	jobs.jobs = append(jobs.jobs, job)
	jobs.mx.Unlock()
}

// Get the first job and deletes from stack
func (jobs *Jobs) get() *iwork {
	jobs.mx.Lock()
	defer jobs.mx.Unlock()

	if len(jobs.jobs) > 0 {
		ret := jobs.jobs[0]
		jobs.jobs = jobs.jobs[1:]

		return &ret
	}

	return nil
}
