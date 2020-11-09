package glog

import (
	"context"
)

// ctxLogKey is used as key to store logger in context
const (
	ctxLogKey = "glog"
)

// WithContext set *Logger into given context and return
func WithContext(ctx context.Context, l *Logger) context.Context {
	// if nil logger was given, return ctx directly
	if l == nil {
		return ctx
	}
	if lp, ok := ctx.Value(ctxLogKey).(*Logger); ok && lp == l {
		// Do not store same logger.
		return ctx
	}
	return context.WithValue(ctx, ctxLogKey, l)
}

// FromContext get *Logger from context
// Notice: It will return a default logger if no Logger was set before
func FromContext(ctx context.Context) *Logger {
	l, ok := ctx.Value(ctxLogKey).(*Logger)
	if !ok {
		return NewDefault()
	}
	return l
}
