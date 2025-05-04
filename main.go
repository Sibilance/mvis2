package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/Sibilance/mvis2/analyze"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

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

	decodedMp3, err := mp3.NewDecoder(file)
	if err != nil {
		log.Fatalf("mp3.NewDecoder failed: %s", err)
	}

	op := &oto.NewContextOptions{}
	op.SampleRate = decodedMp3.SampleRate()
	op.ChannelCount = 2
	op.Format = oto.FormatSignedInt16LE

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		log.Fatalf("oto.NewContext failed: %s", err)
	}
	<-readyChan

	analyzer := analyze.NewAnalyzer(decodedMp3, op)

	player := otoCtx.NewPlayer(analyzer)
	defer player.Close()
	player.Play()
	for player.IsPlaying() {
		player.BufferedSize()
		time.Sleep(time.Millisecond)
	}
}
