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

	jobs          map[atom.ID]*job.Job
	jobsMux       *sync.RWMutex    // for access to jobs maps
	triesMux      *sync.RWMutex    // for access to sequence/job tries maps
	sequenceTries map[atom.ID]uint // Number of sequence retries attempted so far
	totalJobTries map[atom.ID]uint // Number of job retries attempted so far
}

// NewChain ...
func NewChain(d *dag.DAG) *Chain {
	return &Chain{
		DAG:           d,
		jobs:          make(map[atom.ID]*job.Job),
		jobsMux:       &sync.RWMutex{},
		triesMux:      &sync.RWMutex{},
		sequenceTries: make(map[atom.ID]uint),
		totalJobTries: make(map[atom.ID]uint),
	}
}

// JobState returns the state of a given job.
func (c *Chain) JobState(atomID atom.ID) state.State {
	c.jobsMux.RLock()
	defer c.jobsMux.RUnlock()
	return c.jobs[atomID].State
}

// SetJobState set the state of a job in the chain.
func (c *Chain) SetJobState(atomID atom.ID, state state.State) {
	c.jobsMux.Lock()
	defer c.jobsMux.Unlock()
	c.jobs[atomID].State = state
}

// AddJob ...
func (c *Chain) AddJob(j *job.Job) {
	c.jobsMux.Lock()
	defer c.jobsMux.Unlock()
	c.jobs[j.ID()] = j
}

// NextJobs finds all of the jobs adjacent to the given job.
func (c *Chain) NextJobs(jobID atom.ID) []*job.Job {
	c.jobsMux.RLock()
	defer c.jobsMux.RUnlock()
	var nextJobs []*job.Job
	var node *dag.Node
	var ok bool
	if node, ok = c.DAG.Vertices[jobID]; !ok {
		panic(fmt.Sprintf("jobID %v not found in Vertices", jobID))
	}

	for nextJobID := range node.Next {
		j := c.jobs[nextJobID]
		nextJobs = append(nextJobs, j)
	}

	return nextJobs
}

// IsRunnable ...
func (c *Chain) IsRunnable(jobID atom.ID) bool {
	c.DAG.VerticesMux.RLock()
	defer c.DAG.VerticesMux.RUnlock()
	return c.isRunnable(jobID)
}

// RunnableJobs ...
func (c *Chain) RunnableJobs() (runnableJobs []*job.Job) {
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
func (c *Chain) isRunnable(jobID atom.ID) bool {
	// CALLER MUST LOCK c.DAG.VerticesMux!
	var node *dag.Node
	var j *job.Job
	var ok bool
	if node, ok = c.DAG.Vertices[jobID]; !ok {
		panic(fmt.Sprintf("jobID %v not found in Vertices", jobID))
	}
	if j, ok = c.jobs[jobID]; !ok {
		panic(fmt.Sprintf("jobID %v not found in jobs", jobID))
	}
	if j.State != state.StateUnknown && j.State != state.StateUpForRetry && j.State != state.StateMarkRetry {
		return false
	}
	// Check that all previous jobs are complete.
	for prevJobID := range node.Prev {
		j := c.jobs[prevJobID]
		if _, ok := state.JobCompleteState[j.State]; !ok {
			return false
		}
	}
	return true
}

// SequenceStartJob ...
func (c *Chain) SequenceStartJob(jobID atom.ID) *job.Job {
	c.jobsMux.RLock()
	defer c.jobsMux.RUnlock()
	return c.jobs[c.DAG.Vertices[jobID].SequenceID]
}

// IsSequenceStartJob ...
func (c *Chain) IsSequenceStartJob(jobID atom.ID) bool {
	c.jobsMux.RLock()
	defer c.jobsMux.RUnlock()
	return jobID == c.DAG.Vertices[jobID].SequenceID
}


// CanRetrySequence ...
func (c *Chain) CanRetrySequence(jobID atom.ID) bool {
	sequenceStartJob := c.SequenceStartJob(jobID)
	c.triesMux.RLock()
	defer c.triesMux.RUnlock()
	return c.sequenceTries[sequenceStartJob.ID()] <= c.DAG.Vertices[sequenceStartJob.ID()].SequenceRetry
}

// IncrementJobTries ...
func (c *Chain) IncrementJobTries(jobID atom.ID, delta uint) {
	c.triesMux.Lock()
	defer c.triesMux.Unlock()
	// Total job tries can only increase. This is the job try count
	// that's monotonically increasing across all sequence retries.
	c.totalJobTries[jobID] += delta
}

// JobTries ...
func (c *Chain) JobTries(jobID atom.ID) uint {
	c.triesMux.RLock()
	defer c.triesMux.RUnlock()
	return c.totalJobTries[jobID]
}

// IncrementSequenceTries ...
func (c *Chain) IncrementSequenceTries(jobID atom.ID, delta uint) {
	seqID := c.DAG.Vertices[jobID].SequenceID
	c.triesMux.Lock()
	c.sequenceTries[seqID] += delta
	c.triesMux.Unlock()
}

// SequenceTries ...
func (c *Chain) SequenceTries(jobID atom.ID) uint {
	seqID := c.DAG.Vertices[jobID].SequenceID
	c.triesMux.RLock()
	defer c.triesMux.RUnlock()
	return c.sequenceTries[seqID]
}
