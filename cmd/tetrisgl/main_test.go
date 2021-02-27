package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func Test_initWin(t *testing.T) {
	tests := []struct {
		name string
		want *glfw.Window
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initWin(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initWin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlaskjd(t *testing.T) {
	fmt.Printf("bytes: %v\n", "")
}
