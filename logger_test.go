package glog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoggerNewDefault(t *testing.T) {
	l := NewDefault()
	require.Equal(t, l.ctx, context.Background())
	require.Equal(t, l.level, DebugLevel)
	require.Equal(t, l.timeLayout, defaultTimeLayout)
	require.False(t, l.caller)
	require.True(t, reflect.DeepEqual(l.exporter, DefaultExporter))
	require.True(t, reflect.DeepEqual(l.encoderFunc(), TextEncoder()))
	require.True(t, reflect.DeepEqual(l.fields, TextEncoder()))
	require.True(t, reflect.DeepEqual(l.errorOutput, os.Stderr))
}

type loggerWithContext struct {
	record *Record
}

func (exp *loggerWithContext) Export(record *Record) (err error) {
	exp.record = record
	return
}

func (exp *loggerWithContext) Close() (err error) {
	return
}

func TestLogger_WithContext(t *testing.T) {
	type ctxKey struct{}
	ctx := context.WithValue(context.Background(), ctxKey{}, "v1")
	exporter := &loggerWithContext{}

	l := NewDefault().WithContext(ctx).WithExporter(exporter)
	require.Equal(t, l.ctx, ctx)
	require.True(t, reflect.DeepEqual(l.ctx, ctx))

	l.Error().Msg("Hello World").Fire()
	require.Equal(t, exporter.record.Level(), ErrorLevel)
	require.NotNil(t, exporter.record.Context())
	require.Equal(t, exporter.record.Context(), ctx)
	require.True(t, reflect.DeepEqual(exporter.record.Context(), ctx))
}

func TestLogger_WithLevel(t *testing.T) {
	var b bytes.Buffer
	l := NewDefault()
	l.WithExporter(StandardExporter(&b))
	l.WithLevel(InfoLevel)

	l.Debug().Msg("HelloWorld").Fire()

	require.Equal(t, len(b.Bytes()), 0)
}

func TestLogger_WithCaller(t *testing.T) {
	var b bytes.Buffer
	l := NewDefault().WithExporter(StandardExporter(&b)).WithCaller(true)
	l.WithFields().AddString("k1", "v1")
	l.WithFields().AddString("k2", "v2")

	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)

	l.Info().Msg("test caller 1").Fire()
	l.Info().Msg("test caller 2").Fire()

	require.Equal(t, strings.Count(b.String(), file), 2)
}

func TestLogger_WithFields(t *testing.T) {
	var b bytes.Buffer
	l := NewDefault().WithExporter(StandardExporter(&b))

	l.WithFields().AddString("req-k1", "req-v1")
	l.WithFields().AddString("dup-key", "dup-v1")
	l.WithFields().AddString("dup-key", "dup-v2")

	l.Info().Msg("Hello World").Fire()

	s := b.String()

	require.Contains(t, s, "req-k1")
	require.Contains(t, s, "req-v1")
	require.Contains(t, s, "dup-key")
	require.Contains(t, s, "dup-v1")
	require.Contains(t, s, "dup-v2")
	require.Equal(t, strings.Count(s, "dup-key"), 2)
}

func TestLoggerWithTextEncoder(t *testing.T) {
	var eb bytes.Buffer
	l := NewDefault().WithCaller(true).WithErrorOutput(&eb)

	l.Info().Msg("HelloWorld").Fire()

	l.WithFields().AddString("extra-k1", "extra-v1")
	l.WithFields().AddString("extra_k2", "extra-v2")

	l.Info().
		Msg("Test logger out").
		String("String", "Value").
		Strings("Strings", []string{"a", "b", "c"}).
		Byte("Byte", 'a').
		Bytes("Bytes", []byte("abc")).
		Int64("Int64", 64).
		Int64s("Int64s", []int64{123, 456, 789}).
		Complex128("Complex128", complex(1, 2)).
		Float64("Float64", 99.99).
		Strings("Strings", []string{"a", "b", "c"}).
		Bytes("Bytes", []byte("HEllO")).
		Bool("Bool", true).
		Time("Time", time.Now(), defaultTimeLayout).
		Any("Interface1", []string{"i1", "i2", "i3"}).
		Any("Interface2", nil).
		Fire()

	// all level
	l.Debug().Msg("the debug log message").Fire()
	l.Info().Msg("the info log message").Fire()
	l.Warn().Msg("the warn log message").Fire()
	l.Error().Msg("the error log message").Fire()
	l.Fatal().Msg("the fatal log message").Fire()

	require.Equal(t, eb.Len(), 0)
}

