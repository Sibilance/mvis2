package analyze

type Probe struct {
	Buffer     []float32
	Downsample int
	memory     float32
	index      int
}

func (p *Probe) Analyze(signal []float32) {
	index := p.index
	memory := p.memory
	weight := 1 - memory
	for _, sample := range signal {
		p.Buffer[index] = memory*p.Buffer[index] + weight*sample
		index += 1
		if index == len(p.Buffer) {
			index = 0
		}
	}
	p.index = index
}
