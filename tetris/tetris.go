package tetris

import (
	"math/rand"
	"time"
)

type Game struct {
	score        int
	nextTick     time.Time
	blocks       map[Pos]bool
	activeBlocks []Pos
	gameOver     bool
	width        int
	height       int
}
type Pos struct {
	X, Y int
}

const tickLength = time.Millisecond * 200

func New() Game {
	g := Game{nextTick: time.Now().Add(tickLength)}
	g.blocks = map[Pos]bool{}
	g.width = 10
	g.height = 20
	g.activeBlocks = g.newShape()
	return g
}

func (g *Game) Speed() {
	g.nextTick = time.Now().Add(-time.Millisecond)
}

func (g *Game) Right() {
	for i := range g.activeBlocks {
		if g.activeBlocks[i].X >= g.width-1 || g.collides(Pos{g.activeBlocks[i].X + 1, g.activeBlocks[i].Y}) {
			return
		}
	}

	for i := range g.activeBlocks {
		g.activeBlocks[i].X++
	}
}

func (g *Game) Left() {
	for i := range g.activeBlocks {
		if g.activeBlocks[i].X <= 0 || g.collides(Pos{g.activeBlocks[i].X - 1, g.activeBlocks[i].Y}) {
			return
		}
	}

	for i := range g.activeBlocks {
		g.activeBlocks[i].X--
	}
}

func (g *Game) Tick(currentTime time.Time) {
	if !currentTime.After(g.nextTick) || g.gameOver {
		return
	}

	isBlocked := false
	for i := range g.activeBlocks {
		if g.collides(Pos{g.activeBlocks[i].X, g.activeBlocks[i].Y - 1}) {
			isBlocked = true
			break
		}
	}

	if isBlocked {
		for i := range g.activeBlocks {
			g.blocks[g.activeBlocks[i]] = true
		}

		for i := range g.activeBlocks {
			if g.activeBlocks[i].Y > g.height {
				g.gameOver = true
				return
			}
		}

		g.score++
		g.activeBlocks = g.newShape()
	}

	if !isBlocked {
		for i := range g.activeBlocks {
			g.activeBlocks[i].Y--
		}
	}

	g.nextTick = g.nextTick.Add(tickLength)
}

func (g *Game) GetScore() int {
	return g.score
}

func (g *Game) GetBlocks() []Pos {
	res := make([]Pos, 0, len(g.blocks) + len(g.activeBlocks))
	for k, exists := range g.blocks {
		if exists {
			res = append(res, k)
		}
	}
	return append(res, g.activeBlocks...)
}

func (g *Game) collides(currentBlock Pos) bool {
	if currentBlock.Y < 0 {
		return true
	}

	_, exists := g.blocks[currentBlock]
	return exists
}

func (g *Game) IsGameOver() bool {
	return g.gameOver
}

func (g *Game) newShape() []Pos {
	shapes := [][]Pos {
		{
			{5, g.height},
			{6, g.height},
			{4, g.height},
			{5, g.height + 1},
		},
		{
			{5, g.height},
			{6, g.height},
			{5, g.height + 1},
			{6, g.height + 1},
		},

	}

	return shapes[rand.Intn(len(shapes))]
}

func (g *Game) GetDimensions() (int, int) {
	return g.width, g.height
}