func TestLoggerWithJSONEncoder(t *testing.T) {
	var eb bytes.Buffer
	var b bytes.Buffer
	l := NewDefault().WithEncoderFunc(JSONEncoder).WithExporter(StandardExporter(&b)).
		WithCaller(true).WithErrorOutput(&eb)

	l.Info().
		Msg("test logger out").
		String("String", "Value").
		Strings("Strings", []string{"a", "b", "c"}).
		Byte("Byte", 'a').
		Bytes("Bytes", []byte("abc")).
		Int64("Int64", 64).
		Int64s("Int64s", []int64{641, 642, 643}).
		Complex128("Complex128", complex(1, 2)).
		Float64("Float64", 99.99).
		Strings("Strings", []string{"a", "b", "c"}).
		Bytes("Bytes", []byte("HEllO")).
		Bool("Bool", true).
		Time("Time", time.Now(), time.RFC3339Nano).
		Any("Interface1", []string{"i1", "i2", "i3"}).
		Any("Interface2", nil).
		Fire()

	require.Greater(t, b.Len(), 0)
	c := b.Bytes()[:b.Len()-1]

	fmt.Println(string(c))

	m := make(map[string]interface{})
	err := json.Unmarshal(c, &m)
	require.Nil(t, err, "%q", string(c))

	require.Equal(t, eb.Len(), 0)
}

func TestLogger_Clone(t *testing.T) {
	type ctxKey struct{}
	ctx1 := context.WithValue(context.Background(), ctxKey{}, "v1")

	var eb bytes.Buffer
	var b bytes.Buffer
	l := NewDefault().
		WithExporter(StandardExporter(&b)).
		WithErrorOutput(&eb).
		WithContext(ctx1)

	l.WithFields().AddString("filed-k1", "filed-v1")
	require.True(t, l.isRoot)

	nl := l.Clone()
	require.Equal(t, l.ctx, nl.ctx)
	require.Equal(t, l.level, nl.level)
	require.Equal(t, l.timeLayout, nl.timeLayout)
	require.Equal(t, l.caller, nl.caller)
	require.Equal(t, reflect.ValueOf(l.exporter).Pointer(), reflect.ValueOf(nl.exporter).Pointer())
	require.Equal(t, reflect.ValueOf(l.encoderFunc).Pointer(), reflect.ValueOf(nl.encoderFunc).Pointer())
	require.Equal(t, reflect.ValueOf(l.errorOutput).Pointer(), reflect.ValueOf(nl.errorOutput).Pointer())
	require.NotEqual(t, reflect.ValueOf(l.fields).Pointer(), reflect.ValueOf(nl.fields).Pointer())
	require.Equal(t, l.fields.Bytes(), nl.fields.Bytes())
	require.False(t, nl.isRoot)

	// Reset the context in new logger.
	ctx2 := context.WithValue(context.Background(), ctxKey{}, "v2")
	nl.WithContext(ctx2)
	require.Equal(t, l.ctx, ctx1)
	require.True(t, reflect.DeepEqual(l.ctx, ctx1))
	require.Equal(t, nl.ctx, ctx2)
	require.True(t, reflect.DeepEqual(nl.ctx, ctx2))
	require.NotEqual(t, l.ctx, nl.ctx)
	require.False(t, reflect.DeepEqual(l.ctx, nl.ctx))

	nl.WithFields().AddString("filed-k2", "filed-v2")
	b.Reset()

	nl.Info().Msg("HelloWorld").Fire()
	s := b.String()

	// The 'l1' will inheritance the l's fields
	require.Contains(t, s, "filed-k1")
	require.Contains(t, s, "filed-v1")
	require.Contains(t, s, "filed-k2")
	require.Contains(t, s, "filed-v2")

	// Any changed in 'l1' does not affects 'l'
	b.Reset()
	l.Info().Msg("HelloWorld").Fire()
	s = b.String()

	require.NotContains(t, s, "filed-k2")
	require.NotContains(t, s, "filed-v2")

	// Any changed in "l" does not affects 'l1' also
	_ = l.Close()
	b.Reset()

	nl.Info().Msg("HelloWorld").Fire()
	s = b.String()
	require.Contains(t, s, "HelloWorld")
	require.Contains(t, s, "filed-k1")
	require.Contains(t, s, "filed-v1")
	require.Contains(t, s, "filed-k2")
	require.Contains(t, s, "filed-v2")

	require.Equal(t, eb.Len(), 0)
}

