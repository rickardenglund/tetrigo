package soundController

import (
	"Tetrigo/sound"
	"Tetrigo/tetris/shape"
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
