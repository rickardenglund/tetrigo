package main

import (
	"fmt"
	"runtime"

	"github.com/rickardenglund/tetrigo/glhelpers/arraybuffer"
	"github.com/rickardenglund/tetrigo/glhelpers/camera"
	"github.com/rickardenglund/tetrigo/glhelpers/program"
	"github.com/rickardenglund/tetrigo/glhelpers/textures"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	win, err := initWindow()
	camState := camera.New(win)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("GL version: %v\n", version)

	//s1, err := program.New(vertexshader, fragshader)
	//if err != nil {
	//	panic(err)
	//}

	s2, err := program.New(vertexshader, fragshaderyellow)
	if err != nil {
		panic(err)
	}
	defer s2.Delete()
	s2.Use()
	s2.SetUniform1i("ourTexture", 0)
	s2.SetUniform1i("ourTexture2", 1)

	//projection := mgl32.Ortho(-1, 1, -1, 1, -10, 100)
	projection := mgl32.Perspective(mgl32.DegToRad(45), float32(windowWidth)/windowHeight, 0.1, 100)
	projectionUniform := gl.GetUniformLocation(s2.Id, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	//model = model.Mul4(mgl32.Translate3D(0.5, 0, 0))
	//model = model.Mul4(mgl32.Scale3D(140, 140, 0))
	//model = model.Mul4(mgl32.HomogRotate3DX(mgl32.DegToRad(-45)))
	//model = model.Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(-45)))

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(s2.Id, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	glfw.SwapInterval(1)

	triangleVertices := []float32{
		// pos              clr             texture
		// FRONT
		-0.5, -0.5, 0.5, 1, 0, 0, 0, 0,
		0.5, -0.5, 0.5, 0, 1, 0, 1, 0,
		-0.5, 0.5, 0.5, 0, 0, 1, 0, -1,
		0.5, 0.5, 0.5, 0, 1, 1, 1, -1,

		// BACK
		-0.5, -0.5, -0.5, 1, 0, 0, 0, 0,
		0.5, -0.5, -0.5, 0, 1, 0, 1, 0,
		-0.5, 0.5, -0.5, 0, 0, 1, 0, -1,
		0.5, 0.5, -0.5, 0, 1, 1, 1, -1,
	}

	indices := []int32{
		//FRONT
		0, 1, 2,
		1, 2, 3,
		//BACK
		4, 5, 6,
		6, 5, 7,
		//RIGHT
		5, 3, 1,
		5, 7, 3,
		//LEFT
		0, 4, 2,
		2, 4, 6,
		//TOP
		2, 3, 6,
		6, 3, 7,
		//BOTTOM
		0, 1, 4,
		4, 1, 5,
	}

	vbo := arraybuffer.New(triangleVertices)
	defer vbo.Delete()

	var vao1 uint32
	gl.GenVertexArrays(1, &vao1)

	// vao1
	gl.BindVertexArray(vao1)
	s2.Use()

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.Enable(gl.DEPTH_TEST)
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.5, 0.5, 0.5, 0)

	t, err := textures.New("/Users/rickardenglund/code/tetrigo/cmd/learnopengl/butterfly.png")
	if err != nil {
		panic(err)
	}

	t2, err := textures.New("/Users/rickardenglund/code/tetrigo/cmd/learnopengl/butterfly.png")
	if err != nil {
		panic(err)
	}

	cubePoss := [][]float32{
		{0, 0, 0},
		{2.5, 1, -15},
		{-1.5, -2.2, -2.5},
		{-3.8, -3, -12.3},
		{-1.3, 1, -1.5},
		{-1.5, -2.2, -2.5},
		{-3.8, -2, -12.3},
	}
	//glm::vec3( 2.4f, -0.4f, -3.5f),
	//glm::vec3(-1.7f,  3.0f, -7.5f),
	//glm::vec3( 1.3f, -2.0f, -2.5f),
	//glm::vec3( 1.5f,  2.0f, -2.5f),
	//glm::vec3( 1.5f,  0.2f, -1.5f),
	//glm::vec3(-1.3f,  1.0f, -1.5f)

	v := float64(0)
	for !win.ShouldClose() {
		v += 0.01
		x, y := win.GetCursorPos()

		s2.Use()
		s2.SetUniform2f("mouse", float32(x/windowWidth)+0.5, float32(-y/windowHeight))
		t.Bind(gl.TEXTURE0)
		t2.Bind(gl.TEXTURE1)
		//model = mgl32.Translate3D(float32(x), windowHeight-float32(y), 0)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		for _, kp := range cubePoss {
			model = mgl32.Ident4()
			model = model.Mul4(mgl32.Translate3D(kp[0], kp[1], kp[2]))
			//model = model.Mul4(mgl32.Scale3D(140, 140, 0))
			//model = model.Mul4(mgl32.HomogRotate3DX(mgl32.DegToRad(-45)))
			model = model.Mul4(mgl32.HomogRotate3DY(float32(v)))
			gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

			gl.DrawElements(gl.TRIANGLES, int32(len(triangleVertices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
		}

		eye, centre, up := camState.Update()
		fmt.Printf("eye: %v, centre: %v, up: %v\n", eye, centre, up)
		camera := mgl32.LookAtV(eye, centre, up) //mgl32.Translate3D(0.5, -0.5, 0)
		cameraUniform := gl.GetUniformLocation(s2.Id, gl.Str("camera\x00"))
		gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func fit(v, smin, smax, tmin, tmax float32) float32 {
	percent := ((v - smin) / (smax - smin))

	targetdiff := tmax - tmin
	ret := percent*targetdiff + tmin
	return ret

}

func initWindow() (*glfw.Window, error) {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		panic(err)
	}
	return window, err
}

var vertexshader = `
#version 330 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 vColour;
layout (location = 2) in vec2 aTexCoord;
out vec4 vertexColor;
out vec2 pos;
out vec2 texCoord;
uniform vec2 mouse;

uniform mat4 camera;
uniform mat4 projection;
uniform mat4 model;

void main()
{
	gl_Position = projection* camera * model*vec4(aPos, 1.0);
	vertexColor = vec4(vColour, 1);
	pos = aPos.xy;
	texCoord = aTexCoord;
}` + "\x00"

var fragshader = `
#version 330 core
out vec4 FragColor;
in vec4 vertexColor;
in vec2 pos;
in vec2 texCoord;

uniform sampler2D ourTexture;

void main()
{
//    FragColor = vec4(1.0f, 0.0f, 0.0f, 1.0f);
//    FragColor = vertexColor;
//	FragColor = vec4(texCoord.xy, 0,1);//texture(ourTexture, texCoord);
	FragColor = texture(ourTexture, texCoord) * vertexColor;
} ` + "\x00"

var fragshaderyellow = `
#version 330 core
out vec4 FragColor;
in vec2 texCoord;
in vec4 vertexColor;
uniform vec2 mouse;

uniform sampler2D ourTexture;
uniform sampler2D ourTexture2;

void main()
{
	FragColor = texture(ourTexture2,vec2(texCoord.x, texCoord.y));
	if (FragColor.x == 0 && FragColor.y == 0) {
		FragColor = vec4(0.1,0.5,0.5,1);
	}
} ` + "\x00"
