package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFit(t *testing.T) {
	require.Equal(t, float32(0), fit(400, 0, 800, -2, 2))
}
