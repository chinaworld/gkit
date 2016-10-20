package gkit

type Color uint32

func RGBA(r, g, b, a uint8) Color {
	return Color(uint32(r)<<24 | uint32(g)<<16 | uint32(b)<<8 | uint32(a))
}

func (c Color) R() uint8 {
	return uint8(c >> 24)
}

func (c Color) G() uint8 {
	return uint8(c >> 16)
}

func (c Color) B() uint8 {
	return uint8(c >> 8)
}

func (c Color) A() uint8 {
	return uint8(c)
}

func (c Color) vec4() [4]float32 {
	return [4]float32{
		float32(c.R()) / 256.0,
		float32(c.G()) / 256.0,
		float32(c.B()) / 256.0,
		float32(c.A()) / 256.0,
	}
}
