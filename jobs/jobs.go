package jobs

import (
	"errors"
	"github.com/eapache/channels"
)

// Handle multiple jobs
// enqueuing in buffered channel
type Jobs struct {
	ch *channels.InfiniteChannel
}

// Open jobs channel
func (jobs *Jobs) Open() {
	jobs.ch = channels.NewInfiniteChannel()
}

// Parse job to channel
func (jobs *Jobs) Put(job interface{}) {
	jobs.ch.In() <- job
}

// Get the length of jobs to do
func (jobs *Jobs) Len() int {
	return jobs.ch.Len()
}

// Gets a job from the buffered channel
// if error is returned Close() function have
// been called
func (jobs *Jobs) Get() (interface{}, error) {
	r, ok := <-jobs.ch.Out()
	if !ok {
		return nil, errors.New("chan closed")
	}
	return r, nil
}

// stop to receive jobs
func (jobs *Jobs) Close() {
	jobs.ch.Close()
}
