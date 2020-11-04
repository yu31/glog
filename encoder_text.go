package glog

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
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
	enc.appendString(file + ":" + strconv.Itoa(line))
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
func (enc *textEncoder) AddByte(key string, val byte)       { enc.appendKey(key); enc.appendByteRaw(val) }
func (enc *textEncoder) AddString(key string, val string)   { enc.appendKey(key); enc.appendString(val) }
func (enc *textEncoder) AddBool(key string, val bool)       { enc.appendKey(key); enc.appendBool(val) }
func (enc *textEncoder) AddInt64(key string, val int64)     { enc.appendKey(key); enc.appendInt64(val) }
func (enc *textEncoder) AddUnt64(key string, val uint64)    { enc.appendKey(key); enc.appendUint64(val) }
func (enc *textEncoder) AddFloat64(key string, val float64) { enc.appendKey(key); enc.appendFloat(val) }
func (enc *textEncoder) AddTime(key string, val time.Time, layout string) {
	enc.appendKey(key)
	enc.appendTime(val, layout)
}
func (enc *textEncoder) AddComplex128(key string, val complex128) {
	enc.appendKey(key)
	enc.appendComplex128(val)
}
func (enc *textEncoder) AddArray(key string, arr ArrayMarshaler) error {
	enc.appendKey(key)
	return enc.appendArray(arr)
}
func (enc *textEncoder) AddObject(key string, obj ObjectMarshaler) error {
	enc.appendKey(key)
	return enc.appendObject(obj)
}
func (enc *textEncoder) AddInterface(key string, val interface{}) error {
	enc.appendKey(key)
	return enc.appendInterface(val)
}

// Implements FieldEncoder
func (enc *textEncoder) AppendByte(val byte)       { enc.addElementSeparator(); enc.appendByteRaw(val) }
func (enc *textEncoder) AppendString(val string)   { enc.addElementSeparator(); enc.appendString(val) }
func (enc *textEncoder) AppendBool(val bool)       { enc.addElementSeparator(); enc.appendBool(val) }
func (enc *textEncoder) AppendInt64(val int64)     { enc.addElementSeparator(); enc.appendInt64(val) }
func (enc *textEncoder) AppendUnt64(val uint64)    { enc.addElementSeparator(); enc.appendUint64(val) }
func (enc *textEncoder) AppendFloat64(val float64) { enc.addElementSeparator(); enc.appendFloat(val) }
func (enc *textEncoder) AppendTime(val time.Time, layout string) {
	enc.addElementSeparator()
	enc.appendTime(val, layout)
}
func (enc *textEncoder) AppendComplex128(val complex128) {
	enc.addElementSeparator()
	enc.appendComplex128(val)
}
func (enc *textEncoder) AppendArray(arr ArrayMarshaler) error   { return enc.appendArray(arr) }
func (enc *textEncoder) AppendObject(obj ObjectMarshaler) error { return enc.appendObject(obj) }
func (enc *textEncoder) AppendInterface(val interface{}) error  { return enc.appendInterface(val) }

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

func (enc *textEncoder) appendString(val string) {
	AppendStringEscape(enc.buf, val)
}

func (enc *textEncoder) appendByteRaw(val byte) {
	enc.buf.AppendUint(uint64(val))
}

func (enc *textEncoder) appendInt64(val int64) {
	enc.buf.AppendInt(val)
}

func (enc *textEncoder) appendUint64(val uint64) {
	enc.buf.AppendUint(val)
}

func (enc *textEncoder) appendBool(val bool) {
	enc.buf.AppendBool(val)
}

func (enc *textEncoder) appendTime(t time.Time, layout string) {
	switch layout {
	case TimeFormatUnix:
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

func (enc *textEncoder) appendFloat(val float64) {
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, 64)
	}
}

func (enc *textEncoder) appendComplex128(val complex128) {
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendByte('"')
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *textEncoder) appendArray(arr ArrayMarshaler) error {
	enc.buf.AppendByte('[')
	err := arr.MarshalLogArray(enc)
	enc.buf.AppendByte(']')
	return err
}

func (enc *textEncoder) appendObject(obj ObjectMarshaler) error {
	enc.buf.AppendByte('{')
	err := obj.MarshalLogObject(enc)
	enc.buf.AppendByte('}')
	return err
}

func (enc *textEncoder) appendInterface(val interface{}) error {
	enc.appendString(fmt.Sprintf("%v", val))
	return nil
}
