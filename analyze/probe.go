package analyze

type Probe struct {
	Buffer []float32
	memory float32
	index  int
}

func (p *Probe) Analyze(signal []float32) {
	index := p.index
	for _, sample := range signal {
		p.Buffer[index] = p.memory*p.Buffer[index] + sample
		index += 1
		if index == len(p.Buffer) {
			index = 0
		}
	}
	p.index = index
}
