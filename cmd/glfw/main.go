package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("key: %v, scancode: %v, action: %v, mode: %v\n", key, scancode, action, mods)
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(640, 480, "Super window", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("OpenGL version: %v\n", version)
	renderer := gl.GoStr(gl.GetString(gl.RENDERER))
	fmt.Printf("OpenGL renderer: %v\n", renderer)

	window.SetKeyCallback(keyCallback)
	glfw.SwapInterval(1)
	width, height := window.GetFramebufferSize()

	gl.Viewport(0, 0, int32(width), int32(height))

	gl.DepthFunc(gl.LESS)

	//program, err := newProgram(vertex_shader_text, fragment_shader_text)
	//if err != nil {
	//	panic(err)
	//}

	//var vertices = [3][5]float64{
	//	{-0.6, -0.4, 1., 0., 0.},
	//	{0.6, -0.4, 0., 1., 0.},
	//	{0., 0.6, 0., 0., 1.},
	//}
	//
	//mvpLocation := gl.GetUniformLocation(program, gl.Str("MVP\x00"))
	//vposLocation := gl.GetUniformLocation(program, gl.Str("vPos\x00"))
	//vCol := gl.GetUniformLocation(program, gl.Str("vCol\x00"))
	//
	//gl.EnableVertexAttribArray(uint32(vposLocation))
	//gl.VertexAttribPointer(uint32(vposLocation), 2, gl.FLOAT, gl.FALSE, unsafe.Sizeof(vertices[0]), gl.Ptr(0))
	//
	////gl.UseProgram(program)
	//
	//var vertexBuffer uint32
	//gl.GenBuffers(1, &vertexBuffer)
	//gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	//gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), gl.Ptr(vertices), gl.STATIC_DRAW)
	//
	//vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	//gl.ShaderSource(vertexShader, 1, gl.Ptr(vertexShader), len(vertexShader))

	gl.ClearColor(0, 1, 0, 0)

	//previousTime := glfw.GetTime()

	frames := 0
	seconds := time.Tick(time.Second)
	for !window.ShouldClose() {
		//now := glfw.GetTime()
		//elapsed := now - previousTime
		//previousTime = now

		select {
		case <-seconds:
			fmt.Printf("fps: %d\n", frames)
			frames = 0
		default:

		}
		//fmt.Printf("time: %v\n", elapsed)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		//time.Sleep(time.Millisecond * 500)

		window.SwapBuffers()
		glfw.PollEvents()
		frames++
	}

}

//func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
//	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
//	if err != nil {
//		return 0, err
//	}
//
//	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
//	if err != nil {
//		return 0, err
//	}
//
//	program := gl.CreateProgram()
//
//	gl.AttachShader(program, vertexShader)
//	gl.AttachShader(program, fragmentShader)
//	gl.LinkProgram(program)
//
//	var status int32
//	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
//	if status == gl.FALSE {
//		var logLength int32
//		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
//
//		log := strings.Repeat("\x00", int(logLength+1))
//		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
//
//		return 0, fmt.Errorf("failed to link program: %v", log)
//	}
//
//	gl.DeleteShader(vertexShader)
//	gl.DeleteShader(fragmentShader)
//
//	return program, nil
//}
//
//func compileShader(source string, shaderType uint32) (uint32, error) {
//	shader := gl.CreateShader(shaderType)
//
//	csources, free := gl.Strs(source)
//	gl.ShaderSource(shader, 1, csources, nil)
//	free()
//	gl.CompileShader(shader)
//
//	var status int32
//	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
//	if status == gl.FALSE {
//		var logLength int32
//		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
//
//		log := strings.Repeat("\x00", int(logLength+1))
//		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
//
//		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
//	}
//
//	return shader, nil
//}
//
//var vertex_shader_text = `
//#version 110
//uniform mat4 MVP;
//attribute vec3 vCol;
//attribute vec2 vPos;
//varying vec3 color;
//void main()
//{
//   gl_Position = MVP * vec4(vPos, 0.0, 1.0);
//   color = vCol;
//}`
//
//var fragment_shader_text = `
//#version 110
//varying vec3 color;
//void main()
//{
//    gl_FragColor = vec4(color, 1.0);
//}`
