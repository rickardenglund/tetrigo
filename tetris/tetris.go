package tetris

import (
	"fmt"
	"math/rand"
	"sort"
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

		g.activeBlocks = g.newShape()
		g.score += 1
		g.checkForFullLines()
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
	res := make([]Pos, 0, len(g.blocks)+len(g.activeBlocks))
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
	shapes := [][]Pos{
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
		{
			{5, g.height},
			{5, g.height + 1},
			{5, g.height + 2},
			{5, g.height + 3},
		},
	}

	return shapes[rand.Intn(len(shapes))]
}

func (g *Game) GetDimensions() (int, int) {
	return g.width, g.height
}

func (g *Game) checkForFullLines() {
	rows := map[int]map[int]bool{}
	for block, exists := range g.blocks {
		if !exists {
			continue
		}

		if _, exists := rows[block.Y]; !exists {
			rows[block.Y] = map[int]bool{}
		}

		rows[block.Y][block.X] = true
	}

	fullRows := []int{}
	for row := range rows {
		if len(rows[row]) == g.width {
			fullRows = append(fullRows, row)
		}
	}

	for i := len(fullRows) - 1; i >= 0; i-- {
		for b := range g.blocks {
			if b.Y == fullRows[i] {
				delete(g.blocks, b)
			}
		}
	}

	sort.Ints(fullRows)
	for _, fullRow := range fullRows {
		for y := fullRow; y < g.height; y++ {
			for x := 0; x < g.width; x++ {
				p := Pos{x, y}
				if g.blocks[p] {
					delete(g.blocks, p)
					g.blocks[Pos{p.X, p.Y-1}] = true
				}
			}

		}
	}

	g.score += len(fullRows) * len(fullRows) * g.width

	for i := range fullRows {
		fmt.Printf("%d, ", fullRows[i])
	}
	println()
}
