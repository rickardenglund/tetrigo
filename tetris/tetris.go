package tetris

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/rickardenglund/tetrigo/tetris/shape"
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
	startTime    time.Time
	stopTime     time.Time
	pausStart    time.Time
	totalPause   time.Duration
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
	g.startTime = time.Now()
	g.nextTick = time.Now().Add(g.tickLength())
	g.blocks = map[shape.Pos]shape.Block{}
	g.width = 10
	g.height = 20
	g.nextKind = rand.Int() //nolint: gosec // ignore weak rand
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
		if g.activeShape.GetBlocks()[i].Pos.X >= g.width-1 ||
			g.collides(shape.Pos{X: g.activeShape.GetBlocks()[i].Pos.X + 1, Y: g.activeShape.GetBlocks()[i].Pos.Y}) {
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
		if g.activeShape.GetBlocks()[i].Pos.X <= 0 ||
			g.collides(shape.Pos{X: g.activeShape.GetBlocks()[i].Pos.X - 1, Y: g.activeShape.GetBlocks()[i].Pos.Y}) {
			return
		}
	}

	g.activeShape.Left()
}

func (g *Game) Rotate() {
	if g.gameOver {
		return
	}

	rotated := g.activeShape.Rotated()
	for i := range rotated {
		if g.collides(rotated[i].Pos) {
			return
		}
	}

	g.activeShape.Rotate()
}

func (g *Game) Tick(currentTime time.Time) []shape.Block {
	if g.paused {
		g.nextTick = time.Now().Add(g.tickLength())
		return nil
	}

	if !currentTime.After(g.nextTick) || g.gameOver {
		return nil
	}

	var (
		isBlocked      = g.activeIsBlocked()
		explodedBlocks []shape.Block
	)

	if isBlocked {
		for i := range g.activeShape.GetBlocks() {
			b := g.activeShape.GetBlocks()[i]
			g.blocks[b.Pos] = b
		}

		g.newBlock()
		g.score += g.Level()
		explodedBlocks = g.checkForFullLines()
	}

	if !isBlocked {
		g.activeShape.Down()
	}

	g.nextTick = g.nextTick.Add(g.tickLength())
	g.activeShape.Age++

	return explodedBlocks
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
	g.stopTime = time.Now()
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

func (g *Game) checkForFullLines() []shape.Block {
	var (
		rows           = map[int]map[int]bool{}
		explodedBlocks []shape.Block
	)

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

			for x := 0; x < g.width; x++ {
				explodedBlocks = append(explodedBlocks, g.blocks[shape.Pos{X: x, Y: row}])
			}
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

	return explodedBlocks
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

	g.nextKind = rand.Int() //nolint: gosec // ignore weak rand
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
	Score     int
	Level     int
	PlayTime  time.Duration
	StartTime time.Time
	StopTime  time.Time
}

type ByScore []GameResult

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Less(i, j int) bool { return a[i].Score > a[j].Score }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (g *Game) setGameOver() {
	g.stopTime = time.Now()
	g.gameOver = true

	result := GameResult{
		Score:     g.score,
		Level:     g.Level(),
		PlayTime:  g.stopTime.Sub(g.startTime) - g.totalPause,
		StartTime: g.startTime,
		StopTime:  g.stopTime,
	}
	fmt.Printf("playtime: %v\n", result.PlayTime)

	g.Results = append(g.Results, result)
	sort.Sort(ByScore(g.Results))

	const resultsToStore = 5
	if len(g.Results) > resultsToStore-1 {
		g.Results = g.Results[:resultsToStore]
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

func (g *Game) TogglePause() {
	if !g.paused {
		g.pausStart = time.Now()
	} else {
		g.totalPause += time.Since(g.pausStart)
	}

	g.paused = !g.paused
}
