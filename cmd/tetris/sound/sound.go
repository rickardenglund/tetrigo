package sound

import (
	"os"
	"time"

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
}

func New() *Sound {
	fclock, err := os.Open("assets/sound/clock.mp3")
	noErr(err)

	clockStamer, clockFormat, err := mp3.Decode(fclock)
	noErr(err)

	err = speaker.Init(clockFormat.SampleRate, clockFormat.SampleRate.N(time.Second/20))
	noErr(err)

	return &Sound{
		clockStreamer: clockStamer,
		clockFormat:   clockFormat,
	}
}

//func (s *Sound) Background() chan bool {
//done := make(chan bool)
//speaker.Lock()
//err := s.tetrisStreamer.Seek(0)
//noErr(err)
//speaker.Unlock()
//
//myseq := beep.Seq(s.tetrisStreamer, beep.Callback(func() {
//	close(done)
//}))
//speaker.Play(myseq)
//
//return done
//}

func (s *Sound) Click() {
	speaker.Lock()
	err := s.clockStreamer.Seek(s.clockFormat.SampleRate.N(time.Millisecond * 170))
	noErr(err)
	speaker.Unlock()
	speaker.Play(s.clockStreamer)
	//fmt.Printf("clicked\n")
}

func (s *Sound) Close() {
	s.clockStreamer.Close()
}
