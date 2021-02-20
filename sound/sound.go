package sound

import (
	"os"
	"time"

	"github.com/faiface/beep/effects"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Sound struct {
	clockStreamer beep.StreamSeekCloser
	clockFormat   beep.Format
	tickFormat    beep.Format
	tickStreamer  beep.StreamSeekCloser
}

func New() *Sound {
	const (
		clockPath = "assets/sound/clock.mp3"
		tickPath  = "assets/sound/tick.mp3"
	)

	clockStreamer, clockFormat := readSound(clockPath)
	tickStreamer, tickFormat := readSound(tickPath)

	err := speaker.Init(clockFormat.SampleRate, clockFormat.SampleRate.N(time.Second/20)) //nolint: gomnd // init magic
	noErr(err)

	return &Sound{
		clockStreamer: clockStreamer,
		clockFormat:   clockFormat,
		tickStreamer:  tickStreamer,
		tickFormat:    tickFormat,
	}
}

func readSound(path string) (beep.StreamSeekCloser, beep.Format) {
	fclock, err := os.Open(path)
	noErr(err)

	clockStreamer, clockFormat, err := mp3.Decode(fclock)
	noErr(err)

	return clockStreamer, clockFormat
}

func (s *Sound) Click() {
	speaker.Lock()
	const clickTimeOffset = time.Millisecond * 170
	err := s.clockStreamer.Seek(s.clockFormat.SampleRate.N(clickTimeOffset))
	noErr(err)
	speaker.Unlock()
	speaker.Play(s.clockStreamer)
}

func (s *Sound) Tick() {
	speaker.Lock()
	err := s.tickStreamer.Seek(0)
	noErr(err)
	speaker.Unlock()

	volume := &effects.Volume{
		Streamer: s.tickStreamer,
		Base:     2,
		Volume:   -2,
		Silent:   false,
	}

	speaker.Play(volume)
}

func (s *Sound) Close() {
	s.clockStreamer.Close()
}
