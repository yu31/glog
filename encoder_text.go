package glog

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/DataWorkbench/glog/pkg/buffer"
)

var (
	_textBufferPool = buffer.NewPool()
)

// TextEncoder return a new Encoder implements by textEncoder
func TextEncoder() Encoder { return newTextEncoder() }

func newTextEncoder() *textEncoder {
	enc := &textEncoder{
		buf: _textBufferPool.Get(),
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
	enc.appendElementSeparator()
	enc.appendString(msg)
}
func (enc *textEncoder) AddEntryTime(t time.Time, layout string) {
	enc.appendElementSeparator()
	enc.appendTime(t, layout)
}
func (enc *textEncoder) AddLevel(level Level) {
	enc.appendElementSeparator()
	enc.buf.AppendByte('[')
	enc.appendString(level.String())
	enc.buf.AppendByte(']')
}
func (enc *textEncoder) AddCaller(skip int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return
	}
	enc.appendElementSeparator()
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
	enc.appendElementSeparator()
	_, err := enc.buf.Write(p)
	return err
}

// Implements ObjectEncoder
func (enc *textEncoder) AddByte(k string, b byte)       { enc.appendKey(k); enc.appendByteInt(b) }
func (enc *textEncoder) AddString(k string, s string)   { enc.appendKey(k); enc.appendString(s) }
func (enc *textEncoder) AddBool(k string, v bool)       { enc.appendKey(k); enc.appendBool(v) }
func (enc *textEncoder) AddInt64(k string, i int64)     { enc.appendKey(k); enc.appendInt64(i) }
func (enc *textEncoder) AddUnt64(k string, i uint64)    { enc.appendKey(k); enc.appendUint64(i) }
func (enc *textEncoder) AddFloat64(k string, f float64) { enc.appendKey(k); enc.appendFloat(f) }
func (enc *textEncoder) AddComplex128(k string, c complex128) {
	enc.appendKey(k)
	enc.appendComplex128(c)
}
func (enc *textEncoder) AddRawBytes(k string, bs []byte) { enc.appendKey(k); enc.appendRawBytes(bs) }
func (enc *textEncoder) AddRawString(k string, s string) { enc.appendKey(k); enc.appendRawString(s) }
func (enc *textEncoder) AddTime(k string, t time.Time, layout string) {
	enc.appendKey(k)
	enc.appendTime(t, layout)
}
func (enc *textEncoder) AddDuration(k string, d time.Duration, layout int8) {
	enc.appendKey(k)
	enc.appendDuration(d, layout)
}
func (enc *textEncoder) AddArray(k string, am ArrayMarshaler) error {
	enc.appendKey(k)
	return enc.appendArray(am)
}
func (enc *textEncoder) AddObject(k string, om ObjectMarshaler) error {
	enc.appendKey(k)
	return enc.appendObject(om)
}
func (enc *textEncoder) AddInterface(k string, i interface{}) error {
	enc.appendKey(k)
	return enc.appendInterface(i)
}

// Implements FieldEncoder
func (enc *textEncoder) AppendByte(v byte)       { enc.appendElementSeparator(); enc.appendByteInt(v) }
func (enc *textEncoder) AppendString(s string)   { enc.appendElementSeparator(); enc.appendString(s) }
func (enc *textEncoder) AppendBool(v bool)       { enc.appendElementSeparator(); enc.appendBool(v) }
func (enc *textEncoder) AppendInt64(i int64)     { enc.appendElementSeparator(); enc.appendInt64(i) }
func (enc *textEncoder) AppendUnt64(i uint64)    { enc.appendElementSeparator(); enc.appendUint64(i) }
func (enc *textEncoder) AppendFloat64(f float64) { enc.appendElementSeparator(); enc.appendFloat(f) }
func (enc *textEncoder) AppendComplex128(c complex128) {
	enc.appendElementSeparator()
	enc.appendComplex128(c)
}
func (enc *textEncoder) AppendRawBytes(bs []byte) {
	enc.appendElementSeparator()
	enc.appendRawBytes(bs)
}
func (enc *textEncoder) AppendRawString(s string) {
	enc.appendElementSeparator()
	enc.appendRawString(s)
}
func (enc *textEncoder) AppendDuration(d time.Duration, layout int8) {
	enc.appendElementSeparator()
	enc.appendDuration(d, layout)
}
func (enc *textEncoder) AppendTime(t time.Time, layout string) {
	enc.appendElementSeparator()
	enc.appendTime(t, layout)
}
func (enc *textEncoder) AppendArray(am ArrayMarshaler) error {
	enc.appendElementSeparator()
	return enc.appendArray(am)
}
func (enc *textEncoder) AppendObject(om ObjectMarshaler) error {
	enc.appendElementSeparator()
	return enc.appendObject(om)
}
func (enc *textEncoder) AppendInterface(i interface{}) error {
	enc.appendElementSeparator()
	return enc.appendInterface(i)
}

// Add k between ElementSeparator and FieldSeparator
func (enc *textEncoder) appendKey(key string) {
	enc.appendElementSeparator()
	enc.appendString(key)
	enc.appendFieldSeparator()
}

func (enc *textEncoder) appendFieldSeparator() {
	enc.buf.AppendByte('=')
}

// Add elements separator
func (enc *textEncoder) appendElementSeparator() {
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

func (enc *textEncoder) appendRawBytes(bs []byte) {
	_, _ = enc.buf.Write(bs)
}

func (enc *textEncoder) appendRawString(s string) {
	enc.buf.AppendString(s)
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
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
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
