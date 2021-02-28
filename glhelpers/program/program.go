package program

import (
	"fmt"
	"strings"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	Id uint32
}

func (s *Shader) Use() {
	gl.UseProgram(s.Id)
}

func (s *Shader) SetUniform2f(uniformeName string, x, y float32) {
	location := gl.GetUniformLocation(s.Id, gl.Str(uniformeName+"\x00"))
	gl.Uniform2f(location, x, y)
}

func (s *Shader) SetUniform1i(uniformeName string, x int32) {
	location := gl.GetUniformLocation(s.Id, gl.Str(uniformeName+"\x00"))
	gl.Uniform1i(location, x)
}

func (s *Shader) SetUniform4fv(uniformeName string, vs mgl32.Mat4) {
	location := gl.GetUniformLocation(s.Id, gl.Str(uniformeName+"\x00"))
	gl.UniformMatrix4fv(location, 1, false, &vs[0])
}

func (s *Shader) Delete() {
	gl.DeleteProgram(s.Id)
}

func New(vertexShaderSource, fragmentShaderSource string) (Shader, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return Shader{}, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return Shader{}, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return Shader{}, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return Shader{Id: program}, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	source += "\x00"
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
