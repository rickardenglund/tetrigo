package main

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"

	"github.com/faiface/beep/speaker"

	"github.com/faiface/beep/mp3"
)

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	f, err := os.Open("assets/sound/tetris.mp3")
	noErr(err)
	fclock, err := os.Open("assets/sound/clock.mp3")
	noErr(err)

	tetrisStreamer, format, err := mp3.Decode(f)
	noErr(err)
	defer tetrisStreamer.Close()

	clockSteamer, clockFormat, err := mp3.Decode(fclock)
	noErr(err)
	defer tetrisStreamer.Close()

	done := make(chan bool)
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/20))
	noErr(err)

	fast := beep.ResampleRatio(4, 1, tetrisStreamer)
	myseq := beep.Seq(fast, beep.Callback(func() {
		close(done)
	}))
	speaker.Play(myseq)

	go func() {
		seconds := time.Tick(time.Second * 2)
		for {
			<-seconds
			speaker.Lock()
			err := clockSteamer.Seek(clockFormat.SampleRate.N(time.Millisecond * 170))
			noErr(err)
			speaker.Unlock()
			speaker.Play(clockSteamer)

			//fmt.Printf("clock: %v\n", clockSteamer)
			fmt.Printf("clicked\n")
		}
	}()

	<-done
	fmt.Printf("done\n")
	time.Sleep(time.Second * 3)
}
