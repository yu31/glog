package glog

import (
	"bytes"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLogger_NewDefault(t *testing.T) {
	l := NewDefault()
	require.Equal(t, l.level, DebugLevel)
	require.Equal(t, l.timeLayout, defaultTimeLayout)
	require.False(t, l.caller)
	require.True(t, reflect.DeepEqual(l.executor, DefaultExecutor))
	require.True(t, reflect.DeepEqual(l.encoderFunc(), TextEncoder()))
	require.True(t, reflect.DeepEqual(l.fields, TextEncoder()))
	require.True(t, reflect.DeepEqual(l.errorOutput, os.Stderr))
}

func TestLogger_Output(t *testing.T) {
	l := NewDefault().WithCaller(true)

	l.Info().Msg("HelloWorld").Fire()

	l.WithFields().AddString("extra-k1", "extra-v1")
	l.WithFields().AddString("extra_k2", "extra-v2")

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
		Duration("Duration", time.Duration(1)).
		Time("Time", time.Now(), defaultTimeLayout).
		Any("Interface", []string{"i1", "i2", "i3"}).
		Fire()

	// all level
	l.Debug().Msg("the debug log message").Fire()
	l.Info().Msg("the info log message").Fire()
	l.Warn().Msg("the warn log message").Fire()
	l.Error().Msg("the error log message").Fire()
	l.Fatal().Msg("the fatal log message").Fire()
}

func TestLogger_WithLevel(t *testing.T) {
	var b bytes.Buffer
	l := NewDefault()
	l.WithExecutor(MatchExecutor(&b, nil))
	l.WithLevel(InfoLevel)

	l.Debug().Msg("HelloWorld").Fire()

	require.Equal(t, len(b.Bytes()), 0)
}

func TestLogger_WithCaller(t *testing.T) {
	var b bytes.Buffer

	l := NewDefault().WithExecutor(MatchExecutor(&b, nil)).WithCaller(true)
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

	l := NewDefault()
	l.WithExecutor(MatchExecutor(&b, nil))

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

func TestLogger_Clone(t *testing.T) {
	var b bytes.Buffer

	l := NewDefault()
	l.WithExecutor(MatchExecutor(&b, nil))
	l.WithFields().AddString("filed-k1", "filed-v1")

	nl := l.Clone()

	require.Equal(t, l.level, nl.level)
	require.Equal(t, l.timeLayout, nl.timeLayout)
	require.Equal(t, l.caller, nl.caller)
	require.Equal(t, reflect.ValueOf(l.executor).Pointer(), reflect.ValueOf(nl.executor).Pointer())
	require.Equal(t, reflect.ValueOf(l.encoderFunc).Pointer(), reflect.ValueOf(nl.encoderFunc).Pointer())
	require.Equal(t, reflect.ValueOf(l.errorOutput).Pointer(), reflect.ValueOf(nl.errorOutput).Pointer())
	require.NotEqual(t, reflect.ValueOf(l.fields).Pointer(), reflect.ValueOf(nl.fields).Pointer())
	require.Equal(t, l.fields.Bytes(), nl.fields.Bytes())

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
}
