package util

import (
	"errors"
	"io"
)

type RingBuffer interface {
	io.Writer

	Bytes() []byte
	Reset()
}

type ringBuffer struct {
	data         []byte
	size         int64
	writeCursor  int64
	writtenCount int64
}

func NewRingBuffer(size int64) (RingBuffer, error) {
	if size <= 0 {
		return nil, errors.New("Size must be positive")
	}

	b := &ringBuffer{
		size: size,
		data: make([]byte, size),
	}

	return b, nil
}

func (b *ringBuffer) Write(buf []byte) (int, error) {

	n := len(buf)
	b.writtenCount += int64(n)

	if int64(n) > b.size {
		buf = buf[int64(n)-b.size:]
	}

	remain := b.size - b.writeCursor
	copy(b.data[b.writeCursor:], buf)
	if int64(len(buf)) > remain {
		copy(b.data, buf[remain:])
	}

	b.writeCursor = ((b.writeCursor + int64(len(buf))) % b.size)
	return n, nil
}

func (b *ringBuffer) Bytes() []byte {
	switch {
	case b.writtenCount >= b.size && b.writeCursor == 0:
		return b.data
	case b.writtenCount > b.size:
		out := make([]byte, b.size)
		copy(out, b.data[b.writeCursor:])
		copy(out[b.size-b.writeCursor:], b.data[:b.writeCursor])
		return out
	default:
		return b.data[:b.writeCursor]
	}
	return nil
}

func (b *ringBuffer) Reset() {
	b.writeCursor = 0
	b.writtenCount = 0
}
