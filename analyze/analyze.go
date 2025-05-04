package analyze

type Analyzer struct {
	LeftProbes, RightProbes []Probe
}

func NewAnalyzer() *Analyzer {
	analyzer := Analyzer{}

	for i := 12; i < 24; i++ {
		analyzer.LeftProbes = append(analyzer.LeftProbes, Probe{
			Buffer: make([]float32, i),
			memory: 0.9,
		})
		analyzer.RightProbes = append(analyzer.LeftProbes, Probe{
			Buffer: make([]float32, i),
			memory: 0.9,
		})
	}

	return &analyzer
}

func (a *Analyzer) Analyze(leftBuffer, rightBuffer []float32) {
	for _, probe := range a.LeftProbes {
		probe.Analyze(leftBuffer)
	}
	for _, probe := range a.RightProbes {
		probe.Analyze(rightBuffer)
	}
}
