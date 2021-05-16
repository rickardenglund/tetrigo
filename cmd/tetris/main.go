package main

import (
	"flag"
	"fmt"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/rickardenglund/tetrigo/cmd/tetris/renderer"
	"github.com/rickardenglund/tetrigo/cmd/tetris/soundcontroller"
	"github.com/rickardenglund/tetrigo/sound"
	"github.com/rickardenglund/tetrigo/tetris"
	"github.com/rickardenglund/tetrigo/timestat"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type CtlState struct {
	falling     bool
	previousAge int
}

func run(vsync bool) {
	sounds := sound.New()
	defer sounds.Close()

	windowCfg := pixelgl.WindowConfig{
		Title:     "TetriGo",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     vsync,
		Resizable: true,
	}

	win, err := pixelgl.NewWindow(windowCfg)
	if err != nil {
		panic(err)
	}

	r := renderer.New("assets/blocks.png")

	game := tetris.New()
	ctlState := CtlState{}

	const timeMeasurementsToStore = 240

	frames := 0
	ticker := time.NewTicker(time.Second)
	timeBuffer := timestat.New(timeMeasurementsToStore)

	for !win.Closed() {
		startRender := time.Now()
		explodedBlocks := game.Tick(startRender)
		gameInfo := game.GetInfo()

		soundcontroller.ControlSound(ctlState.previousAge > gameInfo.ActiveAge, explodedBlocks, sounds)

		r.Render(win, gameInfo, &game, explodedBlocks, timeBuffer)
		ctlState.handleInput(win, &game, &gameInfo)

		if win.JustPressed(pixelgl.KeyEnter) {
			game = tetris.New()
		}

		select {
		case <-ticker.C:
			win.SetTitle(fmt.Sprintf("fps: %v", frames))
			frames = 0
		default:
		}

		frames++

		timeBuffer.Add(time.Since(startRender))

		win.Update()
	}
}

func (c *CtlState) handleInput(win *pixelgl.Window, game *tetris.Game, gi *tetris.Info) {
	if win.JustPressed(pixelgl.KeyDown) || win.JustPressed(pixelgl.KeySpace) {
		c.falling = true
	}

	if gi.ActiveAge < c.previousAge {
		c.falling = false
	}

	c.previousAge = gi.ActiveAge

	if win.Pressed(pixelgl.KeySpace) || win.Pressed(pixelgl.KeyDown) && c.falling {
		game.Speed()
	}

	if win.JustPressed(pixelgl.KeyLeft) {
		game.Left()
	}

	if win.JustPressed(pixelgl.KeyRight) {
		game.Right()
	}

	if win.JustPressed(pixelgl.KeyUp) {
		game.Rotate()
	}

	if win.JustPressed(pixelgl.KeyP) {
		game.TogglePause()
	}
}

func main() {
	vsync := flag.Bool("vsync", true, "enable vsync")
	flag.Parse()

	rand.Seed(time.Now().Unix())
	pixelgl.Run(func() {
		run(*vsync)
	})
}
