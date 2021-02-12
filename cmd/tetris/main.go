package main

import (
	"Tetrigo/fonts"
	"Tetrigo/tetris"
	"Tetrigo/tetris/shape"
	"fmt"
	"image"
	_ "image/png"
	"math"
	"math/rand"
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

	blockSheet, err := loadPicture("assets/blocks.png")
	if err != nil {
		panic(err)
	}

	blockSprites := getBlockPrites(blockSheet)

	background := createBackground(win)

	nextImd := imdraw.New(nil)

	font := fonts.GetFont()
	atlas := text.NewAtlas(
		font,
		text.ASCII,
	)

	// init texts
	textPos := pixel.V(win.Bounds().Center().X+margin, win.Bounds().H()-margin)
	txt := text.New(textPos, atlas)
	nextBlockPos := pixel.V(600, 400)
	nextBlockTxt := text.New(nextBlockPos, atlas)
	highScoreTxt := text.New(nextBlockPos.Add(pixel.V(0, -100)), atlas)

	game := tetris.New()
	blockBatch := pixel.NewBatch(&pixel.TrianglesData{}, blockSheet)
	for !win.Closed() {

		gameInfo := game.GetInfo()

		win.Clear(colornames.Black)
		background.Draw(win)
		game.Tick(time.Now())

		printHUD(txt, game, gameInfo)
		txt.Draw(win, pixel.IM)

		printHighScore(highScoreTxt, game)
		highScoreTxt.Draw(win, pixel.IM)

		// draw blocks
		boxWidth, boxHeight := getBoxSize(gameInfo.Width, gameInfo.Height, win.Bounds())
		boxScale := getBoxScale(boxWidth, boxHeight, blockSprites[2].Frame().Size())
		blocks := game.GetBlocks()
		drawBlocks(blockBatch, blocks, win, gameInfo, boxWidth, boxHeight, boxScale, blockSprites)
		blockBatch.Draw(win)

		// Draw next block
		ns := game.NextBlock()
		drawNextBlock(nextBlockPos, boxWidth, nextBlockTxt, boxHeight, ns, nextImd)
		nextImd.Draw(win)
		nextBlockTxt.Clear()
		fmt.Fprintf(nextBlockTxt, "Next")
		nextBlockTxt.Draw(win, pixel.IM)

		handleInput(win, &game)
		if win.JustPressed(pixelgl.KeyEnter) {
			game = tetris.New()
		}

		win.Update()
	}

}

func getBlockPrites(sheet pixel.Picture) []*pixel.Sprite {
	const spriteWidth = 64
	sprites := make([]*pixel.Sprite, 0, int(sheet.Bounds().W()/spriteWidth))
	for x := 0.0; x < sheet.Bounds().W(); x += spriteWidth {
		sprite := pixel.NewSprite(sheet, pixel.R(x, 0.0, x+spriteWidth, sheet.Bounds().H()))
		sprites = append(sprites, sprite)
	}

	return sprites
}

func handleInput(win *pixelgl.Window, game *tetris.Game) {
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
	if win.JustPressed(pixelgl.KeyP) {
		game.TogglePaus()
	}
}

func drawNextBlock(nextBlockPos pixel.Vec, boxWidth float64, nextBlockTxt *text.Text, boxHeight float64, ns shape.Shape, nextImd *imdraw.IMDraw) {
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
}

func drawBlocks(batch *pixel.Batch, blocks []shape.Block, win *pixelgl.Window, gameInfo tetris.Info, boxWidth float64, boxHeight float64, boxScale pixel.Vec, blockSprites []*pixel.Sprite) {
	batch.Clear()
	for i := range blocks {
		pos := getBlockPos(win.Bounds(), gameInfo.Width, gameInfo.Height, boxWidth, blocks[i].Pos)
		pos = pos.Add(pixel.V(boxWidth/2, boxHeight/2))
		pos.X = math.Floor(pos.X)
		pos.Y = math.Floor(pos.Y)
		m := pixel.IM.ScaledXY(pixel.ZV, boxScale)
		m = m.Moved(pos)
		blockSprites[blocks[i].Kind].Draw(batch, m)
	}
}

func printHighScore(highScoreTxt *text.Text, game tetris.Game) {
	highScoreTxt.Clear()
	fmt.Fprintf(highScoreTxt, "Score: ")
	highScoreTxt.Dot = highScoreTxt.Dot.Add(pixel.V(50, 0))
	levelPosX := highScoreTxt.Dot.X
	fmt.Fprintf(highScoreTxt, "Level: \n")
	for _, result := range game.Results {
		fmt.Fprintf(highScoreTxt, "%d", result.Score)
		highScoreTxt.Dot.X = levelPosX
		fmt.Fprintf(highScoreTxt, "%d\n", result.Level)
	}
}

func printHUD(txt *text.Text, game tetris.Game, gameInfo tetris.Info) {
	txt.Clear()
	fmt.Fprintf(txt, "Score: %d\n", game.GetScore())
	fmt.Fprintf(txt, "Level: %d", gameInfo.Level)
	if game.IsGameOver() {
		fmt.Fprintf(txt, "\nGame Over")
	}
	if gameInfo.Paused {
		fmt.Fprintf(txt, "\n\nPaused")
	}
}

func createBackground(win *pixelgl.Window) *imdraw.IMDraw {
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
	return background
}

func getBoxScale(desiredWidth, desiredHeight float64, size pixel.Vec) pixel.Vec {
	xScale := desiredWidth / size.X
	yScale := desiredHeight / size.Y

	return pixel.V(xScale, yScale)
}

func getShapePoints(base pixel.Vec, boxWidth, boxHeight float64, blocks []shape.Block) []pixel.Vec {
	res := make([]pixel.Vec, 0, len(blocks)*4)
	for _, block := range blocks {
		pv := pixel.Vec{X: float64(block.Pos.X) * boxWidth, Y: float64(block.Pos.Y) * boxHeight}.Add(base)
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
	rand.Seed(time.Now().Unix())
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
