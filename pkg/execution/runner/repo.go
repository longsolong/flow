// Copyright 2017-2019, Square, Inc.

package runner

import (
	"github.com/orcaman/concurrent-map"
)

// Repo is a small wrapper around a concurrent map that provides the ability to
// store and retrieve Runners in a thread-safe way.
type Repo interface {
	Set(jobID string, runner Runner)
	Get(jobID string) Runner
	Remove(jobID string)
	Items() map[string]Runner
	Count() int
}

type repo struct {
	c cmap.ConcurrentMap
}

// NewRepo ...
func NewRepo() Repo {
	return &repo{
		c: cmap.New(),
	}
}

// Set sets a Runner in the repo.
func (r *repo) Set(jobID string, runner Runner) {
	r.c.Set(jobID, runner)
}

func (r *repo) Get(jobID string) Runner {
	v, ok := r.c.Get(jobID)
	if !ok {
		return nil
	}
	return v.(Runner)
}

// Remove removes a runner from the repo.
func (r *repo) Remove(jobID string) {
	r.c.Remove(jobID)
}

// Items returns a map of jobID => Runner with all the Runners in the repo.
func (r *repo) Items() map[string]Runner {
	runners := map[string]Runner{} // jobID => runner
	for jobID, v := range r.c.Items() {
		runner, ok := v.(Runner)
		if !ok {
			panic("runner for job ID " + jobID + " is not type Runner") // should be impossible
		}
		runners[jobID] = runner
	}
	return runners
}

// Count returns the number of Runners in the repo.
func (r *repo) Count() int {
	return r.c.Count()
}
