package analyze

type Analyzer struct {
	Probes [][2]Probe
}

func NewAnalyzer() *Analyzer {
	analyzer := Analyzer{}
	memory := float32(0.99)

	for i := 10; i < 48; i++ {
		analyzer.Probes = append(analyzer.Probes, [2]Probe{
			{
				Buffer: make([]float32, i),
				memory: memory,
			},
			{
				Buffer: make([]float32, i),
				memory: memory,
			},
		})
	}
	for downsample := 1; downsample < 9; downsample++ {
		for i := 24; i < 48; i++ {
			analyzer.Probes = append(analyzer.Probes, [2]Probe{
				{
					Buffer:     make([]float32, i),
					memory:     memory,
					downsample: downsample,
				},
				{
					Buffer:     make([]float32, i),
					memory:     memory,
					downsample: downsample,
				},
			})
		}
	}
	return &analyzer
}

func (a *Analyzer) Analyze(leftBuffer, rightBuffer []float32) {
	for _, probePair := range a.Probes {
		probePair[0].Analyze(leftBuffer)
		probePair[1].Analyze(rightBuffer)
	}
}
