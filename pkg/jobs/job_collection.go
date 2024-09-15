package jobs

import (
	"context"
	"fmt"
	"sync"
)

type JobCollection struct {
	mux  sync.Mutex
	jobs map[string]*PlotJob
}

func NewJobCollection() *JobCollection {
	return &JobCollection{
		jobs: make(map[string]*PlotJob),
	}
}

func (jc *JobCollection) Add(id, path string) {
	jc.mux.Lock()
	defer jc.mux.Unlock()
	jc.jobs[id] = &PlotJob{
		ID:    id,
		Path:  path,
		State: JobStateNew,
	}
}

func (jc *JobCollection) Remove(id string) {
	jc.mux.Lock()
	defer jc.mux.Unlock()
	delete(jc.jobs, id)
}

func (jc *JobCollection) Get(id string) (*PlotJob, error) {
	jc.mux.Lock()
	defer jc.mux.Unlock()
	job, ok := jc.jobs[id]
	if !ok {
		return nil, fmt.Errorf("job not found: %s", id)
	}
	return job, nil
}

func (jc *JobCollection) UpdateState(id string, state JobState) error {
	jc.mux.Lock()
	defer jc.mux.Unlock()
	job, ok := jc.jobs[id]
	if !ok {
		return fmt.Errorf("job not found: %s", id)
	}
	job.State = state
	return nil
}

func (jc *JobCollection) SetCancel(id string, cancel context.CancelFunc) error {
	jc.mux.Lock()
	defer jc.mux.Unlock()
	job, ok := jc.jobs[id]
	if !ok {
		return fmt.Errorf("job not found: %s", id)
	}
	job.Cancel = &cancel
	return nil
}
