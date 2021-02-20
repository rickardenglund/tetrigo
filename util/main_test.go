package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapRange(t *testing.T) {
	res := MapRange(10, 0, 10, 10, 20)
	require.Equal(t, res, 20.0)

	res = MapRange(5, 0, 10, 10, 20)
	require.Equal(t, res, 15.0)
}
