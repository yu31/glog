package glog

import (
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"time"

	"github.com/DataWorkbench/glog/pkg/buffer"
)

var (
	jp = buffer.NewPool()
)

// JSONEncoder return a new Encoder implements by jsonEncoder
func JSONEncoder() Encoder { return newJSONEncoder() }

func newJSONEncoder() *jsonEncoder {
	enc := &jsonEncoder{
		buf: jp.Get(),
	}
	return enc
}

type jsonEncoder struct {
	buf *buffer.Buffer
}

// Implements Encoder
func (enc *jsonEncoder) Bytes() []byte {
	return enc.buf.Bytes()
}
func (enc *jsonEncoder) Close() error {
	enc.buf.Free()
	enc.buf = nil
	return nil
}

// Implements BuildEncoder
func (enc *jsonEncoder) AddBeginMarker() { enc.buf.AppendByte('{') }
func (enc *jsonEncoder) AddEndMarker()   { enc.buf.AppendByte('}') }
func (enc *jsonEncoder) AddLineBreak()   { enc.buf.AppendByte('\n') }
func (enc *jsonEncoder) AddMsg(msg string) {
	enc.appendKey("message")
	enc.appendString(msg)
}
func (enc *jsonEncoder) AddEntryTime(t time.Time, layout string) {
	enc.appendKey("time")
	enc.appendTime(t, layout)
}
func (enc *jsonEncoder) AddLevel(level Level) {
	enc.appendKey("level")
	enc.appendString(level.String())
}
func (enc *jsonEncoder) AddCaller(skip int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return
	}
	enc.appendKey("caller")
	enc.appendString(file + ":" + strconv.Itoa(line))
}
func (enc *jsonEncoder) WriteIn(p []byte) error {
	if len(p) == 0 {
		return nil
	}
	_, err := enc.buf.Write(p)
	return err
}

// Implements ObjectEncoder
func (enc *jsonEncoder) AddByte(key string, b byte)       { enc.appendKey(key); enc.appendByteInt(b) }
func (enc *jsonEncoder) AddString(key string, s string)   { enc.appendKey(key); enc.appendString(s) }
func (enc *jsonEncoder) AddBool(key string, v bool)       { enc.appendKey(key); enc.appendBool(v) }
func (enc *jsonEncoder) AddInt64(key string, i int64)     { enc.appendKey(key); enc.appendInt64(i) }
func (enc *jsonEncoder) AddUnt64(key string, i uint64)    { enc.appendKey(key); enc.appendUint64(i) }
func (enc *jsonEncoder) AddFloat64(key string, f float64) { enc.appendKey(key); enc.appendFloat(f) }
func (enc *jsonEncoder) AddComplex128(key string, c complex128) {
	enc.appendKey(key)
	enc.appendComplex128(c)
}
func (enc *jsonEncoder) AddTime(key string, t time.Time, layout string) {
	enc.appendKey(key)
	enc.appendTime(t, layout)
}
func (enc *jsonEncoder) AddDuration(key string, d time.Duration, layout int8) {
	enc.appendKey(key)
	enc.appendDuration(d, layout)
}
func (enc *jsonEncoder) AddArray(key string, am ArrayMarshaler) error {
	enc.appendKey(key)
	return enc.appendArray(am)
}
func (enc *jsonEncoder) AddObject(key string, om ObjectMarshaler) error {
	enc.appendKey(key)
	return enc.appendObject(om)
}
func (enc *jsonEncoder) AddInterface(key string, i interface{}) error {
	enc.appendKey(key)
	return enc.appendInterface(i)
}

// Implements FieldEncoder
func (enc *jsonEncoder) AppendByte(b byte)       { enc.addElementSeparator(); enc.appendByteInt(b) }
func (enc *jsonEncoder) AppendString(s string)   { enc.addElementSeparator(); enc.appendString(s) }
func (enc *jsonEncoder) AppendBool(v bool)       { enc.addElementSeparator(); enc.appendBool(v) }
func (enc *jsonEncoder) AppendInt64(i int64)     { enc.addElementSeparator(); enc.appendInt64(i) }
func (enc *jsonEncoder) AppendUnt64(i uint64)    { enc.addElementSeparator(); enc.appendUint64(i) }
func (enc *jsonEncoder) AppendFloat64(f float64) { enc.addElementSeparator(); enc.appendFloat(f) }
func (enc *jsonEncoder) AppendComplex128(c complex128) {
	enc.addElementSeparator()
	enc.appendComplex128(c)
}
func (enc *jsonEncoder) AppendDuration(d time.Duration, layout int8) {
	enc.addElementSeparator()
	enc.appendDuration(d, layout)
}
func (enc *jsonEncoder) AppendTime(t time.Time, layout string) {
	enc.addElementSeparator()
	enc.appendTime(t, layout)
}
func (enc *jsonEncoder) AppendArray(am ArrayMarshaler) error   { return enc.appendArray(am) }
func (enc *jsonEncoder) AppendObject(om ObjectMarshaler) error { return enc.appendObject(om) }
func (enc *jsonEncoder) AppendInterface(i interface{}) error   { return enc.appendInterface(i) }

// build buffer
func (enc *jsonEncoder) appendKey(key string) {
	enc.addElementSeparator()
	enc.appendString(key)
	enc.addFieldSeparator()
}

func (enc *jsonEncoder) addFieldSeparator() {
	enc.buf.AppendByte(':')
}

// Add elements separator
func (enc *jsonEncoder) addElementSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}

	switch enc.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ':
		return
	default:
		enc.buf.AppendByte(',')
	}
}

func (enc *jsonEncoder) appendString(s string) {
	enc.buf.AppendByte('"')
	AppendStringEscape(enc.buf, s)
	enc.buf.AppendByte('"')
}

func (enc *jsonEncoder) appendByteInt(b byte) {
	enc.buf.AppendUint(uint64(b))
}

func (enc *jsonEncoder) appendInt64(i int64) {
	enc.buf.AppendInt(i)
}

func (enc *jsonEncoder) appendUint64(i uint64) {
	enc.buf.AppendUint(i)
}

func (enc *jsonEncoder) appendTime(t time.Time, layout string) {
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
		enc.buf.AppendByte('"')
		enc.buf.AppendTime(t, layout)
		enc.buf.AppendByte('"')
	}
}

func (enc *jsonEncoder) appendDuration(d time.Duration, layout int8) {
	enc.buf.AppendByte('"')
	AppendDuration(enc.buf, d, layout)
	enc.buf.AppendByte('"')
}

func (enc *jsonEncoder) appendFloat(f float64) {
	switch {
	case math.IsNaN(f):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(f, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(f, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(f, 64)
	}
}

func (enc *jsonEncoder) appendBool(val bool) {
	enc.buf.AppendBool(val)
}

func (enc *jsonEncoder) appendComplex128(c complex128) {
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(c)), float64(imag(c))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *jsonEncoder) appendArray(am ArrayMarshaler) error {
	enc.buf.AppendByte('[')
	err := am.MarshalArray(enc)
	enc.buf.AppendByte(']')
	return err
}

func (enc *jsonEncoder) appendObject(om ObjectMarshaler) error {
	enc.buf.AppendByte('{')
	err := om.MarshalObject(enc)
	enc.buf.AppendByte('}')
	return err
}

func (enc *jsonEncoder) appendInterface(i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		enc.appendString(fmt.Sprintf("%v", i))
		return err
	}
	_, _ = enc.buf.Write(b)
	return nil
}
