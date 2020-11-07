package glog

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/DataWorkbench/glog/pkg/buffer"
)

var (
	tp = buffer.NewPool()
)

// TextEncoder return a new Encoder implements by textEncoder
func TextEncoder() Encoder { return newTextEncoder() }

func newTextEncoder() *textEncoder {
	enc := &textEncoder{
		buf: tp.Get(),
	}
	return enc
}

type textEncoder struct {
	buf *buffer.Buffer
}

// Implements Encoder
func (enc *textEncoder) Bytes() []byte {
	return enc.buf.Bytes()
}
func (enc *textEncoder) Close() error {
	enc.buf.Free()
	enc.buf = nil
	return nil
}

// Implements BuildEncoder
func (enc *textEncoder) AddBeginMarker() {}
func (enc *textEncoder) AddEndMarker()   {}
func (enc *textEncoder) AddLineBreak()   { enc.buf.AppendByte('\n') }
func (enc *textEncoder) AddMsg(msg string) {
	enc.addElementSeparator()
	enc.appendString(msg)
}
func (enc *textEncoder) AddEntryTime(t time.Time, layout string) {
	enc.addElementSeparator()
	enc.appendTime(t, layout)
}
func (enc *textEncoder) AddLevel(level Level) {
	enc.addElementSeparator()
	enc.buf.AppendByte('[')
	enc.appendString(level.String())
	enc.buf.AppendByte(']')
}
func (enc *textEncoder) AddCaller(skip int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return
	}
	enc.addElementSeparator()
	enc.buf.AppendByte('(')
	enc.appendString(file)
	enc.buf.AppendByte(':')
	enc.buf.AppendInt(int64(line))
	enc.buf.AppendByte(')')

}
func (enc *textEncoder) WriteIn(p []byte) error {
	if len(p) == 0 {
		return nil
	}
	enc.addElementSeparator()
	_, err := enc.buf.Write(p)
	return err
}

// Implements ObjectEncoder
func (enc *textEncoder) AddByte(key string, b byte)       { enc.appendKey(key); enc.appendByteInt(b) }
func (enc *textEncoder) AddString(key string, s string)   { enc.appendKey(key); enc.appendString(s) }
func (enc *textEncoder) AddBool(key string, v bool)       { enc.appendKey(key); enc.appendBool(v) }
func (enc *textEncoder) AddInt64(key string, i int64)     { enc.appendKey(key); enc.appendInt64(i) }
func (enc *textEncoder) AddUnt64(key string, i uint64)    { enc.appendKey(key); enc.appendUint64(i) }
func (enc *textEncoder) AddFloat64(key string, f float64) { enc.appendKey(key); enc.appendFloat(f) }
func (enc *textEncoder) AddComplex128(key string, c complex128) {
	enc.appendKey(key)
	enc.appendComplex128(c)
}
func (enc *textEncoder) AddTime(key string, t time.Time, layout string) {
	enc.appendKey(key)
	enc.appendTime(t, layout)
}
func (enc *textEncoder) AddDuration(key string, d time.Duration, layout int8) {
	enc.appendKey(key)
	enc.appendDuration(d, layout)
}
func (enc *textEncoder) AddArray(key string, am ArrayMarshaler) error {
	enc.appendKey(key)
	return enc.appendArray(am)
}
func (enc *textEncoder) AddObject(key string, om ObjectMarshaler) error {
	enc.appendKey(key)
	return enc.appendObject(om)
}
func (enc *textEncoder) AddInterface(key string, i interface{}) error {
	enc.appendKey(key)
	return enc.appendInterface(i)
}

