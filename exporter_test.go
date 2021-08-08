package glog

import (
	"bytes"
	"context"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

type CustomExporter struct {
	record *Record
}

func (exp *CustomExporter) Export(record *Record) (err error) {
	exp.record = record
	return
}

func (exp *CustomExporter) Close() (err error) {
	return
}

func TestExporter_Custom(t *testing.T) {
	type ctxKey struct{}
	ctx := context.WithValue(context.Background(), ctxKey{}, "v1")

	export := &CustomExporter{}

	l := NewDefault().WithExporter(export).WithContext(ctx)

	msg1 := "Hello One World"
	l.Debug().Msg(msg1).Fire()
	require.NotNil(t, export.record)
	require.Equal(t, ctx, export.record.Context())
	require.Equal(t, DebugLevel, export.record.Level())
	require.Contains(t, string(export.record.Bytes()), msg1)
	require.Contains(t, string(export.record.Copy()), msg1)

	msg2 := "Hello Two World"
	l.Info().Msg(msg2).Fire()
	require.NotNil(t, export.record)
	require.Equal(t, InfoLevel, export.record.Level())
	require.Contains(t, string(export.record.Bytes()), msg2)
	require.Contains(t, string(export.record.Copy()), msg2)

	// Test copy
	b1 := export.record.Bytes()
	b2 := export.record.Bytes()
	b3 := export.record.Copy()
	p1 := (*reflect.SliceHeader)(unsafe.Pointer(&b1))
	p2 := (*reflect.SliceHeader)(unsafe.Pointer(&b2))
	p3 := (*reflect.SliceHeader)(unsafe.Pointer(&b3))

	require.Equal(t, p1, p2)
	require.NotEqual(t, p1, p3)
}

func TestMultipleExporter(t *testing.T) {
	var b1, b2 bytes.Buffer

	e1 := MatchExporter(&b1, MatchGELevel(DebugLevel))
	e2 := MatchExporter(&b2, MatchGELevel(ErrorLevel))

	l := NewDefault().WithExporter(MultipleExporter(e1, e2))

	l.Debug().Msg("DebugMessage").Fire()
	l.Info().Msg("InfoMessage").Fire()
	l.Error().Msg("ErrorMessage").Fire()
	l.Fatal().Msg("FatalMessage").Fire()

	s1 := b1.String()
	require.Contains(t, s1, "DebugMessage")
	require.Contains(t, s1, "InfoMessage")
	require.Contains(t, s1, "ErrorMessage")
	require.Contains(t, s1, "FatalMessage")

	s2 := b2.String()
	require.NotContains(t, s2, "DebugMessage")
	require.NotContains(t, s2, "InfoMessage")
	require.Contains(t, s2, "ErrorMessage")
	require.Contains(t, s2, "FatalMessage")
}
