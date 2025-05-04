package analyze

type Analyzer struct {
	pos uint64
}

func (a *Analyzer) Analyze(leftBuffer, rightBuffer []float32) {
	a.pos += uint64(len(leftBuffer))
}