type loggerWriterCloser struct {
	data     []byte
	isClosed bool
}

func (w *loggerWriterCloser) Write(p []byte) (n int, err error) {
	if w.isClosed {
		return 0, errors.New("write in a closed writer")
	}
	w.data = append(w.data, p...)
	return len(p), nil
}

func (w *loggerWriterCloser) Close() error {
	w.isClosed = true
	return nil
}

func (w *loggerWriterCloser) resetData() {
	w.data = nil
}

func TestLogger_Close(t *testing.T) {
	t.Run("CloseRootLogger1", func(t *testing.T) {
		writer := &loggerWriterCloser{}
		lp := NewDefault().WithExporter(StandardExporter(writer))
		require.False(t, writer.isClosed)
		lp.Debug().Msg("Hello World").Fire()
		require.True(t, len(writer.data) > 0)
		writer.resetData()
		_ = lp.Close()
		require.True(t, writer.isClosed)
		require.Panics(t, func() {
			lp.Debug().Msg("Hello World").Fire()
		})
	})

	t.Run("CloseCloneLogger1", func(t *testing.T) {
		writer := &loggerWriterCloser{}
		lp := NewDefault().WithExporter(StandardExporter(writer))

		// clone a new logger from root logger.
		lg1 := lp.Clone()
		lg1.Debug().Msg("Hello World").Fire()
		require.True(t, len(writer.data) > 0)
		writer.resetData()

		// clone a new logger from root logger.
		lg2 := lp.Clone()
		lg2.Debug().Msg("Hello World").Fire()
		require.True(t, len(writer.data) > 0)
		writer.resetData()

		// clone a new logger from lg1.
		lg3 := lg1.Clone()
		lg3.Debug().Msg("Hello World").Fire()
		require.True(t, len(writer.data) > 0)
		writer.resetData()

		// clone a new logger from lg2.
		lg4 := lg2.Clone()
		lg4.Debug().Msg("Hello World").Fire()
		require.True(t, len(writer.data) > 0)
		writer.resetData()

		// case1: close lg4
		_ = lg4.Close()
		require.False(t, writer.isClosed)
		require.NotPanics(t, func() {
			lp.Debug().Msg("Hello World").Fire()
			lg1.Debug().Msg("Hello World").Fire()
			lg2.Debug().Msg("Hello World").Fire()
			lg3.Debug().Msg("Hello World").Fire()
		})

		// case2: close lg2
		_ = lg2.Close()
		require.False(t, writer.isClosed)
		require.NotPanics(t, func() {
			lp.Debug().Msg("Hello World").Fire()
			lg1.Debug().Msg("Hello World").Fire()
			lg3.Debug().Msg("Hello World").Fire()
		})

		// case3: close lg1
		_ = lg1.Close()
		require.False(t, writer.isClosed)
		require.NotPanics(t, func() {
			lp.Debug().Msg("Hello World").Fire()
			lg3.Debug().Msg("Hello World").Fire()
		})

		// case4: close lg3
		_ = lg3.Close()
		require.False(t, writer.isClosed)
		require.NotPanics(t, func() {
			lp.Debug().Msg("Hello World").Fire()
		})

		// case5: close root logger.
		_ = lp.Close()
		require.True(t, writer.isClosed)

	})
}
