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
layout(location = 2) in vec2 vUVIn;
layout(location = 3) in vec3 vImageUVIn;
out vec4 vColor;
out vec2 vUV;
out vec3 vImageUV;

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
  vImageUV = vImageUVIn;
}
`
	glPainterFragmentShaderSource = `
#version 330 core

in vec4 vColor;
in vec2 vUV;
in vec3 vImageUV;

out vec4 fColor;

uniform uvec2 maskSize;
uniform sampler2D mask;
uniform sampler2DArray images;

void main() {
  vec2 uv = vUV / maskSize;
  fColor = vColor;
  if (vImageUV.z >= 0) {
    fColor = texture(images, vec3(vImageUV));
  }
  fColor.a *= texture(mask, uv).r;
  if (fColor.a == 0) {
    discard;
  }
}
`
)

const (
	attrCoords uint32 = iota
	attrColors
	attrUv
	attrImageUv

	attrFloatSize     = 4
	attrCoordsCount   = 3
	attrCoordsOffset  = 0
	attrCoordsSize    = attrCoordsCount * attrFloatSize
	attrColorsCount   = 4
	attrColorOffset   = attrCoordsOffset + attrCoordsSize
	attrColorSize     = attrColorsCount * attrFloatSize
	attrUvCount       = 2
	attrUvOffset      = attrColorOffset + attrColorSize
	attrUvSize        = attrUvCount * attrFloatSize
	attrImageUvCount  = 3
	attrImageUvOffset = attrUvOffset + attrUvSize
	attrImageUvSize   = attrImageUvCount * attrFloatSize

	attrStride = attrImageUvOffset + attrImageUvSize
)

type drawingContext struct {
	program uint32

	vao uint32
	vbo uint32

	viewportSizeLocation int32
	maskLocation         int32
	maskSizeLocation     int32
	imagesLocation       int32

	texture      uint32
	textureArray uint32
	sampler      uint32
	samplerArray uint32

	scaleFactor float32
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

	imagesLocation := gl.GetUniformLocation(program, gl.Str("images\x00"))
	if imagesLocation < 0 {
		return nil, glPainterGetUniformLocationError("images")
	}

	var textures [2]uint32
	gl.GenTextures(int32(len(textures)), &textures[0])
	for _, texture := range textures {
		if texture == 0 {
			return nil, glPainterGenTextureError
		}
	}
	defer func() {
		if !ok {
			gl.DeleteTextures(int32(len(textures)), &textures[0])
		}
	}()

	var samplers [2]uint32
	gl.GenSamplers(int32(len(samplers)), &samplers[0])
	for _, sampler := range samplers {
		if sampler == 0 {
			return nil, glPainterGenSamplerError
		}
	}
	defer func() {
		if !ok {
			gl.DeleteSamplers(int32(len(samplers)), &samplers[0])
		}
	}()
	sampler := samplers[0]
	gl.SamplerParameteri(sampler, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.SamplerParameteri(sampler, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	arraySampler := samplers[1]
	gl.SamplerParameteri(arraySampler, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.SamplerParameteri(arraySampler, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	ok = true
	return &drawingContext{
		program:              program,
		vao:                  vao,
		vbo:                  vbo,
		viewportSizeLocation: viewportSizeLocation,
		maskLocation:         maskLocation,
		maskSizeLocation:     maskSizeLocation,
		imagesLocation:       imagesLocation,
		texture:              textures[0],
		textureArray:         textures[1],
		sampler:              sampler,
		samplerArray:         arraySampler,
		scaleFactor:          1,
	}, nil
}

func (g *drawingContext) Destroy() {
	gl.BindVertexArray(g.vao)
	textures := []uint32{
		g.texture,
		g.textureArray,
	}
	samplers := []uint32{
		g.sampler,
		g.samplerArray,
	}
	gl.DeleteSamplers(int32(len(samplers)), &samplers[0])
	gl.DeleteTextures(int32(len(textures)), &textures[0])
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

func textureSideSize(s gkit.Size) uint32 {
	if s.Width > s.Height {
		return 1 << nearestPowerOf2(s.Width)
	}
	return 1 << nearestPowerOf2(s.Height)
}

func (g *drawingContext) BeginPaint(size gkit.Size) gkit.Painter {
	maskSideSize := textureSideSize(size)
	mask := image.NewGray(image.Rectangle{
		Max: image.Point{int(float32(maskSideSize) * g.scaleFactor), int(float32(maskSideSize) * g.scaleFactor)},
	})
	mask.Pix[0] = 0xff
	return &painter{
		context:     g,
		mask:        mask,
		images:      make([]*image.RGBA, 0),
		size:        size,
		vertices:    make([]float32, 0),
		scaleFactor: g.scaleFactor,
	}
}

func (g *drawingContext) EndPaint(gkitPainter gkit.Painter) {
	p := gkitPainter.(*painter)
	if len(p.vertices) == 0 {
		return
	}

	gl.UseProgram(g.program)
	defer gl.UseProgram(0)

	gl.Uniform4ui(g.viewportSizeLocation, p.size.Width, p.size.Height, 256, 1)

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
	gl.VertexAttribPointer(
		attrImageUv, attrImageUvCount, gl.FLOAT, false, attrStride, gl.PtrOffset(attrImageUvOffset))
	gl.EnableVertexAttribArray(attrCoords)
	gl.EnableVertexAttribArray(attrColors)
	gl.EnableVertexAttribArray(attrUv)
	gl.EnableVertexAttribArray(attrImageUv)

	gl.Uniform1i(g.maskLocation, 0)

	gl.ActiveTexture(gl.TEXTURE0)

	gl.BindTexture(gl.TEXTURE_2D, g.texture)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RED, int32(p.mask.Rect.Max.X), int32(p.mask.Rect.Max.Y), 0,
		gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(p.mask.Pix))

	gl.BindSampler(0, g.sampler)
	defer gl.BindSampler(0, 0)

	gl.Uniform1i(g.imagesLocation, 1)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, g.textureArray)
	var images []uint8
	for _, image := range p.images {
		if images == nil {
			images = image.Pix
		} else {
			images = append(images, image.Pix...)
		}
	}
	if len(images) > 0 {
		gl.TexImage3D(
			gl.TEXTURE_2D_ARRAY, 0, gl.RGBA,
			int32(p.mask.Rect.Max.X), int32(p.mask.Rect.Max.Y), int32(len(p.images)), 0,
			gl.RGBA, gl.UNSIGNED_INT_8_8_8_8_REV, gl.Ptr(images))
	}

	gl.BindSampler(1, g.samplerArray)
	defer gl.BindSampler(1, 0)

	gl.Enablei(gl.BLEND, 0)
	gl.BlendEquationSeparate(gl.FUNC_ADD, gl.FUNC_ADD)
	gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.ONE, gl.ZERO)

	gl.Uniform2ui(g.maskSizeLocation, uint32(float32(p.mask.Rect.Max.X)/p.scaleFactor), uint32(float32(p.mask.Rect.Max.Y)/p.scaleFactor))

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(p.vertices)/8))
}
