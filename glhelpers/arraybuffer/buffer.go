package arraybuffer

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type ArrayBuffer struct {
	id uint32
}

func (b *ArrayBuffer) Delete() {
	gl.DeleteBuffers(1, &b.id)
}

func (b *ArrayBuffer) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, b.id)
}

func New(triangleVertices []float32) ArrayBuffer {
	const float32Size = 4

	var vbo uint32

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(triangleVertices)*float32Size, gl.Ptr(triangleVertices), gl.STATIC_DRAW)

	return ArrayBuffer{
		id: vbo,
	}
}
