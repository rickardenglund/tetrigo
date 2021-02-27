package main

import (
	"runtime"

	"Tetrigo/glhelpers/renderer"
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
	r.SetTriangle(0, 0, 1)

	for !r.ShouldClose() {
		r.Draw()
	}
}
