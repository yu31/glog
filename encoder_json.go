package glog

import (
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/DataWorkbench/glog/pkg/buffer"
)

var (
	_jsonBufferPool = buffer.NewPool()
)

// JSONEncoder return a new Encoder implements by jsonEncoder.
func JSONEncoder() Encoder { return newJSONEncoder() }

func newJSONEncoder() *jsonEncoder {
	enc := &jsonEncoder{
		buf: _jsonBufferPool.Get(),
	}
	return enc
}

type jsonEncoder struct {
	buf *buffer.Buffer
}

// Implements Encoder.
func (enc *jsonEncoder) Bytes() []byte {
	return enc.buf.Bytes()
}
func (enc *jsonEncoder) Close() error {
	enc.buf.Free()
	enc.buf = nil
	return nil
}

// Implements BuildEncoder.
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
	enc.appendElementSeparator()
	enc.buf.AppendByte('"')
	AppendStringEscape(enc.buf, file)
	enc.buf.AppendByte(':')
	enc.buf.AppendInt(int64(line))
	enc.buf.AppendByte('"')

	//enc.appendString(file + ":" + strconv.Itoa(line))
}
func (enc *jsonEncoder) WriteIn(p []byte) error {
	if len(p) == 0 {
		return nil
	}
	_, err := enc.buf.Write(p)
	return err
}

// Implements ObjectEncoder.
func (enc *jsonEncoder) AddByte(k string, b byte)       { enc.appendKey(k); enc.appendByteInt(b) }
func (enc *jsonEncoder) AddString(k string, s string)   { enc.appendKey(k); enc.appendString(s) }
func (enc *jsonEncoder) AddBool(k string, v bool)       { enc.appendKey(k); enc.appendBool(v) }
func (enc *jsonEncoder) AddInt64(k string, i int64)     { enc.appendKey(k); enc.appendInt64(i) }
func (enc *jsonEncoder) AddUnt64(k string, i uint64)    { enc.appendKey(k); enc.appendUint64(i) }
func (enc *jsonEncoder) AddFloat64(k string, f float64) { enc.appendKey(k); enc.appendFloat(f) }
func (enc *jsonEncoder) AddComplex128(k string, c complex128) {
	enc.appendKey(k)
	enc.appendComplex128(c)
}
func (enc *jsonEncoder) AddRawBytes(k string, bs []byte) { enc.appendKey(k); enc.appendRawBytes(bs) }
func (enc *jsonEncoder) AddRawString(k string, s string) { enc.appendKey(k); enc.appendRawString(s) }
func (enc *jsonEncoder) AddTime(k string, t time.Time, layout string) {
	enc.appendKey(k)
	enc.appendTime(t, layout)
}
func (enc *jsonEncoder) AddDuration(k string, d time.Duration, layout int8) {
	enc.appendKey(k)
	enc.appendDuration(d, layout)
}
func (enc *jsonEncoder) AddArray(k string, am ArrayMarshaler) error {
	enc.appendKey(k)
	return enc.appendArray(am)
}
func (enc *jsonEncoder) AddObject(k string, om ObjectMarshaler) error {
	enc.appendKey(k)
	return enc.appendObject(om)
}
func (enc *jsonEncoder) AddInterface(k string, i interface{}) error {
	enc.appendKey(k)
	return enc.appendInterface(i)
}

// Implements FieldEncoder.
func (enc *jsonEncoder) AppendByte(b byte)       { enc.appendElementSeparator(); enc.appendByteInt(b) }
func (enc *jsonEncoder) AppendString(s string)   { enc.appendElementSeparator(); enc.appendString(s) }
func (enc *jsonEncoder) AppendBool(v bool)       { enc.appendElementSeparator(); enc.appendBool(v) }
func (enc *jsonEncoder) AppendInt64(i int64)     { enc.appendElementSeparator(); enc.appendInt64(i) }
func (enc *jsonEncoder) AppendUnt64(i uint64)    { enc.appendElementSeparator(); enc.appendUint64(i) }
func (enc *jsonEncoder) AppendFloat64(f float64) { enc.appendElementSeparator(); enc.appendFloat(f) }
func (enc *jsonEncoder) AppendComplex128(c complex128) {
	enc.appendElementSeparator()
	enc.appendComplex128(c)
}
func (enc *jsonEncoder) AppendRawBytes(bs []byte) {
	enc.appendElementSeparator()
	enc.appendRawBytes(bs)
}
func (enc *jsonEncoder) AppendRawString(s string) {
	enc.appendElementSeparator()
	enc.appendRawString(s)
}
func (enc *jsonEncoder) AppendDuration(d time.Duration, layout int8) {
	enc.appendElementSeparator()
	enc.appendDuration(d, layout)
}
func (enc *jsonEncoder) AppendTime(t time.Time, layout string) {
	enc.appendElementSeparator()
	enc.appendTime(t, layout)
}
func (enc *jsonEncoder) AppendArray(am ArrayMarshaler) error {
	enc.appendElementSeparator()
	return enc.appendArray(am)
}
func (enc *jsonEncoder) AppendObject(om ObjectMarshaler) error {
	enc.appendElementSeparator()
	return enc.appendObject(om)
}
func (enc *jsonEncoder) AppendInterface(i interface{}) error {
	enc.appendElementSeparator()
	return enc.appendInterface(i)
}

// Add k between ElementSeparator and FieldSeparator.
func (enc *jsonEncoder) appendKey(key string) {
	enc.appendElementSeparator()
	enc.appendString(key)
	enc.appendFieldSeparator()
}

// Add field separator.
func (enc *jsonEncoder) appendFieldSeparator() {
	enc.buf.AppendByte(':')
}

// Add elements separator.
func (enc *jsonEncoder) appendElementSeparator() {
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

func (enc *jsonEncoder) appendRawBytes(bs []byte) {
	_, _ = enc.buf.Write(bs)
}

func (enc *jsonEncoder) appendRawString(s string) {
	enc.buf.AppendString(s)
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
	err := am.MarshalGLogArray(enc)
	enc.buf.AppendByte(']')
	return err
}

func (enc *jsonEncoder) appendObject(om ObjectMarshaler) error {
	enc.buf.AppendByte('{')
	err := om.MarshalGLogObject(enc)
	enc.buf.AppendByte('}')
	return err
}

func (enc *jsonEncoder) appendInterface(i interface{}) error {
	var err error
	var b []byte

	switch m := i.(type) {
	case json.Marshaler:
		b, err = m.MarshalJSON()
	case nil:
		enc.buf.AppendString("null")
		return nil
	default:
		b, err = json.Marshal(i)
	}

	if err != nil {
		enc.appendString(fmt.Sprintf("%+v", i))
		return err
	}
	_, err = enc.buf.Write(b)
	return err
}
