package analyze

type Probe struct {
	Buffer     []float32
	memory     float32
	index      int
	downsample int
}

func (p *Probe) Analyze(signal []float32) {
	index := p.index
	downsample := p.downsample
	memory := p.memory
	weight := 1 - memory
	for _, sample := range signal {
		bufferIndex := index >> downsample
		p.Buffer[bufferIndex] = memory*p.Buffer[bufferIndex] + weight*sample

		index += 1
		bufferIndex = index >> downsample
		if bufferIndex == len(p.Buffer) {
			index = 0
		}
	}
	p.index = index
}
