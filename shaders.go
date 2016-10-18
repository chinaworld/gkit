package gkit

import (
	"bytes"
	"fmt"

	"github.com/go-gl/gl/v3.2-core/gl"
)

const (
	vertexShaderSource = `
`

	fragmentShaderSource = `
`
)

func loadShader(shaderType uint32, source string) (uint32, error) {
	content := gl.Str(source + "\x00")

	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, 1, &content, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := bytes.Repeat([]byte{0}, int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("Failed to compile shader: %s", log)
	}
}
