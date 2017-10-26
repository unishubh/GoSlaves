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

func (jobs *Jobs) Open() {
	jobs.ch = channels.NewInfiniteChannel()
}

func (jobs *Jobs) Put(job interface{}) {
	jobs.ch.In() <- job
}

func (jobs *Jobs) Len() int {
	return jobs.ch.Len()
}

func (jobs *Jobs) Get() (interface{}, error) {
	r, ok := <-jobs.ch.Out()
	if !ok {
		return nil, errors.New("chan closed")
	}
	return r, nil
}

func (jobs *Jobs) Close() {
	jobs.ch.Close()
}
