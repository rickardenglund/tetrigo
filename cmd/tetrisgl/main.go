package main

import (
	"runtime"

	"github.com/rickardenglund/tetrigo/glhelpers/renderer"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {

	//game := tetris.New()

	r := renderer.New()
	defer r.Cleanup()

	r.SetShader()
	r.SetTriangle(0, 0, 2)

	for !r.ShouldClose() {
		r.Draw()
	}
}
