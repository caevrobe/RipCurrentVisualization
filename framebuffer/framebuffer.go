package framebuffer

type FrameBuffer struct {
	width      int
	height     int
	unitSize   int
	components int
	Buffer     []byte
	cursor     int
}

func New(width, height, unitSize, components int) *FrameBuffer {
	return &FrameBuffer{
		width,
		height,
		unitSize,
		components,
		make([]byte, width*height*unitSize*components),
		0,
	}
}

func (b *FrameBuffer) Reset() {
	b.cursor = 0
}

// todo ensure cursor is in range
func (b *FrameBuffer) Next() interface{} {
	defer func() { b.cursor += b.unitSize }()

	slice := b.Buffer[b.cursor : b.cursor+b.unitSize]
	value := int32(0)

	for i := 0; i < b.unitSize; i++ {
		value += int32(slice[i])

		if i != b.unitSize-1 {
			value <<= 8
		}
	}

	return value
}
