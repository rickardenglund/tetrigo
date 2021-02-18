package shape

type Shape struct {
	pos      Pos
	kind     int
	rotation int
	Age      int
}

type Block struct {
	Pos  Pos
	Kind int
}

func (s *Shape) GetBlocks() []Block {
	return s.getBlocks(s.rotation)
}

func (s *Shape) getBlocks(rotation int) []Block {
	blocks := make([]Block, 0, len(tetronimos[s.kind]))
	for _, p := range tetronimos[s.kind][rotation] {
		newPos := p.Add(s.pos)
		blocks = append(blocks, Block{Pos: newPos, Kind: s.kind})
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

func (s *Shape) Rotated() []Block {
	return s.getBlocks((s.rotation + 1) % 4)
}

func GetShape(kind int, pos Pos) Shape {
	s := Shape{
		pos:  pos,
		kind: kind % len(tetronimos),
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
		{0, 1},
		{-1, 1},
		{1, 1},
		{0, 2},
	},
	{
		{-1, 1},
		{0, 0},
		{0, 1},
		{0, 2},
	},
	{
		{0, 1},
		{-1, 1},
		{1, 1},
		{0, 0},
	},
	{
		{1, 1},
		{0, 0},
		{0, 1},
		{0, 2},
	},
}

var iTetronimo = [][]Pos{
	{
		{0, 0},
		{0, 1},
		{0, 2},
		{0, 3},
	},
	{
		{-2, 2},
		{-1, 2},
		{0, 2},
		{1, 2},
	},
	{
		{0, 0},
		{0, 1},
		{0, 2},
		{0, 3},
	},
	{
		{-2, 2},
		{-1, 2},
		{0, 2},
		{1, 2},
	},
}

var oTetronimo = [][]Pos{
	{
		{0, 0},
		{1, 0},
		{0, 1},
		{1, 1},
	},
	{
		{0, 0},
		{1, 0},
		{0, 1},
		{1, 1},
	},
	{
		{0, 0},
		{1, 0},
		{0, 1},
		{1, 1},
	},
	{
		{0, 0},
		{1, 0},
		{0, 1},
		{1, 1},
	},
}
var jTetronimo = [][]Pos{
	{
		{-1, 0},
		{-1, 1},
		{0, 0},
		{1, 0},
	},
	{
		{1, 1},
		{0, 1},
		{0, 0},
		{0, -1},
	},
	{
		{-1, 0},
		{0, 0},
		{1, 0},
		{1, -1},
	},
	{
		{-1, -1},
		{0, -1},
		{0, 0},
		{0, 1},
	},
}
var lTetronimo = [][]Pos{
	{
		{-1, 0},
		{0, 0},
		{1, 0},
		{1, 1},
	},
	{
		{1, -1},
		{0, -1},
		{0, 1},
		{0, 0},
	},
	{
		{-1, 0},
		{0, 0},
		{1, 0},
		{-1, -1},
	},
	{
		{-1, 1},
		{0, 1},
		{0, 0},
		{0, -1},
	},
}

var sTetronimo = [][]Pos{
	{
		{0, 0},
		{-1, 0},
		{0, 1},
		{1, 1},
	},
	{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, -1},
	},
	{
		{0, 0},
		{-1, 0},
		{0, 1},
		{1, 1},
	},
	{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, -1},
	},
}

var zTetronimo = [][]Pos{
	{
		{0, 1},
		{1, 1},
		{1, 2},
		{0, 0},
	},
	{
		{0, 1},
		{0, 2},
		{-1, 2},
		{1, 1},
	},
	{
		{0, 1},
		{1, 1},
		{1, 2},
		{0, 0},
	},
	{
		{0, 1},
		{0, 2},
		{-1, 2},
		{1, 1},
	},
}
