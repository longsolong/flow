package state

// State ...
type State byte

// State const ...
const (
	StateSuccess     State = iota // completed successfully
	StateFail                     // failed
	StateException                // exception
	StateStopped                  // stopped
	StateRunning                  // running
	StateCanceled                 // canceled
	StateWaitInput                // wait input
	StateMarkSkipped              // mark as skipped by user
	StateIgnored                  // ignored due to conditional
	StateUpForRetry               // up for retry
	StateMarkRetry                // mark as retry by user

	StateUnknown State = 0xff
)

var (
	// JobState ...
	JobState = map[State]bool{
		StateSuccess:     true,
		StateFail:        true,
		StateException:   true,
		StateStopped:     true,
		StateRunning:     true,
		StateCanceled:    true,
		StateWaitInput:   true,
		StateMarkSkipped: true,
		StateIgnored:     true,
		StateUpForRetry:  true,
		StateMarkRetry:   true,
		StateUnknown:     true,
	}
	// JobUndoneState ...
	JobUndoneState = map[State]bool{
		StateRunning:    true,
		StateWaitInput:  true,
		StateUpForRetry: true,
		StateMarkRetry:  true,
		StateUnknown:    true,
	}
	// JobDoneState ...
	JobDoneState = map[State]bool{
		StateSuccess:     true,
		StateFail:        true,
		StateStopped:     true,
		StateException:   true,
		StateCanceled:    true,
		StateMarkSkipped: true,
		StateIgnored:     true,
	}
	// JobCompleteState ...
	JobCompleteState = map[State]bool{
		StateSuccess:     true,
		StateMarkSkipped: true,
		StateIgnored:     true,
	}
)

// StateText ...
var StateText = map[State]string{
	StateUnknown:     "UNKNOWN",
	StateRunning:     "RUNNING",
	StateSuccess:     "SUCCESS",
	StateFail:        "FAIL",
	StateException:   "EXCEPTION",
	StateStopped:     "STOPPED",
	StateCanceled:    "CANCELED",
	StateWaitInput:   "WAIT_INPUT",
	StateMarkSkipped: "MARK_SKIPPED",
	StateIgnored:     "IGNORED",
	StateUpForRetry:  "UP_FOR_RETRY",
	StateMarkRetry:   "MARK_RETRY",
}

// StateValue ...
var StateValue = map[string]State{
	"UNKNOWN":      StateUnknown,
	"RUNNING":      StateRunning,
	"SUCCESS":      StateSuccess,
	"FAIL":         StateFail,
	"EXCEPTION":    StateException,
	"STOPPED":      StateStopped,
	"CANCELED":     StateCanceled,
	"WAIT_INPUT":   StateWaitInput,
	"MARK_SKIPPED": StateMarkSkipped,
	"IGNORED":      StateIgnored,
	"UP_FOR_RETRY": StateUpForRetry,
	"MARK_RETRY":   StateMarkRetry,
}
