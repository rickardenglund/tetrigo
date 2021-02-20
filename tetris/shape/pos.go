package shape

type Pos struct {
	X, Y int
}

func (p *Pos) Add(other Pos) Pos {
	return Pos{p.X + other.X, p.Y + other.Y}
}
