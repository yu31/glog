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
func (enc *jsonEncoder) AddByte(key string, val byte)       { enc.appendKey(key); enc.appendByteRaw(val) }
func (enc *jsonEncoder) AddString(key string, val string)   { enc.appendKey(key); enc.appendString(val) }
func (enc *jsonEncoder) AddBool(key string, val bool)       { enc.appendKey(key); enc.appendBool(val) }
func (enc *jsonEncoder) AddInt64(key string, val int64)     { enc.appendKey(key); enc.appendInt64(val) }
func (enc *jsonEncoder) AddUnt64(key string, val uint64)    { enc.appendKey(key); enc.appendUint64(val) }
func (enc *jsonEncoder) AddFloat64(key string, val float64) { enc.appendKey(key); enc.appendFloat(val) }
func (enc *jsonEncoder) AddTime(key string, val time.Time, layout string) {
	enc.appendKey(key)
	enc.appendTime(val, layout)
}
func (enc *jsonEncoder) AddComplex128(key string, val complex128) {
	enc.appendKey(key)
	enc.appendComplex128(val)
}
func (enc *jsonEncoder) AddArray(key string, arr ArrayMarshaler) error {
	enc.appendKey(key)
	return enc.appendArray(arr)
}
func (enc *jsonEncoder) AddObject(key string, obj ObjectMarshaler) error {
	enc.appendKey(key)
	return enc.appendObject(obj)
}
func (enc *jsonEncoder) AddInterface(key string, val interface{}) error {
	enc.appendKey(key)
	return enc.appendInterface(val)
}

// Implements FieldEncoder
func (enc *jsonEncoder) AppendByte(val byte)       { enc.addElementSeparator(); enc.appendByteRaw(val) }
func (enc *jsonEncoder) AppendString(val string)   { enc.addElementSeparator(); enc.appendString(val) }
func (enc *jsonEncoder) AppendBool(val bool)       { enc.addElementSeparator(); enc.appendBool(val) }
func (enc *jsonEncoder) AppendInt64(val int64)     { enc.addElementSeparator(); enc.appendInt64(val) }
func (enc *jsonEncoder) AppendUnt64(val uint64)    { enc.addElementSeparator(); enc.appendUint64(val) }
func (enc *jsonEncoder) AppendFloat64(val float64) { enc.addElementSeparator(); enc.appendFloat(val) }
func (enc *jsonEncoder) AppendTime(val time.Time, layout string) {
	enc.addElementSeparator()
	enc.appendTime(val, layout)
}
func (enc *jsonEncoder) AppendComplex128(val complex128) {
	enc.addElementSeparator()
	enc.appendComplex128(val)
}
func (enc *jsonEncoder) AppendArray(arr ArrayMarshaler) error   { return enc.appendArray(arr) }
func (enc *jsonEncoder) AppendObject(obj ObjectMarshaler) error { return enc.appendObject(obj) }
func (enc *jsonEncoder) AppendInterface(val interface{}) error  { return enc.appendInterface(val) }

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

func (enc *jsonEncoder) appendString(val string) {
	enc.buf.AppendByte('"')
	AppendStringEscape(enc.buf, val)
	enc.buf.AppendByte('"')
}

func (enc *jsonEncoder) appendByteRaw(val byte) {
	enc.buf.AppendUint(uint64(val))
}

func (enc *jsonEncoder) appendInt64(val int64) {
	enc.buf.AppendInt(val)
}

func (enc *jsonEncoder) appendUint64(val uint64) {
	enc.buf.AppendUint(val)
}

func (enc *jsonEncoder) appendTime(t time.Time, layout string) {
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
		enc.buf.AppendByte('"')
		enc.buf.AppendTime(t, layout)
		enc.buf.AppendByte('"')
	}
}

func (enc *jsonEncoder) appendFloat(val float64) {
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

func (enc *jsonEncoder) appendBool(val bool) {
	enc.buf.AppendBool(val)
}

func (enc *jsonEncoder) appendComplex128(val complex128) {
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *jsonEncoder) appendArray(arr ArrayMarshaler) error {
	enc.buf.AppendByte('[')
	err := arr.MarshalLogArray(enc)
	enc.buf.AppendByte(']')
	return err
}

func (enc *jsonEncoder) appendObject(obj ObjectMarshaler) error {
	enc.buf.AppendByte('{')
	err := obj.MarshalLogObject(enc)
	enc.buf.AppendByte('}')
	return err
}

func (enc *jsonEncoder) appendInterface(val interface{}) error {
	b, err := json.Marshal(val)
	if err != nil {
		enc.appendString(fmt.Sprintf("%v", val))
		return err
	}
	_, _ = enc.buf.Write(b)
	return nil
}
