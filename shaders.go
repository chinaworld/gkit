package gkit

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
)

const (
	vertexShaderSource = `
#version 150

in vec4 vPosition;

void main() {
  gl_Position = vPosition;
}
`

	fragmentShaderSource = `
#version 150

out vec4 fColor;

void main() {
  fColor = vec4(0.0, 0.0, 1.0, 1.0);
}
`
)

func loadShader(shaderType uint32, source string) (uint32, error) {
	ok := false
	content, free := gl.Strs(source + "\x00")
	defer free()

	shader := gl.CreateShader(shaderType)
	defer func() {
		if !ok {
			gl.DeleteShader(shader)
		}
	}()
	gl.ShaderSource(shader, 1, content, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("Failed to compile shader: %s", log)
	}

	ok = true
	return shader, nil
}

func loadShaders() (uint32, error) {
	ok := false
	program := gl.CreateProgram()
	defer func() {
		if !ok {
			gl.DeleteProgram(program)
		}
	}()

	vertexShader, err := loadShader(gl.VERTEX_SHADER, vertexShaderSource)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := loadShader(gl.FRAGMENT_SHADER, fragmentShaderSource)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fragmentShader)

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

		return 0, fmt.Errorf("Failed to link program: %s", log)
	}

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, gl.Ptr(nil))

	ok = true
	return program, nil
}
