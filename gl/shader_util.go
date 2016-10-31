package gl

import (
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
)

type compileError string

func (e compileError) Error() string {
	return string(e)
}

var (
	glPainterGenVAOError        error = compileError("Failed to create Vertex Array.")
	glPainterGenVBOError              = compileError("Failed to create Vertex Buffers.")
	createShaderError                 = compileError("Failed to create Shader.")
	glPainterCreateProgramError       = compileError("Failed to create Program.")
)

func glPainterGetUniformLocationError(name string) error {
	return compileError("Failed to get uniform location: " + name)
}

func loadShader(shaderType uint32, source string) (uint32, error) {
	ok := false
	content, free := gl.Strs(source + "\x00")
	defer free()

	shader := gl.CreateShader(shaderType)
	if shader == 0 {
		return 0, compileError("Could not create shader.")
	}
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

		return 0, compileError(log)
	}

	ok = true
	return shader, nil
}

func linkProgram(glProgram uint32) error {
	gl.LinkProgram(glProgram)

	var status int32
	gl.GetProgramiv(glProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(glProgram, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(glProgram, logLength, nil, gl.Str(log))

		return compileError(log)
	}

	return nil
}
