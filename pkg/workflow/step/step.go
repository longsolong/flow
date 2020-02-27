package step

// Step ...
type Step struct {
	ID              string
	ExpansionDigest string
}

// NewStep ...
func NewStep(id, expansionDigest string) *Step {
	return &Step{
		ID:              id,
		ExpansionDigest: expansionDigest,
	}
}