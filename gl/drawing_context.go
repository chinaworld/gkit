package gl

import (
	"image"
	"image/color"

	"github.com/go-gl/gl/v3.2-core/gl"

	"github.com/alex-ac/gkit"
)

const (
	glPainterVertexShaderSource = `
#version 330 core
layout(location = 0) in vec4 vPosition;
layout(location = 1) in vec4 vColorIn;
layout(location = 2) in vec2 vUVIn;
out vec4 vColor;
out vec2 vUV;

uniform uvec4 viewportSize;

void main() {
  vec4 position = vPosition * 2.0f / viewportSize - 1.0f;
  gl_Position = position * mat4(
      1.0, 0.0, 0.0, 0.0,
      0.0, -1.0, 0.0, 0.0,
      0.0, 0.0, -1.0, 0.0,
      0.0, 0.0, 0.0, 1.0);
  vColor = vColorIn;
  vUV = vUVIn;
}
`
	glPainterFragmentShaderSource = `
#version 330 core

in vec4 vColor;
in vec2 vUV;

out vec4 fColor;

uniform uvec2 maskSize;
uniform sampler2D mask;

void main() {
  mat3 translate = mat3(
    1.0, 0.0, 0.0,
    0.0, -1.0, 1.0,
    0.0, 0.0, 1.0);
  vec2 uv = (vec3(vUV / maskSize, 1.0) * translate).xy;
  if (uv.x > 0 && uv.x < 1 && uv.y > 0 && uv.y < 1) {
    // discard;
    fColor = texture(mask, uv);
  } else {
    vec4 color = vColor;
    fColor = color;
  }
}
`
)

const (
	attrCoords uint32 = iota
	attrColors
	attrUv

	attrFloatSize    = 4
	attrCoordsCount  = 3
	attrCoordsOffset = 0
	attrCoordsSize   = attrCoordsCount * attrFloatSize
	attrColorsCount  = 4
	attrColorOffset  = attrCoordsOffset + attrCoordsSize
	attrColorSize    = attrColorsCount * attrFloatSize
	attrUvCount      = 2
	attrUvOffset     = attrColorOffset + attrColorSize
	attrUvSize       = attrUvCount * attrFloatSize
	attrStride       = attrUvOffset + attrUvSize
)

type drawingContext struct {
	program uint32

	vao uint32
	vbo uint32

	viewportSizeLocation int32
	maskLocation         int32
	maskSizeLocation     int32

	texture uint32
	sampler uint32
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

	maskLocation := gl.GetUniformLocation(program, gl.Str("mask\x00"))
	if maskLocation < 0 {
		return nil, glPainterGetUniformLocationError("mask")
	}

	maskSizeLocation := gl.GetUniformLocation(program, gl.Str("maskSize\x00"))
	if maskSizeLocation < 0 {
		return nil, glPainterGetUniformLocationError("maskSize")
	}

	var texture uint32
	gl.GenTextures(1, &texture)
	if texture == 0 {
		return nil, glPainterGenTextureError
	}
	defer func() {
		if !ok {
			gl.DeleteTextures(1, &texture)
		}
	}()

	var sampler uint32
	gl.GenSamplers(1, &sampler)
	if sampler == 0 {
		return nil, glPainterGenSamplerError
	}
	defer func() {
		if !ok {
			gl.DeleteSamplers(1, &sampler)
		}
	}()

	ok = true
	return &drawingContext{
		program:              program,
		vao:                  vao,
		vbo:                  vbo,
		viewportSizeLocation: viewportSizeLocation,
		maskLocation:         maskLocation,
		maskSizeLocation:     maskSizeLocation,
		texture:              texture,
		sampler:              sampler,
	}, nil
}

func (g *drawingContext) Destroy() {
	gl.BindVertexArray(g.vao)
	gl.DeleteSamplers(1, &g.sampler)
	gl.DeleteTextures(1, &g.texture)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &g.vbo)
	gl.BindVertexArray(0)
	gl.DeleteVertexArrays(1, &g.vao)
	gl.UseProgram(0)
	gl.DeleteProgram(g.program)
}

func nearestPowerOf2(x uint32) uint32 {
	if x == 0 {
		return 1
	}
	return nearestPowerOf2(x>>1) + 1
}

func textureSideSize(w, h uint32) uint32 {
	if w > h {
		return 1 << nearestPowerOf2(w)
	}
	return 1 << nearestPowerOf2(h)
}

func (g *drawingContext) BeginPaint(width, height uint32) gkit.Painter {
	maskSideSize := textureSideSize(width, height)
	mask := image.NewGray(image.Rectangle{
		Max: image.Point{int(maskSideSize), int(maskSideSize)},
	})
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			mask.SetGray(x, y, color.Gray{0xff})
		}
	}
	return &painter{
		context:  g,
		mask:     mask,
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
	gl.VertexAttribPointer(
		attrCoords, attrCoordsCount, gl.FLOAT, false, attrStride, gl.PtrOffset(attrCoordsOffset))
	gl.VertexAttribPointer(
		attrColors, attrColorsCount, gl.FLOAT, false, attrStride, gl.PtrOffset(attrColorOffset))
	gl.VertexAttribPointer(
		attrUv, attrUvCount, gl.FLOAT, false, attrStride, gl.PtrOffset(attrUvOffset))
	gl.EnableVertexAttribArray(attrCoords)
	gl.EnableVertexAttribArray(attrColors)
	gl.EnableVertexAttribArray(attrUv)

	gl.Uniform1i(g.maskLocation, 0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, g.texture)
	// defer gl.BindTexture(gl.TEXTURE_2D, 0)

	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RED, int32(p.mask.Rect.Max.X), int32(p.mask.Rect.Max.Y), 0,
		gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(p.mask.Pix))

	gl.BindSampler(0, g.sampler)
	//defer gl.BindSampler(0, 0)

	gl.Uniform2ui(g.maskSizeLocation, uint32(p.mask.Rect.Max.X), uint32(p.mask.Rect.Max.Y))

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(p.vertices)/8))
}
