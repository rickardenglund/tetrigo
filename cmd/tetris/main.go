package main

import (
	"Tetrigo/cmd/tetris/renderer"
	"Tetrigo/cmd/tetris/soundController"
	"Tetrigo/sound"
	"Tetrigo/tetris"
	"fmt"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type CtlState struct {
	falling     bool
	previousAge int
}

func run() {
	sounds := sound.New()
	defer sounds.Close()

	windowCfg := pixelgl.WindowConfig{
		Title:     "TetriGo",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(windowCfg)
	if err != nil {
		panic(err)
	}

	r := renderer.New("assets/blocks.png")

	game := tetris.New()
	ctlState := CtlState{}

	frames := 0
	ticker := time.NewTicker(time.Second)

	for !win.Closed() {
		startRender := time.Now()
		explodedBlocks := game.Tick(startRender)
		gameInfo := game.GetInfo()

		soundController.ControlSound(ctlState.previousAge > gameInfo.ActiveAge, explodedBlocks, sounds)

		r.Render(win, gameInfo, &game, explodedBlocks)
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
	rand.Seed(time.Now().Unix())
	pixelgl.Run(run)
}
