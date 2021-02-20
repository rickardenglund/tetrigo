package hud

import (
	"Tetrigo/tetris"
	"Tetrigo/tetris/shape"
	"Tetrigo/timestat"
	"fmt"
	"time"

	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"

	"github.com/faiface/pixel/text"
)

type Hud struct {
	txt          *text.Text
	highScoreTxt *text.Text
	nextBlockTxt *text.Text
	nextBlockPos pixel.Vec
	imd          *imdraw.IMDraw
}

func New(pos pixel.Vec, atlas *text.Atlas) Hud {
	const (
		yOffset = 400
		xOffset = 50
	)

	nextBlockPos := pixel.V(pos.X+xOffset, yOffset)
	highScorePos := nextBlockPos.Add(pixel.V(0, -100))

	hudTxt := text.New(pos, atlas)
	highScoreTxt := text.New(highScorePos, atlas)
	nextBlockTxt := text.New(nextBlockPos, atlas)

	imd := imdraw.New(nil)

	return Hud{txt: hudTxt, highScoreTxt: highScoreTxt, nextBlockTxt: nextBlockTxt, nextBlockPos: nextBlockPos, imd: imd}
}

func (h *Hud) DrawHUD(win *pixelgl.Window, game *tetris.Game, gameInfo tetris.Info,
	boxWidth float64, boxHeight float64, buffer timestat.TimeStat) {
	fmt.Fprintf(h.txt, "Score: %d\n", game.GetScore())
	fmt.Fprintf(h.txt, "Level: %d", gameInfo.Level)

	h.txt.Clear()

	if game.IsGameOver() {
		fmt.Fprintf(h.txt, "\nGame Over")
	}

	if gameInfo.Paused {
		fmt.Fprintf(h.txt, "\n\nPaused")
	}

	h.txt.Draw(win, pixel.IM)

	printHighScore(h.highScoreTxt, game)
	h.highScoreTxt.Draw(win, pixel.IM)

	// Draw next block
	ns := game.NextBlock()
	drawNextBlock(h.nextBlockPos, boxWidth, h.nextBlockTxt, boxHeight, ns, h.imd)
	h.imd.Draw(win)
	h.nextBlockTxt.Clear()
	fmt.Fprintf(h.nextBlockTxt, "Next")
	h.nextBlockTxt.Draw(win, pixel.IM)

	// Draw stats
	const statMargin = 20
	p := pixel.V(win.Bounds().Center().X+statMargin, statMargin)

	h.imd.Clear()

	for i, duration := range buffer.Values() {
		if duration > time.Second/60 {
			h.imd.Color = colornames.Red
		} else if duration > time.Second/120 {
			h.imd.Color = colornames.Yellow
		} else {
			h.imd.Color = colornames.Green
		}

		const (
			heightScaling = 2
			statLineWidth = 2
		)

		i *= 2
		h.imd.Push(p.Add(pixel.V(float64(i), float64(duration.Milliseconds()*heightScaling))))
		h.imd.Push(p.Add(pixel.V(float64(i), 0)))
		h.imd.Line(statLineWidth)
	}

	h.imd.Draw(win)
}

func printHighScore(highScoreTxt *text.Text, game *tetris.Game) {
	const widthBetweenRows = 50

	highScoreTxt.Clear()
	fmt.Fprintf(highScoreTxt, "Score: ")
	highScoreTxt.Dot = highScoreTxt.Dot.Add(pixel.V(widthBetweenRows, 0))
	levelPosX := highScoreTxt.Dot.X

	for _, result := range game.Results {
		fmt.Fprintf(highScoreTxt, "%d", result.Score)
		highScoreTxt.Dot.X = levelPosX
		fmt.Fprintf(highScoreTxt, "%d\n", result.Level)
	}
}

func drawNextBlock(nextBlockPos pixel.Vec, boxWidth float64, nextBlockTxt *text.Text,
	boxHeight float64, ns shape.Shape, nextImd *imdraw.IMDraw) {
	const lineThickness = 3

	points := getShapePoints(
		nextBlockPos.Add(pixel.V(boxWidth, nextBlockTxt.LineHeight*1.5)), // nolint: gomnd // offset
		boxWidth, boxHeight, ns.GetBlocks())

	nextImd.Clear()

	nextImd.Color = colornames.Greenyellow

	i := 0
	for i < len(points) {
		for j := 0; j < 4; j++ {
			nextImd.Push(points[i])
			i++
		}
		nextImd.Polygon(lineThickness)
	}
}

func getShapePoints(base pixel.Vec, boxWidth, boxHeight float64, blocks []shape.Block) []pixel.Vec {
	const pointsPerBlock = 4

	res := make([]pixel.Vec, 0, len(blocks)*pointsPerBlock)

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
