package gkit

import (
	"github.com/go-gl/gl/v3.2-core/gl"
)

type View struct {
	vao uint32
	vbo uint32
}

func (v *View) setupGl() bool {
	ok := false
	gl.GenVertexArrays(1, &v.vao)
	if v.vao == 0 {
		return false
	}
	defer func() {
		if !ok {
			gl.DeleteVertexArrays(1, &v.vao)
			v.vao = 0
		}
	}()
	gl.BindVertexArray(v.vao)
	defer gl.BindVertexArray(0)

	gl.GenBuffers(1, &v.vbo)
	if v.vbo == 0 {
		return false
	}
	defer func() {
		if !ok {
			gl.DeleteBuffers(1, &v.vbo)
			v.vbo = 0
		}
	}()
	gl.BindBuffer(gl.ARRAY_BUFFER, v.vbo)
	defer gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	vertices := []float32{
		-1.0, -1.0,
		1.0, -1.0,
		-1.0, 1.0,
		1.0, -1.0,
		1.0, 1.0,
		-1.0, 1.0,
	}
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	ok = true
	return ok
}

func (v *View) Draw() {
	if v.vao == 0 && !v.setupGl() {
		return
	}

	gl.BindVertexArray(v.vao)
	defer gl.BindVertexArray(0)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func (v *View) Destroy() {
	if v.vao != 0 {
		gl.BindVertexArray(v.vao)
		if v.vbo != 0 {
			gl.BindBuffer(gl.ARRAY_BUFFER, 0)
			gl.DeleteBuffers(1, &v.vbo)
			v.vbo = 0
		}
		gl.BindVertexArray(0)
		gl.DeleteVertexArrays(1, &v.vao)
		v.vao = 0
	}
}
