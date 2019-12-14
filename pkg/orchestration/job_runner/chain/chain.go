package chain

import (
	"fmt"
	"github.com/longsolong/flow/pkg/orchestration/job"
	"github.com/longsolong/flow/pkg/workflow/atom"
	"github.com/longsolong/flow/pkg/workflow/dag"
	"github.com/longsolong/flow/pkg/workflow/state"
	"sync"
)

// JobChain represents a job chain and some meta information about it.
type JobChain struct {
	*dag.DAG

	jobs          map[atom.ID]*job.Job
	jobsMux       *sync.RWMutex    // for access to jobs maps
	triesMux      *sync.RWMutex    // for access to sequence/job tries maps
	sequenceTries map[atom.ID]uint // Number of sequence retries attempted so far
	jobTries      map[atom.ID]uint // Number of job retries attempted so far
}

// NewJobChain ...
func NewJobChain(d *dag.DAG) *JobChain {
	return &JobChain{
		DAG:           d,
		jobs:          make(map[atom.ID]*job.Job),
		jobsMux:       &sync.RWMutex{},
		triesMux:      &sync.RWMutex{},
		sequenceTries: make(map[atom.ID]uint),
		jobTries:      make(map[atom.ID]uint),
	}
}

// AddJob ...
func (c *JobChain) AddJob(j *job.Job) {
	c.jobsMux.Lock()
	defer c.jobsMux.Unlock()
	c.jobs[j.ID()] = j
}

// IsRunnable ...
func (c *JobChain) IsRunnable(jobID atom.ID) bool {
	c.DAG.VerticesMux.RLock()
	defer c.DAG.VerticesMux.RUnlock()
	return c.isRunnable(jobID)
}

// RunnableJobs ...
func (c *JobChain) RunnableJobs() (runnableJobs []*job.Job) {
	for jobID, j := range c.jobs {
		if !c.IsRunnable(jobID) {
			continue
		}
		runnableJobs = append(runnableJobs, j)
	}
	return runnableJobs
}

// IsDoneRunning returns two booleans: done indicates if there are running or
// runnable jobs, and complete indicates if all jobs finished successfully
// (StateSuccess, StateMarkSkipped, StateIgnored).
//
// A chain is complete if every job finished successfully (StateSuccess, StateMarkSkipped, StateIgnored).
//
// A chain is done running if there are no running or runnable jobs.
// The reaper waits for running jobs to reap them. Reapers roll back failed jobs
// if the sequence can be retried. Consequently, failed jobs do not mean the chain
// is done, and they do not immediately fail the whole chain.
//
// For chain A -> B -> C, if B is stopped, C is not runnable; the chain is done.
// But add job D off A (A -> D) and although B is stopped, if D is pending then
// the chain is not done. This is a side-effect of not stopping/failing
// the whole chain when a job stops/fails. Instead, the chain continues to run
// independent sequences.
func (c *JobChain) IsDoneRunning() (done bool, complete bool) {
	c.DAG.VerticesMux.RLock()
	defer c.DAG.VerticesMux.RUnlock()
	complete = true
	for _, j := range c.jobs {
		if _, ok := state.JobCompleteState[j.State]; ok {
			continue
		}
		if j.State == state.StateUnknown {
			if c.isRunnable(j.ID()) {
				return false, false
			}
		} else if _, ok := state.JobUndoneState[j.State]; ok {
			return false, false
		}

		complete = false
	}
	return true, complete
}

// isRunnable returns true if the job is runnable. A job is runnable iff its
// state is StateUnknown || StateUpForRetry || StateMarkRetry and all immediately previous jobs are state COMPLETE.
func (c *JobChain) isRunnable(jobID atom.ID) bool {
	// CALLER MUST LOCK c.DAG.VerticesMux!
	var node *dag.Node
	var j *job.Job
	var ok bool
	if node, ok = c.DAG.Vertices[jobID]; !ok {
		panic(fmt.Sprintf("jobID %v not found", jobID))
	}
	if j, ok = c.jobs[jobID]; !ok {
		panic(fmt.Sprintf("jobID %v not found", jobID))
	}
	if j.State != state.StateUnknown && j.State != state.StateUpForRetry && j.State != state.StateMarkRetry {
		return false
	}
	// Check that all previous jobs are complete.
	for _, prev := range node.Prev {
		j := c.jobs[prev.Datum.ID()]
		if _, ok := state.JobCompleteState[j.State]; !ok {
			return false
		}
	}
	return true
}

