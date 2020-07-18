package chain

import (
	"fmt"
	"github.com/longsolong/flow/pkg/orchestration/job"
	"github.com/longsolong/flow/pkg/workflow/atom"
	"github.com/longsolong/flow/pkg/workflow/dag"
	"github.com/longsolong/flow/pkg/workflow/state"
	"sync"
)

// Chain represents a job chain and some meta information about it.
type Chain struct {
	*dag.DAG

	jobs    map[atom.AtomID]*job.Job
	jobsMux *sync.RWMutex // for access to jobs maps

	triesMux      *sync.RWMutex        // for access to sequence/job tries maps
	sequenceTries map[atom.AtomID]uint // Number of sequence retries attempted so far
	totalJobTries map[atom.AtomID]uint // Number of job retries attempted so far
}

// NewChain ...
func NewChain(d *dag.DAG) *Chain {
	return &Chain{
		DAG:           d,
		jobs:          make(map[atom.AtomID]*job.Job),
		jobsMux:       &sync.RWMutex{},
		triesMux:      &sync.RWMutex{},
		sequenceTries: make(map[atom.AtomID]uint),
		totalJobTries: make(map[atom.AtomID]uint),
	}
}

// JobState returns the state of a given job.
func (c *Chain) JobState(atomID atom.AtomID) state.State {
	c.jobsMux.RLock()
	defer c.jobsMux.RUnlock()
	return c.jobs[atomID].State
}

// SetJobState set the state of a job in the chain.
func (c *Chain) SetJobState(atomID atom.AtomID, state state.State) {
	c.jobsMux.Lock()
	defer c.jobsMux.Unlock()
	c.jobs[atomID].State = state
}

// AddJob ...
func (c *Chain) AddJob(j *job.Job) {
	c.jobsMux.Lock()
	defer c.jobsMux.Unlock()
	c.jobs[j.AtomID()] = j
}

// AllJobs ...
func (c *Chain) AllJobs() (allJobs []*job.Job) {
	for _, j := range c.jobs {
		allJobs = append(allJobs, j)
	}
	return allJobs
}

// NextJobs finds all of the jobs adjacent to the given job.
func (c *Chain) NextJobs(jobID atom.AtomID) []*job.Job {
	c.jobsMux.RLock()
	defer c.jobsMux.RUnlock()
	var nextJobs []*job.Job
	node := c.DAG.MustGetNode(jobID)

	for nextJobID := range node.Downstream() {
		j := c.jobs[nextJobID]
		nextJobs = append(nextJobs, j)
	}

	return nextJobs
}

// IsRunnable ...
func (c *Chain) IsRunnable(jobID atom.AtomID) bool {
	c.DAG.VerticesMux.RLock()
	defer c.DAG.VerticesMux.RUnlock()
	return c.isRunnable(jobID)
}

// RunnableJobs ...
func (c *Chain) RunnableJobs() (runnableJobs []*job.Job) {
	c.jobsMux.RLock()
	defer c.jobsMux.RUnlock()
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
func (c *Chain) IsDoneRunning() (done bool, complete bool) {
	c.DAG.VerticesMux.RLock()
	c.jobsMux.RLock()
	defer c.jobsMux.RUnlock()
	defer c.DAG.VerticesMux.RUnlock()
	complete = true
	for _, j := range c.jobs {
		if _, ok := state.JobCompleteState[j.State]; ok {
			continue
		}
		if j.State == state.StateUnknown {
			if c.isRunnable(j.AtomID()) {
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
func (c *Chain) isRunnable(jobID atom.AtomID) bool {
	// CALLER MUST LOCK c.DAG.VerticesMux!
	var j *job.Job
	var ok bool
	if j, ok = c.jobs[jobID]; !ok {
		panic(fmt.Sprintf("jobID %v not found in jobs", jobID))
	}
	if j.State != state.StateUnknown && j.State != state.StateUpForRetry && j.State != state.StateMarkRetry {
		return false
	}
	// Check that all previous jobs are complete.
	node := c.DAG.MustGetNode(jobID)
	for prevJobID := range node.Upstream() {
		j := c.jobs[prevJobID]
		if _, ok := state.JobCompleteState[j.State]; !ok {
			return false
		}
	}
	return true
}

// SequenceStartJob ...
func (c *Chain) SequenceStartJob(jobID atom.AtomID) *job.Job {
	node := c.DAG.MustGetNode(jobID)
	if node.SequenceID.IsEmpty() {
		return nil
	}
	return c.jobs[node.SequenceID]
}

// IsSequenceStartJob ...
func (c *Chain) IsSequenceStartJob(jobID atom.AtomID) bool {
	node := c.DAG.MustGetNode(jobID)
	return jobID == node.SequenceID
}

// CanRetrySequence ...
func (c *Chain) CanRetrySequence(jobID atom.AtomID) bool {
	sequenceStartJob := c.SequenceStartJob(jobID)
	if sequenceStartJob == nil {
		return false
	}
	c.triesMux.RLock()
	defer c.triesMux.RUnlock()
	node := c.DAG.MustGetNode(sequenceStartJob.AtomID())
	return c.sequenceTries[sequenceStartJob.AtomID()] <= node.SequenceRetry
}

// IncrementJobTries ...
func (c *Chain) IncrementJobTries(jobID atom.AtomID, delta uint) {
	c.triesMux.Lock()
	defer c.triesMux.Unlock()
	// Total job tries can only increase. This is the job try count
	// that's monotonically increasing across all sequence retries.
	c.totalJobTries[jobID] += delta
}

// JobTries ...
func (c *Chain) JobTries(jobID atom.AtomID) uint {
	c.triesMux.RLock()
	defer c.triesMux.RUnlock()
	return c.totalJobTries[jobID]
}

// IncrementSequenceTries ...
func (c *Chain) IncrementSequenceTries(jobID atom.AtomID, delta uint) {
	node := c.DAG.MustGetNode(jobID)
	c.triesMux.Lock()
	defer c.triesMux.Unlock()
	c.sequenceTries[node.SequenceID] += delta
}

// SequenceTries ...
func (c *Chain) SequenceTries(jobID atom.AtomID) uint {
	node := c.DAG.MustGetNode(jobID)
	c.triesMux.RLock()
	defer c.triesMux.RUnlock()
	return c.sequenceTries[node.SequenceID]
}
