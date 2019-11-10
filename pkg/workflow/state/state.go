package state

// State ...
const (
	StateSuccess     byte = iota // completed successfully
	StateFail                    // failed
	StateRunning                 // running
	StateCanceled                // canceled
	StateWaitInput               // wait input
	StateMarkSkipped             // mark as skipped by user
	StateIgnored                 // ignored due to conditional
	StateUpForRetry              // up for retry
	StateMarkRetry               // mark as retry by user

	StateUnknown byte = 0xff
)

// StateText ...
var StateText = map[byte]string{
	StateUnknown:     "UNKNOWN",
	StateRunning:     "RUNNING",
	StateSuccess:     "SUCCESS",
	StateFail:        "FAIL",
	StateCanceled:    "CANCELED",
	StateWaitInput:   "WAIT_INPUT",
	StateMarkSkipped: "MARK_SKIPPED",
	StateIgnored:     "IGNORED",
	StateUpForRetry:  "UP_FOR_RETRY",
	StateMarkRetry:   "MARK_RETRY",
}

// StateValue ...
var StateValue = map[string]byte{
	"UNKNOWN":      StateUnknown,
	"RUNNING":      StateRunning,
	"SUCCESS":      StateSuccess,
	"FAIL":         StateFail,
	"CANCELED":     StateCanceled,
	"WAIT_INPUT":   StateWaitInput,
	"MARK_SKIPPED": StateMarkSkipped,
	"IGNORED":      StateIgnored,
	"UP_FOR_RETRY": StateUpForRetry,
	"MARK_RETRY":   StateMarkRetry,
}
