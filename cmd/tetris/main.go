package main

import (
	"Tetrigo/fonts"
	"Tetrigo/tetris"
	"Tetrigo/tetris/shape"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"time"
)

func run() {
	windowCfg := pixelgl.WindowConfig{
		Title:  "TetriGo",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(windowCfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	imd.Color = colornames.Gray
	imd.Push(pixel.V(margin, win.Bounds().H()-margin))
	imd.Push(pixel.V(win.Bounds().Center().X-margin, win.Bounds().H()-margin))
	imd.Push(pixel.V(win.Bounds().Center().X-margin, margin))
	imd.Push(pixel.V(margin, margin))
	imd.Polygon(5)

	imd.Push(pixel.V(win.Bounds().Center().X, win.Bounds().H()))
	imd.Push(pixel.V(win.Bounds().Center().X, 0))
	imd.Line(10)

	activeBlock := imdraw.New(nil)

	font := fonts.GetFont()
	atlas := text.NewAtlas(
		font,
		text.ASCII,
	)
	textPos := pixel.V(win.Bounds().Center().X+margin, win.Bounds().H()-margin)
	txt := text.New(textPos, atlas)

	game := tetris.New()
	gameWidth, gameHeight := game.GetDimensions()
	for !win.Closed() {
		win.Clear(colornames.Black)
		imd.Draw(win)
		game.Tick(time.Now())

		txt.Clear()
		fmt.Fprintf(txt, "Score: %d", game.GetScore())
		if game.IsGameOver() {
			fmt.Fprintf(txt, "\nGame Over")
		}
		txt.Draw(win, pixel.IM)

		// draw blocks
		blocks := game.GetBlocks()
		//fmt.Printf("blocks: %v\n", blocks)
		activeBlock.Clear()
		for i := range blocks {
			pos, boxWidth, boxHeight := getBlockPos(win.Bounds(), gameWidth, gameHeight, blocks[i])
			activeBlock.Color = colornames.Red
			activeBlock.Push(
				pos,
				pos.Add(pixel.V(boxWidth, boxHeight)),
			)
			activeBlock.Rectangle(0)
		}
		activeBlock.Draw(win)

		if win.Pressed(pixelgl.KeySpace) || win.Pressed(pixelgl.KeyDown) {
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
		if win.JustPressed(pixelgl.KeyEnter) {
			game = tetris.New()
		}

		win.Update()
	}

}

const margin = 50

func getBlockPos(bounds pixel.Rect, gameWidth, gameHeight int, pos shape.Pos) (pixel.Vec, float64, float64) {
	boardLeft := float64(margin)
	boardTop := bounds.H() - margin
	boardBottom := float64(margin)
	boardRight := bounds.Center().X - margin

	boardWidth := boardRight - boardLeft
	boardHeight := boardTop - boardBottom

	boxWidth := boardWidth / float64(gameWidth)
	boxHeight := boardHeight / float64(gameHeight)

	out := pixel.Vec{}
	out.X = mapRange(float64(pos.X), 0, float64(gameWidth), boardLeft, boardRight)
	out.Y = mapRange(float64(pos.Y), 0, float64(gameHeight), boardBottom, boardTop-boxWidth)
	return out, boxWidth, boxHeight
}

func mapRange(input, inputStart, inputEnd, outputStart, outputEnd float64) float64 {
	return ((input-inputStart)/(inputEnd-inputStart))*(outputEnd-outputStart) + outputStart
}

func main() {
	pixelgl.Run(run)
}
