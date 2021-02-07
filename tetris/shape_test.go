package tetris

import (
	"Tetrigo/tetris/shape"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShape_Rotate(t *testing.T) {
//	 #
//	###
	s := shape.Shape{
		pos:      shape.Pos{5,5},
		kind:     0,
		rotation: 0,
	}
	require.ElementsMatchf(t, s.getBlocks(), []shape.Pos{{4,5},{5,5},{6,5},{5,6}}, "content does not match")

	s.Rotate()
	require.ElementsMatchf(t, s.getBlocks(), []shape.Pos{{4,5},{5,5},{5,4},{5,6}}, "content does not match")

	s.Rotate()

	require.ElementsMatchf(t, s.getBlocks(), []shape.Pos{{4,5},{5,5},{6,5},{5,4}}, "content does not match")

//	#
//	#
//	#
//	#
}
