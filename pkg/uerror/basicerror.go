package uerror

import (
	"encoding/binary"
	"errors"
)

type BasicError struct {
	errorCategory ErrorCategory
	code          int32
	message       string
}

func (b *BasicError) Code() int32 {
	return b.code
}
func (b *BasicError) ErrorCategory() ErrorCategory {
	return b.errorCategory
}

func (b *BasicError) Error() string {
	return b.message
}

var _ error = (*BasicError)(nil)

// Size 返回序列化后的字节大小
func (b *BasicError) Size() int {
	return 1 + 4 + 4 + len(b.message) // 1 byte for level, 4 bytes for code, 4 bytes for message length, and message itself
}

// MarshalToSizedBuffer 将 BasicError 序列化到提供的字节切片中
func (b *BasicError) MarshalToSizedBuffer(buf []byte) (int, error) {
	if len(buf) < b.Size() {
		return 0, errors.New("buffer is too small")
	}

	i := len(buf)

	// Write message
	i -= len(b.message)
	copy(buf[i:], b.message)

	// Write message length
	i -= 4
	binary.BigEndian.PutUint32(buf[i:], uint32(len(b.message)))

	// Write code
	i -= 4
	binary.BigEndian.PutUint32(buf[i:], uint32(b.code))

	// Write level
	i--
	buf[i] = byte(b.errorCategory)

	return len(buf) - i, nil
}

// Unmarshal 从字节切片中反序列化 BasicError
func (b *BasicError) Unmarshal(data []byte) error {
	if len(data) < 9 { // At least 1 byte for level, 4 for code, and 4 for message length
		return errors.New("buffer is too small")
	}

	b.errorCategory = ErrorCategory(data[0])
	b.code = int32(binary.BigEndian.Uint32(data[1:5]))
	messageLen := binary.BigEndian.Uint32(data[5:9])

	if len(data) < 9+int(messageLen) {
		return errors.New("data is too short for message")
	}

	b.message = string(data[9 : 9+messageLen])

	return nil
}
