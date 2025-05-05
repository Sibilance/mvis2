package main

import (
	"flag"
	"io"
	"log"
	"math"
	"os"

	"github.com/Sibilance/mvis2/analyze"
	"github.com/Sibilance/mvis2/display"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/go-mp3"
)

const MaxInt16 = float32(math.MaxInt16)

type SampleExtractor struct {
	reader                  io.Reader
	analyzer                *analyze.Analyzer
	leftBuffer, rightBuffer []float32
}

func (se SampleExtractor) Read(p []byte) (n int, err error) {
	n, err = se.reader.Read(p)
	nSamples := n >> 2
	if cap(se.leftBuffer) < nSamples {
		se.leftBuffer = make([]float32, nSamples)
	}
	if cap(se.rightBuffer) < nSamples {
		se.rightBuffer = make([]float32, nSamples)
	}
	for i := 0; i < nSamples; i++ {
		sampleIndex := i << 2
		se.leftBuffer[i] = float32(
			int16(
				uint16(p[sampleIndex])|uint16(p[sampleIndex+1])<<8,
			),
		) / MaxInt16
		se.rightBuffer[i] = float32(
			int16(
				uint16(p[sampleIndex+2])|uint16(p[sampleIndex+3])<<8,
			),
		) / MaxInt16
	}
	se.analyzer.Analyze(se.leftBuffer, se.rightBuffer)
	return
}

func main() {
	fileName := flag.String("file", "", "mp3 file to play")
	flag.Parse()

	if *fileName == "" {
		flag.Usage()
		log.Fatal("missing file argument")
	}

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatalf("error opening '%s': %s", *fileName, err)
	}
	defer file.Close()

	mp3Decoder, err := mp3.NewDecoder(file)
	if err != nil {
		log.Fatalf("mp3.NewDecoder failed: %s", err)
	}

	op := &oto.NewContextOptions{}
	op.SampleRate = mp3Decoder.SampleRate()
	op.ChannelCount = 2
	op.Format = oto.FormatSignedInt16LE

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		log.Fatalf("oto.NewContext failed: %s", err)
	}
	<-readyChan

	analyzer := analyze.NewAnalyzer()
	sampleExtractor := SampleExtractor{
		reader:   mp3Decoder,
		analyzer: analyzer,
	}

	player := otoCtx.NewPlayer(sampleExtractor)
	defer player.Close()

	game := display.NewDisplay(
		*analyzer,
		func() bool { return !player.IsPlaying() },
	)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Music Visualizer")

	player.Play()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
