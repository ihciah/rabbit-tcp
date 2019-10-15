package connection

type ByteRingBuffer struct {
	buffer []byte
	head   int
	tail   int
}

func NewByteRingBuffer(size uint32) ByteRingBuffer {
	buffer := make([]byte, size)
	return ByteRingBuffer{
		buffer: buffer,
	}
}

func (rb *ByteRingBuffer) OverWrite(data []byte) {
	if len(rb.buffer) < len(data) {
		rb.buffer = make([]byte, len(data))
	}
	n := copy(rb.buffer, data)
	rb.head = 0
	rb.tail = n
}

func (rb *ByteRingBuffer) Read(data []byte) int {
	n := len(data)
	if n > rb.tail-rb.head {
		n = rb.tail - rb.head
	}
	copy(data, rb.buffer[rb.head:rb.tail])
	rb.head += n
	return n
}

func (rb *ByteRingBuffer) Empty() bool {
	return rb.tail-rb.head == 0
}
