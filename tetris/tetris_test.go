package tetris

import (
	"fmt"
	"testing"
)

func TestRowScore(t *testing.T) {
	const w = 10
	for level := 1; level <= 20; level++ {
		for rows := 1; rows < 5; rows++ {
			fmt.Printf("L: %d, rows: %d - score: %d\n", level, rows, rowScore(rows, w, level))
		}
		println()
	}

}
