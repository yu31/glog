package buffer

import "sync"

const (
	size = 256 // by default, create 256 bits buffers
	//maxSize = 1 << 16 // 64KiB
	maxSize = 1 << 14 // 16KiB
)

// A Pool is a type-safe wrapper around a sync.Pool.
type Pool struct {
	p *sync.Pool
}

// NewPool constructs a new Pool.
func NewPool() Pool {
	return Pool{p: &sync.Pool{
		New: func() interface{} {
			return &Buffer{bs: make([]byte, 0, size)}
		},
	}}
}

// Get retrieves a Buffer from the pool, creating one if necessary.
func (p Pool) Get() *Buffer {
	buf := p.p.Get().(*Buffer)
	buf.Reset()
	buf.pool = p
	return buf
}

func (p Pool) put(buf *Buffer) {
	// See https://golang.org/issue/23199
	if cap(buf.bs) > maxSize {
		return
	}
	p.p.Put(buf)
}
