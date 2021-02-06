package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMapRange(t *testing.T) {
	res := mapRange(10, 0, 10, 10, 20)
	require.Equal(t, res, 20.0)

	res = mapRange(5, 0, 10, 10, 20)
	require.Equal(t, res, 15.0)
}
