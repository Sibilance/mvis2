package analyze

import (
	math "github.com/chewxy/math32"
)

type Probe struct {
	Buffer     []float32
	memory     float32
	index      int
	downsample int
}

func (p *Probe) Analyze(signal []float32) {
	index := p.index
	downsample := p.downsample
	memory := math.Pow(p.memory, float32(int(1)<<downsample))
	weight := 1 - memory
	for _, sample := range signal {
		p.Buffer[index>>downsample] = memory*p.Buffer[index>>downsample] + weight*sample
		index += 1
		if index>>downsample == len(p.Buffer) {
			index = 0
		}
	}
	p.index = index
}
