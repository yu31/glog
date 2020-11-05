package glog

import (
	"fmt"
	"time"
)

// Entry used to build a log record
type Entry struct {
	Level   Level
	Encoder Encoder

	l *Logger
}

// newEntry will create a new entry with level and fields.
func newEntry(l *Logger, level Level) *Entry {
	e := &Entry{
		l:       l,
		Level:   level,
		Encoder: l.encoderFunc(),
	}
	e.encodeHeads()
	return e
}

// withError handle any error if happen in entry inside
func (e *Entry) withError(err error) {
	if err == nil {
		return
	}
	_, _ = fmt.Fprintf(e.l.errorOutput, "entry error: %v\n", err)
}

func (e *Entry) encodeHeads() {
	e.Encoder.AddBeginMarker()
	e.Encoder.AddEntryTime(time.Now(), e.l.timeLayout)
	e.Encoder.AddLevel(e.Level)
}

func (e *Entry) encodeEnds() {
	e.withError(e.Encoder.WriteIn(e.l.fields.Bytes()))
	if e.l.caller {
		e.Encoder.AddCaller(2)
	}
	e.Encoder.AddEndMarker()
	e.Encoder.AddLineBreak()
}

func (e *Entry) free() {
	if e == nil {
		return
	}
	e.withError(e.Encoder.Close())
	e.l = nil
	e.Encoder = nil
}

// Fire sends the *Entry to Logger's executor.
//
// NOTICE: once this method is called, the *Entry should be disposed.
// Calling Fire twice can have unexpected result.
func (e *Entry) Fire() {
	if e == nil {
		return
	}
	e.encodeEnds()
	// NOTICE: once the `Execute` returns, the *Entry should be disposed,
	// if not can have unexpected result.
	e.withError(e.l.executor.Execute(e))
	// Release resources
	e.free()
}

func (e *Entry) Msg(msg string) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddMsg(msg)
	return e
}

func (e *Entry) Duration(key string, val time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(key, int64(val))
	return e
}

func (e *Entry) Array(key string, arr ArrayMarshaler) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, arr))
	return e
}

func (e *Entry) Object(key string, obj ObjectMarshaler) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddObject(key, obj))
	return e
}

func (e *Entry) Byte(key string, val byte) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddByte(key, val)
	return e
}

func (e *Entry) Bytes(key string, val []byte) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, byteArray(val)))
	return e
}

func (e *Entry) String(key string, val string) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddString(key, val)
	return e
}

func (e *Entry) Strings(key string, val []string) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, stringArray(val)))
	return e
}

func (e *Entry) Bool(key string, val bool) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddBool(key, val)
	return e
}

func (e *Entry) Boos(key string, val []bool) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, bools(val)))
	return e
}

func (e *Entry) Uintptr(key string, val uintptr) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(key, uint64(val))
	return e
}

func (e *Entry) Uintptrs(key string, val []uintptr) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, uintptrs(val)))
	return e
}

func (e *Entry) Rune(key string, val rune) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(key, int64(val))
	return e
}

func (e *Entry) Runes(key string, val []rune) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, runes(val)))
	return e
}

func (e *Entry) Int(key string, val int) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(key, int64(val))
	return e
}

func (e *Entry) Ints(key string, val []int) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, ints(val)))
	return e
}

func (e *Entry) Int8(key string, val int8) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(key, int64(val))
	return e
}

func (e *Entry) Int8s(key string, val []int8) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, int8s(val)))
	return e
}

func (e *Entry) Int16(key string, val int16) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(key, int64(val))
	return e
}

func (e *Entry) Int16s(key string, val []int16) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, int16s(val)))
	return e
}

func (e *Entry) Int32(key string, val int32) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(key, int64(val))
	return e
}

func (e *Entry) Int32s(key string, val []int32) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, int32s(val)))
	return e
}

func (e *Entry) Int64(key string, val int64) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(key, int64(val))
	return e
}

func (e *Entry) Int64s(key string, val []int64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, int64s(val)))
	return e
}

func (e *Entry) Uint(key string, val uint) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(key, uint64(val))
	return e
}

func (e *Entry) Uints(key string, val []uint) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, uints(val)))
	return e
}

func (e *Entry) Uint8(key string, val uint8) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(key, uint64(val))
	return e
}

func (e *Entry) Uint8s(key string, val []uint8) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, uint8s(val)))
	return e
}

func (e *Entry) Uint16(key string, val uint16) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(key, uint64(val))
	return e
}

func (e *Entry) Uint16s(key string, val []uint16) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, uint16s(val)))
	return e
}

func (e *Entry) Uint32(key string, val uint32) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(key, uint64(val))
	return e
}

func (e *Entry) Uint32s(key string, val []uint32) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, uint32s(val)))
	return e
}

func (e *Entry) Uint64(key string, val uint64) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(key, uint64(val))
	return e
}

func (e *Entry) Uint64s(key string, val []uint64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, uint64s(val)))
	return e
}

func (e *Entry) Float32(key string, val float32) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddFloat64(key, float64(val))
	return e
}

func (e *Entry) Float32s(key string, val []float32) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, float32s(val)))
	return e
}

func (e *Entry) Float64(key string, val float64) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddFloat64(key, val)
	return e
}

func (e *Entry) Float64s(key string, val []float64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, float64s(val)))
	return e
}

func (e *Entry) Complex64(key string, val complex64) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddComplex128(key, complex128(val))
	return e
}

func (e *Entry) Complex64s(key string, val []complex64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, complex64s(val)))
	return e
}

func (e *Entry) Complex128(key string, val complex128) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddComplex128(key, val)
	return e
}

func (e *Entry) Complex128s(key string, val []complex128) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, complex128s(val)))
	return e
}

func (e *Entry) Time(key string, val time.Time, layout string) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddTime(key, val, layout)
	return e
}

func (e *Entry) Error(key string, err error) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddString(key, err.Error())
	return e
}

// Any uses reflection to serialize arbitrary objects, so it can be
// slow and allocation-heavy.
func (e *Entry) Any(key string, val interface{}) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddInterface(key, val))
	return e
}
