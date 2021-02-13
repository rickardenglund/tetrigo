package main

import (
	"Tetrigo/cmd/tetris/hud"
	"Tetrigo/cmd/tetris/sound"
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

type CtlState struct {
	falling     bool
	previousAge int
}

type FallingBlock struct {
	pos      pixel.Vec
	rotation float64
	kind     int
	xs       float64
	ys       float64
}

func run() {
	sounds := sound.New()
	defer sounds.Close()

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

	blockSprites := getBlockSprites(blockSheet)

	background := createBackground(win)

	nextImd := imdraw.New(nil)

	font := fonts.GetFont()
	atlas := text.NewAtlas(
		font,
		text.ASCII,
	)

	bcnt := 0
	fallingBlocks := map[int]FallingBlock{}

	// init texts
	textPos := pixel.V(win.Bounds().Center().X+margin, win.Bounds().H()-margin)
	nextBlockPos := pixel.V(600, 400)
	nextBlockTxt := text.New(nextBlockPos, atlas)

	hud := hud.New(textPos, atlas, nextBlockPos.Add(pixel.V(0, -100)))

	game := tetris.New()
	blockBatch := pixel.NewBatch(&pixel.TrianglesData{}, blockSheet)
	ctlState := CtlState{}
	frames := 0
	ticker := time.NewTicker(time.Second)

	for !win.Closed() {
		startRender := time.Now()
		explodedBlocks := game.Tick(startRender)
		gameInfo := game.GetInfo()

		boxWidth, boxHeight := getBoxSize(gameInfo.Width, gameInfo.Height, win.Bounds())
		boxScale := getBoxScale(boxWidth, boxHeight, blockSprites[2].Frame().Size())

		if len(explodedBlocks) > 0 {
			sounds.Click()
			for _, b := range explodedBlocks {
				fb := FallingBlock{
					pos:      getBlockPos(win.Bounds(), gameInfo.Width, gameInfo.Height, boxWidth, b.Pos),
					rotation: 0,
					kind:     b.Kind,
					xs:       rand.Float64()*-2 + 1,
					ys:       rand.Float64()*2 - 1,
				}

				bcnt++
				fallingBlocks[bcnt] = fb
			}
		}

		movaFallingBlocks(fallingBlocks)

		win.Clear(colornames.Black)
		background.Draw(win)

		hud.DrawHUD(win, game, gameInfo)

		blockBatch.Clear()
		for _, v := range fallingBlocks {
			m := pixel.IM.ScaledXY(pixel.ZV, boxScale)
			m = m.Moved(v.pos)
			blockSprites[v.kind].Draw(blockBatch, m)
		}

		// draw blocks
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

func movaFallingBlocks(blocks map[int]FallingBlock) {
	for k, v := range blocks {
		v.ys -= 1.0
		v.xs *= 0.9
		v.pos.Y += v.ys
		v.pos.X += v.xs
		blocks[k] = v
		if v.pos.Y < 0 {
			delete(blocks, k)
		}
	}
}

func getBlockSprites(sheet pixel.Picture) []*pixel.Sprite {
	const spriteWidth = 64
	sprites := make([]*pixel.Sprite, 0, int(sheet.Bounds().W()/spriteWidth))
	for x := 0.0; x < sheet.Bounds().W(); x += spriteWidth {
		sprite := pixel.NewSprite(sheet, pixel.R(x, 0.0, x+spriteWidth, sheet.Bounds().H()))
		sprites = append(sprites, sprite)
	}

	return sprites
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
