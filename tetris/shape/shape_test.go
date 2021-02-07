package shape

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShape_Rotate(t *testing.T) {
	//	 #
	//	###
	s := Shape{
		pos:      Pos{5, 5},
		kind:     0,
		rotation: 0,
	}
	require.ElementsMatchf(t, s.GetBlocks(), []Pos{{4, 5}, {5, 5}, {6, 5}, {5, 6}}, "content does not match")

	s.Rotate()
	require.ElementsMatchf(t, s.GetBlocks(), []Pos{{4, 5}, {5, 5}, {5, 4}, {5, 6}}, "content does not match")

	s.Rotate()

	require.ElementsMatchf(t, s.GetBlocks(), []Pos{{4, 5}, {5, 5}, {6, 5}, {5, 4}}, "content does not match")

	//	#
	//	#
	//	#
	//	#
}
