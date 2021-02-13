package hud

import (
	"Tetrigo/tetris"
	"fmt"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"

	"github.com/faiface/pixel/text"
)

type Hud struct {
	txt          *text.Text
	highScoreTxt *text.Text
}

func New(pos pixel.Vec, atlas *text.Atlas, highScorePos pixel.Vec) Hud {
	hudTxt := text.New(pos, atlas)
	highScoreTxt := text.New(highScorePos, atlas)

	return Hud{txt: hudTxt, highScoreTxt: highScoreTxt}
}

func (h *Hud) DrawHUD(win *pixelgl.Window, game tetris.Game, gameInfo tetris.Info) {
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
