package analyze

import (
	"io"

	"github.com/ebitengine/oto/v3"
)

type Analyzer struct {
	reader io.Reader
	op     *oto.NewContextOptions
	pos    uint64
}

func NewAnalyzer(r io.Reader, op *oto.NewContextOptions) *Analyzer {
	return &Analyzer{
		reader: r,
		op:     op,
	}
}

func (a *Analyzer) Read(p []byte) (n int, err error) {
	n, err = a.reader.Read(p)
	a.pos += uint64(n)
	return
}
