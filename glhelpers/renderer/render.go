package renderer

import (
	_ "embed" // nolint: golint
	"fmt"

	"github.com/rickardenglund/tetrigo/glhelpers/program"

	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Renderer struct {
	shader   program.Shader
	vao      uint32
	vbo      uint32
	win      *glfw.Window
	vertices []float32
}

const (
	windowWidth  = 800
	windowHeight = 600
)

func New() Renderer {
	r := Renderer{}
	r.initWin()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// nolint: gocritic // no
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	gl.ClearColor(0.5, 0.5, 0.5, 1)

	return r
}

func (r *Renderer) SetShader() {
	s, err := program.New(vertexShader, fragShader)
	if err != nil {
		panic(err)
	}

	r.shader = s

	s.Use()
}

func (r *Renderer) SetTriangle(x, y, size float32) {
	gl.GenVertexArrays(1, &r.vao)
	gl.GenBuffers(1, &r.vbo)
	gl.BindVertexArray(r.vao)

	w := size / 2
	r.vertices = []float32{
		x - w, y - w, 0,
		x + w, y - w, 0,
		x, y + w, 0,
	}

	const floatSize = 4

	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(r.vertices)*floatSize, gl.Ptr(r.vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
}

func (r *Renderer) Draw() {
	gl.ClearColor(0.5, 0.5, 0.5, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	r.shader.Use()
	gl.BindVertexArray(r.vao)

	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	r.win.SwapBuffers()
	glfw.PollEvents()
}

func (r *Renderer) initWin() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Super window", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("GL version: %v\n", version)

	glfw.SwapInterval(1)

	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	r.win = window
}

func (r *Renderer) ShouldClose() bool {
	return r.win.ShouldClose()
}

func (r *Renderer) Cleanup() {
	r.shader.Delete()
	gl.DeleteBuffers(1, &r.vbo)
	gl.DeleteVertexArrays(1, &r.vao)
	r.win.Destroy()
	glfw.Terminate()
}

//go:embed shaders/vertexShader.glsl
var vertexShader string

//go:embed shaders/fragmentShader.glsl
var fragShader string
