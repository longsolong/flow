package state

// State ...
const (
	StateSuccess  byte = iota // completed successfully
	StateFail                 // failed
	StateRunning              // running
	StateCanceled             // canceled

	StateUnknown byte = 0xff
)

// StateText ...
var StateText = map[byte]string{
	StateUnknown:  "UNKNOWN",
	StateRunning:  "RUNNING",
	StateSuccess:  "SUCCESS",
	StateFail:     "FAIL",
	StateCanceled: "CANCELED",
}

// StateValue ...
var StateValue = map[string]byte{
	"UNKNOWN":  StateUnknown,
	"RUNNING":  StateRunning,
	"SUCCESS":  StateSuccess,
	"FAIL":     StateFail,
	"CANCELED": StateCanceled,
}
