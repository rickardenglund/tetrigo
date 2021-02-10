package tetris

import (
	"Tetrigo/tetris/shape"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

type Game struct {
	score        int
	nextTick     time.Time
	blocks       map[shape.Pos]bool
	activeBlocks shape.Shape
	gameOver     bool
	width        int
	height       int
	nextKind     int
	explodedRows int
	Results      []GameResult
	paused       bool
}

type Info struct {
	Level    int
	Width    int
	Height   int
	NextKind int
	Paused   bool
}

const maxLevel = 20
const highScoreFileName = "highscore.json"

func New() Game {
	g := Game{}
	g.nextTick = time.Now().Add(g.tickLength())
	g.blocks = map[shape.Pos]bool{}
	g.width = 10
	g.height = 20
	g.nextKind = rand.Int()
	g.newBlock()
	g.activeBlocks = shape.GetShape(g.nextKind, shape.Pos{X: 5, Y: g.height})

	if f, err := os.Open(highScoreFileName); err == nil {
		defer f.Close()

		var results []GameResult
		err = json.NewDecoder(f).Decode(&results)
		if err != nil {
			fmt.Printf("Failed to read highscores: %v\n", err)
		}
		g.Results = results
	}
	return g
}

func (g *Game) Speed() {
	g.nextTick = time.Now().Add(-time.Millisecond)
}

func (g *Game) Right() {
	for i := range g.activeBlocks.GetBlocks() {
		if g.activeBlocks.GetBlocks()[i].X >= g.width-1 || g.collides(shape.Pos{X: g.activeBlocks.GetBlocks()[i].X + 1, Y: g.activeBlocks.GetBlocks()[i].Y}) {
			return
		}
	}

	g.activeBlocks.Right()
}

func (g *Game) Left() {
	for i := range g.activeBlocks.GetBlocks() {
		if g.activeBlocks.GetBlocks()[i].X <= 0 || g.collides(shape.Pos{X: g.activeBlocks.GetBlocks()[i].X - 1, Y: g.activeBlocks.GetBlocks()[i].Y}) {
			return
		}
	}

	g.activeBlocks.Left()
}

func (g *Game) Rotate() {
	rotated := g.activeBlocks.Rotated()
	for i := range rotated {
		if g.collides(rotated[i]) {
			return
		}
	}

	g.activeBlocks.Rotate()
}

func (g *Game) Tick(currentTime time.Time) {
	if g.paused {
		g.nextTick = time.Now().Add(g.tickLength())
		return
	}
	if !currentTime.After(g.nextTick) || g.gameOver {
		return
	}

	isBlocked := false
	for i := range g.activeBlocks.GetBlocks() {
		if g.collides(shape.Pos{X: g.activeBlocks.GetBlocks()[i].X, Y: g.activeBlocks.GetBlocks()[i].Y - 1}) {
			isBlocked = true
			break
		}
	}

	if isBlocked {
		for i := range g.activeBlocks.GetBlocks() {
			g.blocks[g.activeBlocks.GetBlocks()[i]] = true
		}

		for i := range g.activeBlocks.GetBlocks() {
			if g.activeBlocks.GetBlocks()[i].Y > g.height {
				g.setGameOver()
				return
			}
		}

		g.newBlock()
		g.score += 1
		g.checkForFullLines()
	}

	if !isBlocked {
		g.activeBlocks.Down()
	}

	g.nextTick = g.nextTick.Add(g.tickLength())
}

func (g *Game) GetScore() int {
	return g.score
}

func (g *Game) GetBlocks() []shape.Pos {
	res := make([]shape.Pos, 0, len(g.blocks)+len(g.activeBlocks.GetBlocks()))
	for k, exists := range g.blocks {
		if exists {
			res = append(res, k)
		}
	}
	return append(res, g.activeBlocks.GetBlocks()...)
}

func (g *Game) collides(currentBlock shape.Pos) bool {
	if currentBlock.Y < 0 {
		return true
	}

	if currentBlock.X >= g.width || currentBlock.X < 0 {
		return true
	}

	_, exists := g.blocks[currentBlock]
	return exists
}

func (g *Game) IsGameOver() bool {
	return g.gameOver
}

func (g *Game) GetInfo() Info {
	i := Info{
		Height:   g.height,
		Width:    g.width,
		Level:    g.Level(),
		NextKind: g.nextKind,
		Paused:   g.paused,
	}

	return i
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

	var fullRows []int
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
	for i := len(fullRows) - 1; i >= 0; i-- {
		fullRow := fullRows[i]

		for y := fullRow; y < g.height; y++ {
			for x := 0; x < g.width; x++ {
				p := shape.Pos{X: x, Y: y}
				if g.blocks[p] {
					delete(g.blocks, p)
					g.blocks[shape.Pos{X: p.X, Y: p.Y - 1}] = true
				}
			}

		}
	}

	g.score += len(fullRows) * len(fullRows) * g.width

	for i := range fullRows {
		fmt.Printf("%d, ", fullRows[i])
	}

	g.explodedRows += len(fullRows)
}

func (g *Game) newBlock() {
	g.activeBlocks = shape.GetShape(g.nextKind, shape.Pos{X: 5, Y: g.height})
	g.nextKind = rand.Int()
}

func (g *Game) NextBlock() shape.Shape {
	return shape.GetShape(g.nextKind, shape.Pos{})
}

func (g *Game) Level() int {
	l := (g.explodedRows / 20) + 1
	if l > maxLevel {
		return maxLevel
	}

	return l
}
func (g *Game) tickLength() time.Duration {
	tickLength := 600*time.Millisecond - time.Duration(g.Level())*time.Millisecond*25
	return tickLength
}

type GameResult struct {
	Score int
	Level int
}

type ByScore []GameResult

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Less(i, j int) bool { return a[i].Score > a[j].Score }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (g *Game) setGameOver() {
	g.gameOver = true

	result := GameResult{
		Score: g.score,
		Level: g.Level(),
	}
	g.Results = append(g.Results, result)
	sort.Sort(ByScore(g.Results))
	if len(g.Results) > 4 {
		g.Results = g.Results[:5]
	}

	f, err := os.OpenFile(highScoreFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(g.Results)
	if err != nil {
		panic(err)
	}
}

func (g *Game) TogglePaus() {
	g.paused = !g.paused
}
