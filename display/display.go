package display

import (
	math "github.com/chewxy/math32"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/Sibilance/mvis2/analyze"
)

type Display struct {
	Analyzer      analyze.Analyzer
	width, height int
	pixels        []byte
	finished      func() bool
}

func NewDisplay(analyzer analyze.Analyzer, finished func() bool) *Display {
	display := Display{
		Analyzer: analyzer,
		height:   len(analyzer.Probes),
		finished: finished,
	}
	for _, probePair := range analyzer.Probes {
		display.width = max(display.width, 4*len(probePair[0].Buffer))
	}
	return &display
}

func (d *Display) Update() error {
	if d.finished() {
		return ebiten.Termination
	}
	return nil
}

func (d *Display) Draw(screen *ebiten.Image) {
	if d.pixels == nil {
		d.pixels = make([]byte, 4*d.width*d.height)
	}
	pix := d.pixels
	width := d.width
	halfWidth := width / 2
	for y, probePair := range d.Analyzer.Probes {
		for x := range halfWidth {
			probeIndex := x * len(probePair[0].Buffer) / halfWidth
			leftValue := math.Abs(probePair[0].Buffer[probeIndex])
			rightValue := math.Abs(probePair[1].Buffer[probeIndex])
			pix[4*(y*width+x)] = byte(math.Min(leftValue*256*1, 255))
			pix[4*(y*width+x)+1] = byte(math.Min(leftValue*256*2, 255))
			pix[4*(y*width+x)+2] = byte(math.Min(leftValue*256*3, 255))
			pix[4*(y*width+x)+3] = 0xff
			pix[4*(y*width+x+halfWidth)] = byte(math.Min(rightValue*256*1, 255))
			pix[4*(y*width+x+halfWidth)+1] = byte(math.Min(rightValue*256*2, 255))
			pix[4*(y*width+x+halfWidth)+2] = byte(math.Min(rightValue*256*3, 255))
			pix[4*(y*width+x+halfWidth)+3] = 0xff
		}
	}
	screen.WritePixels(pix)
}

func (d *Display) Layout(outsideWidth, outsideHeight int) (int, int) {
	return d.width, d.height
}
