package hud

import (
	"Tetrigo/tetris"
	"Tetrigo/tetris/shape"
	"fmt"

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
	nextBlockPos := pixel.V(pos.X+50, 400)
	highScorePos := nextBlockPos.Add(pixel.V(0, -100))

	hudTxt := text.New(pos, atlas)
	highScoreTxt := text.New(highScorePos, atlas)
	nextBlockTxt := text.New(nextBlockPos, atlas)

	imd := imdraw.New(nil)

	return Hud{txt: hudTxt, highScoreTxt: highScoreTxt, nextBlockTxt: nextBlockTxt, nextBlockPos: nextBlockPos, imd: imd}
}

func (h *Hud) DrawHUD(win *pixelgl.Window, game *tetris.Game, gameInfo tetris.Info, boxWidth float64, boxHeight float64) {
	h.txt.Clear()
	fmt.Fprintf(h.txt, "Score: %d\n", game.GetScore())
	fmt.Fprintf(h.txt, "Level: %d", gameInfo.Level)
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
}

func printHighScore(highScoreTxt *text.Text, game *tetris.Game) {
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

func drawNextBlock(nextBlockPos pixel.Vec, boxWidth float64, nextBlockTxt *text.Text, boxHeight float64, ns shape.Shape, nextImd *imdraw.IMDraw) {
	points := getShapePoints(
		nextBlockPos.Add(pixel.V(boxWidth, nextBlockTxt.LineHeight*1.5)),
		boxWidth, boxHeight, ns.GetBlocks())
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