// Implements FieldEncoder
func (enc *textEncoder) AppendByte(v byte)       { enc.addElementSeparator(); enc.appendByteInt(v) }
func (enc *textEncoder) AppendString(s string)   { enc.addElementSeparator(); enc.appendString(s) }
func (enc *textEncoder) AppendBool(v bool)       { enc.addElementSeparator(); enc.appendBool(v) }
func (enc *textEncoder) AppendInt64(i int64)     { enc.addElementSeparator(); enc.appendInt64(i) }
func (enc *textEncoder) AppendUnt64(i uint64)    { enc.addElementSeparator(); enc.appendUint64(i) }
func (enc *textEncoder) AppendFloat64(f float64) { enc.addElementSeparator(); enc.appendFloat(f) }
func (enc *textEncoder) AppendComplex128(c complex128) {
	enc.addElementSeparator()
	enc.appendComplex128(c)
}
func (enc *textEncoder) AppendDuration(d time.Duration, layout int8) {
	enc.addElementSeparator()
	enc.appendDuration(d, layout)
}
func (enc *textEncoder) AppendTime(t time.Time, layout string) {
	enc.addElementSeparator()
	enc.appendTime(t, layout)
}
func (enc *textEncoder) AppendArray(am ArrayMarshaler) error {
	enc.addElementSeparator()
	return enc.appendArray(am)
}
func (enc *textEncoder) AppendObject(om ObjectMarshaler) error {
	enc.addElementSeparator()
	return enc.appendObject(om)
}
func (enc *textEncoder) AppendInterface(i interface{}) error {
	enc.addElementSeparator()
	return enc.appendInterface(i)
}

// build buffer
func (enc *textEncoder) appendKey(key string) {
	enc.addElementSeparator()
	enc.appendString(key)
	enc.addFieldSeparator()
}

func (enc *textEncoder) addFieldSeparator() {
	enc.buf.AppendByte('=')
}

// Add elements separator
func (enc *textEncoder) addElementSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}

	switch enc.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ':
		return
	default:
		enc.buf.AppendByte(' ')
	}
}

func (enc *textEncoder) appendString(s string) {
	AppendStringEscape(enc.buf, s)
}

func (enc *textEncoder) appendByteInt(b byte) {
	enc.buf.AppendUint(uint64(b))
}

func (enc *textEncoder) appendInt64(i int64) {
	enc.buf.AppendInt(i)
}

func (enc *textEncoder) appendUint64(i uint64) {
	enc.buf.AppendUint(i)
}

func (enc *textEncoder) appendBool(v bool) {
	enc.buf.AppendBool(v)
}

func (enc *textEncoder) appendTime(t time.Time, layout string) {
	switch layout {
	case TimeFormatUnixSecond:
		enc.buf.AppendInt(t.Unix())
	case TimeFormatUnixMilli:
		enc.buf.AppendInt(t.UnixNano() / 1e6)
	case TimeFormatUnixMicro:
		enc.buf.AppendInt(t.UnixNano() / 1e3)
	case TimeFormatUnixNano:
		enc.buf.AppendInt(t.UnixNano())
	default:
		enc.buf.AppendTime(t, layout)
	}
}

func (enc *textEncoder) appendDuration(d time.Duration, layout int8) {
	AppendDuration(enc.buf, d, layout)
}

func (enc *textEncoder) appendFloat(f float64) {
	switch {
	case math.IsNaN(f):
		enc.buf.AppendString(`NaN`)
	case math.IsInf(f, 1):
		enc.buf.AppendString(`+Inf`)
	case math.IsInf(f, -1):
		enc.buf.AppendString(`-Inf`)
	default:
		enc.buf.AppendFloat(f, 64)
	}
}

func (enc *textEncoder) appendComplex128(c complex128) {
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(c)), float64(imag(c))
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendByte('"')
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *textEncoder) appendArray(am ArrayMarshaler) error {
	enc.buf.AppendByte('[')
	err := am.MarshalGLogArray(enc)
	enc.buf.AppendByte(']')
	return err
}

func (enc *textEncoder) appendObject(om ObjectMarshaler) error {
	enc.buf.AppendByte('{')
	err := om.MarshalGLogObject(enc)
	enc.buf.AppendByte('}')
	return err
}

func (enc *textEncoder) appendInterface(i interface{}) error {
	switch i.(type) {
	case nil:
		enc.buf.AppendString("<nil>")
		return nil
	default:
		enc.appendString(fmt.Sprintf("%+v", i))
	}
	return nil
}
