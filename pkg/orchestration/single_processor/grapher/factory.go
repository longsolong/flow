package grapher

// A Factory makes Graphers.
type Factory interface {
	// Make makes a Grapher. A new grapher should be made for every request.
	Make(name string, version int, rawRequestArgs []byte) (*Grapher, error)
}
