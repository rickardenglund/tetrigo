package tetris

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileWrite(t *testing.T) {
	const filename = "apa.dat"
	defer os.Remove(filename)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	require.NoError(t, err)

	_, err = fmt.Fprintf(f, "hej, %s\n", filename)
	require.NoError(t, err)

}
