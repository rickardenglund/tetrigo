package main

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestName(t *testing.T) {
	size([]float64{1, 2, 3}[0])

	a := [...]uint64{2, 3, 4}
	size(a)
}

func size(v interface{}) {
	fmt.Printf("%T: %v s: %v\n", v, v, unsafe.Sizeof(v))
}
