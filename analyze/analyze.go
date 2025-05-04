package analyze

type Analyzer struct {
	LeftProbes, RightProbes []Probe
	downsampleBuffer        []float32
}

func NewAnalyzer() *Analyzer {
	analyzer := Analyzer{}

	for i := 10; i < 48; i++ {
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
	nSamples := len(leftBuffer) // always the same as the right buffer
	downsampleBuffer := a.downsampleBuffer

	// TODO: the downsampling won't work unless the buffer size is always
	// a power of two, which it won't be because it is measured in milliseconds.
	// So I need a more sophisticated strategy for downsampling.
	// Apply a simple window filter to downsample by half repeatedly.
	if len(downsampleBuffer) < nSamples {
		a.downsampleBuffer = make([]float32, nSamples)
		downsampleBuffer = a.downsampleBuffer
	}
	for i := 0; i < nSamples; i += 2 {
		downsampleBuffer[i>>1] = (leftBuffer[i] + leftBuffer[i+1]) / 2
	}
	for i := 0; i < nSamples; i += 2 {
		downsampleBuffer[(nSamples+i)>>1] = (downsampleBuffer[i] + downsampleBuffer[i+1]) / 2
	}

	for _, probe := range a.LeftProbes {
		if probe.Downsample == 0 {
			probe.Analyze(leftBuffer)
		} else { // If there aren't enough downsamples, whatever. Length will be 0.
			length := nSamples >> probe.Downsample
			end := nSamples - length
			start := end - length
			probe.Analyze(downsampleBuffer[start:end])
		}
	}

	// Apply a simple window filter to downsample by half repeatedly.
	if len(downsampleBuffer) < nSamples {
		a.downsampleBuffer = make([]float32, nSamples)
		downsampleBuffer = a.downsampleBuffer
	}
	for i := 0; i < nSamples; i += 2 {
		downsampleBuffer[i>>1] = (rightBuffer[i] + rightBuffer[i+1]) / 2
	}
	for i := 0; i < nSamples; i += 2 {
		downsampleBuffer[(nSamples+i)>>1] = (downsampleBuffer[i] + downsampleBuffer[i+1]) / 2
	}

	for _, probe := range a.RightProbes {
		if probe.Downsample == 0 {
			probe.Analyze(rightBuffer)
		} else { // If there aren't enough downsamples, whatever. Length will be 0.
			length := nSamples >> probe.Downsample
			end := nSamples - length
			start := end - length
			probe.Analyze(downsampleBuffer[start:end])
		}
	}
}
