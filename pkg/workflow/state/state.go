package state

const (
	STATE_UNKNOWN byte = 0

	// Normal states, in order
	STATE_RUNNING  byte = 1 // running
	STATE_SUCCESS  byte = 2 // completed successfully
	STATE_FAIL     byte = 3 // failed
)

var StateName = map[byte]string{
	STATE_UNKNOWN:   "UNKNOWN",
	STATE_RUNNING:   "RUNNING",
	STATE_SUCCESS:   "SUCCESS",
	STATE_FAIL:      "FAIL",
}

var StateValue = map[string]byte{
	"UNKNOWN":   STATE_UNKNOWN,
	"RUNNING":   STATE_RUNNING,
	"SUCCESS":   STATE_SUCCESS,
	"FAIL":      STATE_FAIL,
}
