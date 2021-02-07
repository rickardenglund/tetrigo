package shape

import (
	"math/rand"
)

type Shape struct {
	pos      Pos
	blocks   []Pos
	kind     int
	rotation int
}

func (s *Shape) GetBlocks() []Pos {
	return s.getBlocks(s.rotation)
}

func (s *Shape) getBlocks(rotation int) []Pos {
	blocks := make([]Pos, 0, len(tetronimos[s.kind]))
	for _, p := range tetronimos[s.kind][rotation] {
		newPos := p.Add(s.pos)
		blocks = append(blocks, newPos)
	}

	return blocks
}

func (s *Shape) Right() {
	s.pos.X++
}

func (s *Shape) Left() {
	s.pos.X--
}

func (s *Shape) Down() {
	s.pos.Y--
}

func (s *Shape) Rotate() {
	s.rotation = (s.rotation + 1) % 4
}

func (s *Shape) Rotated() []Pos {
	return s.getBlocks((s.rotation + 1) % 4)
}

func NewShape(pos Pos) Shape {
	s := Shape{
		pos:  pos,
		kind: rand.Intn(len(tetronimos)),
	}

	return s
}

var tetronimos = [][][]Pos{
	tTetronimo,
	iTetronimo,
	oTetronimo,
	jTetronimo,
	lTetronimo,
	sTetronimo,
	zTetronimo,
}

//  #
// ###
var tTetronimo = [][]Pos{
	{
		{0, 0},
		{-1, 0},
		{1, 0},
		{0, 1},
	},
	{
		{-1, 0},
		{0, -1},
		{0, 0},
		{0, 1},
	},
	{
		{0, 0},
		{-1, 0},
		{1, 0},
		{0, -1},
	},
	{
		{1, 0},
		{0, -1},
		{0, 0},
		{0, 1},
	},
}

var iTetronimo = [][]Pos {
	{
		{0,-2},
		{0,-1},
		{0,0},
		{0,1},
	},
	{
		{-2,0},
		{-1,0},
		{0,0},
		{1,0},
	},
	{
		{0,-2},
		{0,-1},
		{0,0},
		{0,1},
	},
	{
		{-2,0},
		{-1,0},
		{0,0},
		{1,0},
	},
}

var oTetronimo = [][]Pos {
	{
		{0,0},
		{1,0},
		{0,1},
		{1,1},
	},
	{
		{0,0},
		{1,0},
		{0,1},
		{1,1},
	},
	{
		{0,0},
		{1,0},
		{0,1},
		{1,1},
	},
	{
		{0,0},
		{1,0},
		{0,1},
		{1,1},
	},
}
var jTetronimo = [][]Pos {
	{
		{-1,0},
		{-1,1},
		{0,0},
		{1,0},
	},
	{
		{1,1},
		{0,1},
		{0,0},
		{0,-1},
	},
	{
		{-1,0},
		{0,0},
		{1,0},
		{1,-1},
	},
	{
		{-1,-1},
		{0,-1},
		{0,0},
		{0,1},
	},
}
var lTetronimo = [][]Pos {
	{
		{-1,0},
		{0,0},
		{1,0},
		{1,1},
	},
	{
		{1,-1},
		{0,-1},
		{0,1},
		{0,0},
	},
	{
		{-1,0},
		{0,0},
		{1,0},
		{-1,-1},
	},
	{
		{-1,1},
		{0,1},
		{0,0},
		{0,-1},
	},
}

var sTetronimo = [][]Pos {
	{
		{0,0},
		{-1,0},
		{0,1},
		{1,1},
	},
	{
		{0,0},
		{0,1},
		{1,0},
		{1, -1},
	},
	{
		{0,0},
		{-1,0},
		{0,1},
		{1,1},
	},
	{
		{0,0},
		{0,1},
		{1,0},
		{1, -1},
	},
}

var zTetronimo = [][]Pos {
	{
		{0,0},
		{1,0},
		{1,1},
		{0,-1},
	},
	{
		{0,0},
		{0,1},
		{-1,1},
		{1,0},
	},
	{
		{0,0},
		{1,0},
		{1,1},
		{0,-1},
	},
	{
		{0,0},
		{0,1},
		{-1,1},
		{1,0},
	},
}
