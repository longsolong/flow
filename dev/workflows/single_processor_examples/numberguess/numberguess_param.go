package numberguess

type numberGuessParam struct {
	serect, low, high int
}

func (p *numberGuessParam) MustSetSecret(secret int) {
	p.serect = secret
}

func (p *numberGuessParam) Secret() int {
	return p.serect
}

func (p *numberGuessParam) MustSetLow(low int) {
	p.low = low
}

func (p *numberGuessParam) Low() int {
	return p.low
}

func (p *numberGuessParam) MustSetHigh(high int) {
	p.high = high
}

func (p *numberGuessParam) High() int {
	return p.high
}
