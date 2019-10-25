package state

const (
	STATE_UNKNOWN byte = 0xff

	// Normal states, in order
	STATE_SUCCESS  byte = 0 // completed successfully
	STATE_FAIL     byte = 1 // failed
	STATE_RUNNING  byte = 2 // running
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
