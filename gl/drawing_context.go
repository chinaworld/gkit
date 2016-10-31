package gl

import (
	"image"

	"github.com/go-gl/gl/v3.2-core/gl"

	"github.com/alex-ac/gkit"
)

const (
	glPainterVertexShaderSource = `
#version 330 core
layout(location = 0) in vec4 vPosition;
layout(location = 1) in vec4 vColorIn;
out vec4 vColor;

uniform uvec4 viewportSize;

void main() {
  vec4 position = vPosition * 2.0f / viewportSize - 1.0f;
  gl_Position = position * mat4(
      1.0, 0.0, 0.0, 0.0,
      0.0, -1.0, 0.0, 0.0,
      0.0, 0.0, -1.0, 0.0,
      0.0, 0.0, 0.0, 1.0);
  vColor = vColorIn;
}
`
	glPainterFragmentShaderSource = `
#version 330 core

// uniform vec4 viewColor;

in vec4 vColor;
out vec4 fColor;

void main() {
  vec4 color = vColor;
  fColor = color;
}
`
)

type drawingContext struct {
	program uint32

	vao uint32
	vbo uint32

	viewportSizeLocation int32
}

func newDrawingContext() (*drawingContext, error) {
	ok := false
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	if vao == 0 {
		return nil, glPainterGenVAOError
	}
	defer func() {
		if !ok {
			gl.DeleteVertexArrays(1, &vao)
		}
	}()

	gl.BindVertexArray(vao)
	defer gl.BindVertexArray(0)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	if vbo == 0 {
		return nil, glPainterGenVBOError
	}

	program := gl.CreateProgram()
	if program == 0 {
		return nil, glPainterCreateProgramError
	}
	defer func() {
		if !ok {
			gl.DeleteProgram(program)
		}
	}()

	vertexShader, err := loadShader(gl.VERTEX_SHADER, glPainterVertexShaderSource)
	if err != nil {
		return nil, err
	}
	defer gl.DeleteShader(vertexShader)
	gl.AttachShader(program, vertexShader)

	fragmentShader, err := loadShader(gl.FRAGMENT_SHADER, glPainterFragmentShaderSource)
	if err != nil {
		return nil, err
	}
	defer gl.DeleteShader(fragmentShader)
	gl.AttachShader(program, fragmentShader)

	err = linkProgram(program)
	if err != nil {
		return nil, err
	}

	viewportSizeLocation := gl.GetUniformLocation(program, gl.Str("viewportSize\x00"))
	if viewportSizeLocation < 0 {
		return nil, glPainterGetUniformLocationError("viewportSize")
	}

	ok = true
	return &drawingContext{
		program:              program,
		vao:                  vao,
		vbo:                  vbo,
		viewportSizeLocation: viewportSizeLocation,
	}, nil
}

func (g *drawingContext) Destroy() {
	gl.BindVertexArray(g.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &g.vbo)
	gl.BindVertexArray(0)
	gl.DeleteVertexArrays(1, &g.vao)
	gl.UseProgram(0)
	gl.DeleteProgram(g.program)
}

func (g *drawingContext) BeginPaint(width, height uint32) gkit.Painter {
	return &painter{
		context:  g,
		mask:     image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{1, 1}}),
		width:    width,
		height:   height,
		vertices: make([]float32, 0, 6),
	}
}

func (g *drawingContext) EndPaint(gkitPainter gkit.Painter) {
	p := gkitPainter.(*painter)
	if len(p.vertices) == 0 {
		return
	}

	gl.UseProgram(g.program)
	defer gl.UseProgram(0)

	gl.Uniform4ui(g.viewportSizeLocation, p.width, p.height, 256, 1)

	gl.BindVertexArray(g.vao)
	defer gl.BindVertexArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, g.vbo)
	defer gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BufferData(gl.ARRAY_BUFFER, len(p.vertices)*4, gl.Ptr(p.vertices), gl.STREAM_DRAW)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 32, gl.Ptr(nil))
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 32, gl.PtrOffset(16))
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(p.vertices)/8))
}
