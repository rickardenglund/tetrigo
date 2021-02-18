package renderer

import (
	"Tetrigo/cmd/tetris/hud"
	"Tetrigo/cmd/tetris/sprites"
	"Tetrigo/fonts"
	"Tetrigo/tetris"
	"Tetrigo/tetris/shape"
	"image/color"
	"math"
	"math/rand"

	"github.com/faiface/pixel/text"

	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"
)

type Renderer struct {
	blockSprites  []*pixel.Sprite
	blockBatch    *pixel.Batch
	fallingBlocks map[int]FallingBlock
	atlas         *text.Atlas
}

type FallingBlock struct {
	pos      pixel.Vec
	rotation float64
	kind     int
	xs       float64
	ys       float64
}

func New(spritePath string) Renderer {
	blockSheet, err := sprites.LoadPicture(spritePath)
	if err != nil {
		panic(err)
	}

	font := fonts.GetFont()
	atlas := text.NewAtlas(
		font,
		text.ASCII,
	)

	blockSprites := sprites.GetBlockSprites(blockSheet)

	blockBatch := pixel.NewBatch(&pixel.TrianglesData{}, blockSheet)

	return Renderer{
		blockSprites:  blockSprites,
		blockBatch:    blockBatch,
		atlas:         atlas,
		fallingBlocks: map[int]FallingBlock{},
	}
}

func (r *Renderer) Render(win *pixelgl.Window, gameInfo tetris.Info, game *tetris.Game, explodedBlocks []shape.Block) {
	boxSize := calculateBoxSize(gameInfo.Width, gameInfo.Height, win.Bounds())
	boxScale := getBoxScale(boxSize, boxSize, r.blockSprites[2].Frame().Size())

	background := createBackground(win.Bounds(), boxSize, float64(gameInfo.Width), float64(gameInfo.Height))
	textPos := pixel.V(win.Bounds().Center().X+margin, win.Bounds().H()-margin)
	hud := hud.New(textPos, r.atlas)

	win.Clear(colornames.Black)
	background.Draw(win)
	hud.DrawHUD(win, game, gameInfo, boxSize, boxSize)

	blockIDCounter := 0
	toFallingBlocks(explodedBlocks, win, gameInfo, boxSize, blockIDCounter, r.fallingBlocks)
	movaFallingBlocks(r.fallingBlocks)

	r.blockBatch.Clear()
	for _, v := range r.fallingBlocks {
		m := pixel.IM.ScaledXY(pixel.ZV, boxScale)
		m = m.Moved(v.pos)
		r.blockSprites[v.kind].Draw(r.blockBatch, m)
	}

	// draw blocks
	drawBlocks(r.blockBatch, game.GetBlocks(), win, gameInfo, boxSize, boxScale, r.blockSprites)
	r.blockBatch.Draw(win)
}

func drawBlocks(batch *pixel.Batch, blocks []shape.Block, win *pixelgl.Window, gameInfo tetris.Info, boxSize float64, boxScale pixel.Vec, blockSprites []*pixel.Sprite) {
	for i := range blocks {
		pos := getBlockPos(win.Bounds(), gameInfo, boxSize, blocks[i].Pos)
		pos = pos.Add(pixel.V(boxSize/2, boxSize/2))
		pos.X = math.Floor(pos.X)
		pos.Y = math.Floor(pos.Y)
		m := pixel.IM.ScaledXY(pixel.ZV, boxScale)
		m = m.Moved(pos)
		blockSprites[blocks[i].Kind].Draw(batch, m)
	}
}

func getBoxScale(desiredWidth, desiredHeight float64, size pixel.Vec) pixel.Vec {
	xScale := (desiredWidth - 1) / size.X
	yScale := (desiredHeight - 1) / size.Y

	return pixel.V(xScale, yScale)
}

func toFallingBlocks(explodedBlocks []shape.Block, win *pixelgl.Window, gameInfo tetris.Info, boxWidth float64, blockIDCounter int, fallingBlocks map[int]FallingBlock) {
	const (
		xSpread = 8
		ySpread = 4
	)
	for _, b := range explodedBlocks {
		fb := FallingBlock{
			pos:      getBlockPos(win.Bounds(), gameInfo, boxWidth, b.Pos),
			rotation: 0,
			kind:     b.Kind,
			xs:       rand.Float64()*xSpread - xSpread/2,
			ys:       rand.Float64()*ySpread - ySpread,
		}

		blockIDCounter++
		fallingBlocks[blockIDCounter] = fb
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

func getBlockPos(bounds pixel.Rect, info tetris.Info, boxSize float64, pos shape.Pos) pixel.Vec {
	gameWidth := float64(info.Width) * boxSize
	space := bounds.W()/2 - gameWidth

	return pixel.V(float64(pos.X)*boxSize, float64(pos.Y)*boxSize).Add(pixel.V(space/2, margin))
}

func createBackground(bounds pixel.Rect, size, width, height float64) *imdraw.IMDraw {
	const lineWidth = 8

	gameWidth := width * size
	spaceX := bounds.W()/2 - gameWidth
	gameHeight := height * size
	spaceY := bounds.H() - gameHeight

	background := imdraw.New(nil)
	background.Color = colornames.Gray
	vertices := []pixel.Vec{
		pixel.V(spaceX/2-lineWidth/2, spaceY/2+gameHeight),
		pixel.V(spaceX/2-lineWidth/2, spaceY/2-lineWidth/2),
		pixel.V(spaceX/2+size*width+lineWidth/2, spaceY/2-lineWidth/2),
		pixel.V(spaceX/2+size*width+lineWidth/2, spaceY/2+gameHeight),
	}
	background.Push(vertices...)
	background.Line(lineWidth / 2)

	background.Color = color.RGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xFF}
	for i := 0; i < int(width); i++ {
		vertices := []pixel.Vec{
			pixel.V(spaceX/2+size*float64(i), spaceY/2+gameHeight),
			pixel.V(spaceX/2+size*float64(i), spaceY/2-lineWidth/2),
		}
		background.Push(vertices...)
		background.Line(1)

	}

	background.Push(pixel.V(bounds.Center().X, bounds.H()))
	background.Push(pixel.V(bounds.Center().X, 0))
	background.Line(lineWidth * 2)

	return background
}

const margin = 50

func calculateBoxSize(gameWidth, gameHeight int, bounds pixel.Rect) float64 {
	boardLeft := float64(margin)
	boardTop := bounds.H() - margin
	boardBottom := float64(margin)
	boardRight := bounds.Center().X - margin

	boardWidth := boardRight - boardLeft
	boardHeight := boardTop - boardBottom

	boxWidth := boardWidth / float64(gameWidth)
	boxHeight := boardHeight / float64(gameHeight)

	return math.Min(boxWidth, boxHeight)
}
