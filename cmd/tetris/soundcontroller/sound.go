package soundcontroller

import (
	"github.com/rickardenglund/tetrigo/sound"
	"github.com/rickardenglund/tetrigo/tetris/shape"
)

func ControlSound(newBlock bool, explodedBlocks []shape.Block, sounds *sound.Sound) {
	if len(explodedBlocks) > 0 {
		sounds.Click()
		println("click")
	}

	if newBlock {
		sounds.Tick()
		println("tick")
	}
}
