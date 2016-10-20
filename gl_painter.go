package gkit

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

type glPainterVBOIndex int

const (
	glPainterNoVBO glPainterVBOIndex = iota - 1
	glPainterRectVBO

	glPainterVBOCount
)

const (
	glPainterVertexShaderSource = `
#version 140
in uvec4 vPosition;

uniform uvec4 viewportSize;
uniform uvec4 viewOrigin;
uniform uvec4 viewSize;

void main() {
  vec4 position = vec4(vPosition * viewSize + viewOrigin) * 2.0f / viewportSize - 1.0f;
  gl_Position = position * mat4(
      1.0, 0.0, 0.0, 0.0,
      0.0, -1.0, 0.0, 0.0,
      0.0, 0.0, 1.0, 0.0,
      0.0, 0.0, 0.0, 1.0);
}
`
	glPainterFragmentShaderSource = `
#version 140

uniform vec4 viewColor;
out vec4 fColor;

void main() {
  fColor = viewColor;
}
`
)

type glPainterInternal interface {
	setColor(c Color)
	drawRect(x, y, z, width, height uint32)
}

func normalizeCoords(x, y, width, height, maxWidth, maxHeight uint32) (uint32, uint32, uint32, uint32) {
	if x > maxWidth {
		x = maxWidth
	}
	if y > maxHeight {
		y = maxHeight
	}
	if x+width > maxWidth {
		width = maxWidth - x
	}
	if y+height > maxHeight {
		height = maxHeight - y
	}

	return x, y, width, height
}

type glPainter struct {
	program uint32

	vao uint32
	vbo [glPainterVBOCount]uint32

	currentVBO glPainterVBOIndex

	width, height uint32

	colorLocation        int32
	viewOriginLocation   int32
	viewSizeLocation     int32
	viewportSizeLocation int32
}

var _ glPainterInternal = &glPainter{}

func newGlPainter() (*glPainter, error) {
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

	var vbos [glPainterVBOCount]uint32
	gl.GenBuffers(int32(glPainterVBOCount), &vbos[0])
	for _, vbo := range vbos {
		if vbo == 0 {
			return nil, glPainterGenVBOError
		}
	}

	rectVertices := []uint32{
		0, 0,
		1, 0,
		1, 1,
		0, 0,
		1, 1,
		0, 1,
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, vbos[glPainterRectVBO])
	defer gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BufferData(gl.ARRAY_BUFFER, len(rectVertices)*4, gl.Ptr(rectVertices), gl.STATIC_DRAW)

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

	colorLocation := gl.GetUniformLocation(program, gl.Str("viewColor\x00"))
	if colorLocation < 0 {
		return nil, glPainterGetUniformLocationError("viewColor")
	}
	viewOriginLocation := gl.GetUniformLocation(program, gl.Str("viewOrigin\x00"))
	if viewOriginLocation < 0 {
		return nil, glPainterGetUniformLocationError("viewOrigin")
	}
	viewSizeLocation := gl.GetUniformLocation(program, gl.Str("viewSize\x00"))
	if viewSizeLocation < 0 {
		return nil, glPainterGetUniformLocationError("viewSize")
	}
	viewportSizeLocation := gl.GetUniformLocation(program, gl.Str("viewportSize\x00"))
	if viewportSizeLocation < 0 {
		return nil, glPainterGetUniformLocationError("viewportSize")
	}

	ok = true
	return &glPainter{
		program:              program,
		vao:                  vao,
		vbo:                  vbos,
		currentVBO:           glPainterNoVBO,
		colorLocation:        colorLocation,
		viewOriginLocation:   viewOriginLocation,
		viewSizeLocation:     viewSizeLocation,
		viewportSizeLocation: viewportSizeLocation,
	}, nil
}

func (g *glPainter) Destroy() {
	gl.BindVertexArray(g.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	g.currentVBO = glPainterNoVBO
	gl.DeleteBuffers(int32(glPainterVBOCount), &g.vbo[0])
	gl.BindVertexArray(0)
	gl.DeleteVertexArrays(1, &g.vao)
	gl.UseProgram(0)
	gl.DeleteProgram(g.program)
}

func (g *glPainter) BeginPaint() (Painter, func()) {
	gl.UseProgram(g.program)
	gl.BindVertexArray(g.vao)

	viewportSize := make([]int32, 4)
	gl.GetIntegerv(gl.VIEWPORT, &viewportSize[0])
	g.width, g.height = uint32(viewportSize[2]), uint32(viewportSize[3])
	gl.Uniform4ui(g.viewportSizeLocation, g.width, g.height, 256, 1)

	return &glPainterProxy{
			impl:   g,
			x:      0,
			y:      0,
			width:  g.width,
			height: g.height,
		}, func() {
			gl.BindVertexArray(0)
			gl.UseProgram(0)
		}
}

type glPainterProxy struct {
	impl glPainterInternal

	x, y, width, height uint32
}

var _ Painter = &glPainterProxy{}
var _ glPainterInternal = &glPainterProxy{}

func (g *glPainter) setColor(c Color) {
	colorVec4 := c.vec4()
	gl.Uniform4f(g.colorLocation,
		colorVec4[0], colorVec4[1], colorVec4[2], colorVec4[3])
}

func (g *glPainter) drawRect(x, y, z, width, height uint32) {
	if g.currentVBO != glPainterRectVBO {
		gl.BindBuffer(gl.ARRAY_BUFFER, g.vbo[glPainterRectVBO])
		g.currentVBO = glPainterRectVBO
		gl.VertexAttribPointer(0, 2, gl.UNSIGNED_INT, false, 0, gl.Ptr(nil))
		gl.EnableVertexAttribArray(0)
	}

	gl.Uniform4ui(g.viewOriginLocation, x, y, z, 0)
	gl.Uniform4ui(g.viewSizeLocation, width, height, 1, 1)

	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func (p *glPainterProxy) DrawRect(x, y, width, height uint32) {
	p.drawRect(x, y, 0, width, height)
}

func (p *glPainterProxy) drawRect(x, y, z, width, height uint32) {
	x, y, width, height = normalizeCoords(
		x, y, width, height, p.width, p.height)
	p.impl.drawRect(p.x+x, p.y+y, z+1, width, height)
}

func (p *glPainterProxy) SubPainter(x, y, width, height uint32) Painter {
	x, y, width, height = normalizeCoords(
		x, y, width, height, p.width, p.height)
	return &glPainterProxy{
		impl:   p,
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}

func (p *glPainterProxy) SetColor(c Color) {
	p.setColor(c)
}

func (p *glPainterProxy) setColor(c Color) {
	p.impl.setColor(c)
}
