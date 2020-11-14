package glog

import (
	"context"
)

// ctxLogKey is used as key to store *Logger in context
type ctxLogKey struct{}

// WithContext set *Logger into context and returned with given ctx.
func WithContext(ctx context.Context, l *Logger) context.Context {
	// If nil logger was given, return ctx directly
	if l == nil {
		return ctx
	}
	// Do not store the same logger.
	if lp, ok := ctx.Value(ctxLogKey{}).(*Logger); ok && lp == l {
		return ctx
	}
	return context.WithValue(ctx, ctxLogKey{}, l)
}

// FromContext get *Logger from context.
// NOTICE: This must be called after WithContext, if not a nil pointer is returned.
func FromContext(ctx context.Context) *Logger {
	l, ok := ctx.Value(ctxLogKey{}).(*Logger)
	if !ok {
		return nil
	}
	return l
}

// FromContextDefault get *Logger from context.
// And it will return a default *Logger if no *Logger was set before.
func FromContextDefault(ctx context.Context) *Logger {
	l, ok := ctx.Value(ctxLogKey{}).(*Logger)
	if !ok {
		return NewDefault()
	}
	return l
}
