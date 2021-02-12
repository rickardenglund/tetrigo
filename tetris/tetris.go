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
	blocks       map[shape.Pos]shape.Block
	activeShape  shape.Shape
	gameOver     bool
	width        int
	height       int
	nextKind     int
	explodedRows int
	Results      []GameResult
	paused       bool
}

type Info struct {
	Level     int
	Width     int
	Height    int
	NextKind  int
	Paused    bool
	ActiveAge int
}

const maxLevel = 20
const highScoreFileName = "highscore.json"

func New() Game {
	g := Game{}
	g.nextTick = time.Now().Add(g.tickLength())
	g.blocks = map[shape.Pos]shape.Block{}
	g.width = 10
	g.height = 20
	g.nextKind = rand.Int()
	g.newBlock()

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
	if g.paused || g.gameOver {
		return
	}

	for i := range g.activeShape.GetBlocks() {
		if g.activeShape.GetBlocks()[i].Pos.X >= g.width-1 || g.collides(shape.Pos{X: g.activeShape.GetBlocks()[i].Pos.X + 1, Y: g.activeShape.GetBlocks()[i].Pos.Y}) {
			return
		}
	}

	g.activeShape.Right()
}

func (g *Game) Left() {
	if g.paused || g.gameOver {
		return
	}

	for i := range g.activeShape.GetBlocks() {
		if g.activeShape.GetBlocks()[i].Pos.X <= 0 || g.collides(shape.Pos{X: g.activeShape.GetBlocks()[i].Pos.X - 1, Y: g.activeShape.GetBlocks()[i].Pos.Y}) {
			return
		}
	}

	g.activeShape.Left()
}

func (g *Game) Rotate() {
	rotated := g.activeShape.Rotated()
	for i := range rotated {
		if g.collides(rotated[i].Pos) {
			return
		}
	}

	g.activeShape.Rotate()
}

func (g *Game) Tick(currentTime time.Time) {
	if g.paused {
		g.nextTick = time.Now().Add(g.tickLength())
		return
	}
	if !currentTime.After(g.nextTick) || g.gameOver {
		return
	}

	isBlocked := g.activeIsBlocked()

	if isBlocked {
		for i := range g.activeShape.GetBlocks() {
			b := g.activeShape.GetBlocks()[i]
			g.blocks[b.Pos] = b
		}

		g.newBlock()
		g.score += g.Level()
		g.checkForFullLines()
	}

	if !isBlocked {
		g.activeShape.Down()
	}

	g.nextTick = g.nextTick.Add(g.tickLength())
	g.activeShape.Age++
}

func (g *Game) activeIsBlocked() bool {
	isBlocked := false
	for i := range g.activeShape.GetBlocks() {
		if g.collides(shape.Pos{X: g.activeShape.GetBlocks()[i].Pos.X, Y: g.activeShape.GetBlocks()[i].Pos.Y - 1}) {
			isBlocked = true
			break
		}
	}
	return isBlocked
}

func (g *Game) GetScore() int {
	return g.score
}

func (g *Game) GetBlocks() []shape.Block {
	res := make([]shape.Block, 0, len(g.blocks)+len(g.activeShape.GetBlocks()))
	for _, block := range g.blocks {
		res = append(res, block)
	}
	return append(res, g.activeShape.GetBlocks()...)
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
		Height:    g.height,
		Width:     g.width,
		Level:     g.Level(),
		NextKind:  g.nextKind,
		Paused:    g.paused,
		ActiveAge: g.activeShape.Age,
	}

	return i
}

func (g *Game) checkForFullLines() {
	rows := map[int]map[int]bool{}
	for pos := range g.blocks {
		if _, exists := rows[pos.Y]; !exists {
			rows[pos.Y] = map[int]bool{}
		}

		rows[pos.Y][pos.X] = true
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
				b, exists := g.blocks[p]
				if exists {
					delete(g.blocks, p)
					newPos := shape.Pos{X: p.X, Y: p.Y - 1}
					g.blocks[newPos] = shape.Block{Pos: newPos, Kind: b.Kind}
				}
			}

		}
	}

	g.score += rowScore(len(fullRows), g.width, g.Level())

	g.explodedRows += len(fullRows)
	fmt.Printf("nRows: %d\n", g.explodedRows)
}

func rowScore(nRows, width, level int) int {
	s := nRows * nRows * width * level
	return s
}

func (g *Game) newBlock() {
	g.activeShape = shape.GetShape(g.nextKind, shape.Pos{X: g.width / 2, Y: g.height})
	if g.activeIsBlocked() {
		g.setGameOver()
	}

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
