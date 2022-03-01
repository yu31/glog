package glog

import (
	"context"
)

// Record represents the Entry's content.
type Record struct {
	ctx   context.Context
	level Level
	data  []byte
}

// Context returns context where in Logger.
func (r *Record) Context() context.Context {
	return r.ctx
}

// Level returns the log level of the entry.
func (r *Record) Level() Level {
	return r.level
}

// Bytes returns the Entry's content.
func (r *Record) Bytes() []byte {
	return r.data
}

// Copy returns an copy of Entry's content.
func (r *Record) Copy() []byte {
	bs := make([]byte, len(r.data))
	copy(bs, r.data)
	return bs
}
