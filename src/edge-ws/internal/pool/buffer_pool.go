package pool

import "sync"

var BufferPool = sync.Pool{
	New: func() any {
		b := make([]byte, 4096)
		return &b
	},
}

func GetBuffer() *[]byte {
	buffer, ok := BufferPool.Get().(*[]byte)
	if !ok {
		return nil
	}
	return buffer
}

func PutBuffer(b *[]byte) {
	BufferPool.Put(b)
}
