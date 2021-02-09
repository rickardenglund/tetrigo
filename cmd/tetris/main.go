package main

import (
	"Tetrigo/fonts"
	"Tetrigo/tetris"
	"Tetrigo/tetris/shape"
	"fmt"
	"image"
	_ "image/png"
	"math"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
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

	blockPic, err := loadPicture("assets/blocks/red2.png")
	if err != nil {
		panic(err)
	}

	blockSprite := pixel.NewSprite(blockPic, blockPic.Bounds())

	background := imdraw.New(nil)
	background.Color = colornames.Gray
	background.Push(pixel.V(margin, win.Bounds().H()-margin))
	background.Push(pixel.V(win.Bounds().Center().X-margin, win.Bounds().H()-margin))
	background.Push(pixel.V(win.Bounds().Center().X-margin, margin))
	background.Push(pixel.V(margin, margin))
	background.Polygon(5)

	background.Push(pixel.V(win.Bounds().Center().X, win.Bounds().H()))
	background.Push(pixel.V(win.Bounds().Center().X, 0))
	background.Line(10)

	gameImd := imdraw.New(nil)
	nextImd := imdraw.New(nil)

	font := fonts.GetFont()
	atlas := text.NewAtlas(
		font,
		text.ASCII,
	)
	textPos := pixel.V(win.Bounds().Center().X+margin, win.Bounds().H()-margin)
	txt := text.New(textPos, atlas)
	nextBlockPos := pixel.V(600, 400)
	nextBlockTxt := text.New(nextBlockPos, atlas)

	game := tetris.New()
	paused := false
	gameWidth, gameHeight := game.GetDimensions()
	batch := pixel.NewBatch(&pixel.TrianglesData{}, blockPic)
	for !win.Closed() {
		win.Clear(colornames.Black)
		background.Draw(win)
		if !paused {
			game.Tick(time.Now())
		}

		txt.Clear()
		fmt.Fprintf(txt, "Score: %d", game.GetScore())
		if game.IsGameOver() {
			fmt.Fprintf(txt, "\nGame Over")
		}
		if paused {
			fmt.Fprintf(txt, "\n\nPaused")
		}
		txt.Draw(win, pixel.IM)

		// draw blocks
		blocks := game.GetBlocks()
		//fmt.Printf("blocks: %v\n", blocks)
		gameImd.Clear()
		boxWidth, boxHeight := getBoxSize(gameWidth, gameHeight, win.Bounds())
		boxScale := getBoxScale(boxWidth, boxHeight, blockSprite.Picture().Bounds())
		batch.Clear()
		for i := range blocks {
			pos := getBlockPos(win.Bounds(), gameWidth, gameHeight, boxWidth, blocks[i])
			pos = pos.Add(pixel.V(boxWidth/2, boxHeight/2))
			pos.X = math.Floor(pos.X)
			pos.Y = math.Floor(pos.Y)
			m := pixel.IM.ScaledXY(pixel.ZV, boxScale)
			m = m.Moved(pos)
			blockSprite.Draw(batch, m)
			//gameImd.Color = colornames.Red
			//gameImd.Push(
			//	pos,
			//	pos.Add(pixel.V(boxWidth, boxHeight)),
			//)
			//gameImd.Rectangle(0)
		}
		batch.Draw(win)
		//gameImd.Draw(win)

		// Draw next block
		ns := game.NextBlock()
		points := getShapePoints(nextBlockPos.Add(pixel.V(boxWidth, nextBlockTxt.LineHeight*2.5)), boxWidth, boxHeight, ns.GetBlocks())
		nextImd.Clear()
		i := 0
		nextImd.Color = colornames.Greenyellow
		for i < len(points) {
			for j := 0; j < 4; j++ {
				nextImd.Push(points[i])
				i++
			}
			nextImd.Polygon(3)
		}
		nextImd.Draw(win)

		nextBlockTxt.Clear()
		fmt.Fprintf(nextBlockTxt, "Next")
		nextBlockTxt.Draw(win, pixel.IM)

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
		if win.JustPressed(pixelgl.KeyP) {
			paused = !paused
		}

		win.Update()
	}

}

func getBoxScale(desiredWidth, desiredHeight float64, bounds pixel.Rect) pixel.Vec {
	xScale := desiredWidth / bounds.W()
	yScale := desiredHeight / bounds.H()

	return pixel.V(xScale, yScale)
}

func getShapePoints(base pixel.Vec, boxWidth, boxHeight float64, blocks []shape.Pos) []pixel.Vec {
	res := make([]pixel.Vec, 0, len(blocks)*4)
	for _, p := range blocks {
		pv := pixel.Vec{X: float64(p.X) * boxWidth, Y: float64(p.Y) * boxHeight}.Add(base)
		res = append(res,
			pv,
			pv.Add(pixel.V(boxWidth, 0)),
			pv.Add(pixel.V(boxWidth, boxHeight)),
			pv.Add(pixel.V(0, boxHeight)),
		)
	}
	return res
}

const margin = 50

func getBoxSize(gameWidth, gameHeight int, bounds pixel.Rect) (float64, float64) {
	boardLeft := float64(margin)
	boardTop := bounds.H() - margin
	boardBottom := float64(margin)
	boardRight := bounds.Center().X - margin

	boardWidth := boardRight - boardLeft
	boardHeight := boardTop - boardBottom

	boxWidth := boardWidth / float64(gameWidth)
	boxHeight := boardHeight / float64(gameHeight)

	return boxWidth, boxHeight
}

func getBlockPos(bounds pixel.Rect, gameWidth, gameHeight int, boxWidth float64, pos shape.Pos) pixel.Vec {
	boardLeft := float64(margin)
	boardTop := bounds.H() - margin
	boardBottom := float64(margin)
	boardRight := bounds.Center().X - margin

	out := pixel.Vec{}
	out.X = mapRange(float64(pos.X), 0, float64(gameWidth), boardLeft, boardRight)
	out.Y = mapRange(float64(pos.Y), 0, float64(gameHeight)-2, boardBottom, boardTop-boxWidth)
	return out
}

func mapRange(input, inputStart, inputEnd, outputStart, outputEnd float64) float64 {
	return ((input-inputStart)/(inputEnd-inputStart))*(outputEnd-outputStart) + outputStart
}

func main() {
	pixelgl.Run(run)
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
