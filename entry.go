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

// Nanosecond encode time.Duration to an string nanoseconds; format sample "1004854348ns".
func (e *Entry) Nanosecond(key string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddDuration(key, d, DurationFormatNano)
	//e.Encoder.AddDuration(key, durationNano(val))
	return e
}

// Microsecond encode time.Duration to an string microseconds, format sample "1004854us".
func (e *Entry) Microsecond(key string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddDuration(key, d, DurationFormatMicro)
	//e.Encoder.AddDuration(key, durationMicro(val))
	return e
}

// Millisecond encode time.Duration to an string milliseconds, format sample "1004ms".
func (e *Entry) Millisecond(key string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddDuration(key, d, DurationFormatMilli)
	//e.Encoder.AddDuration(key, durationMilli(val))
	return e
}

// Second encode time.Duration to an string seconds, format sample "1.004854348s".
func (e *Entry) Second(key string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddDuration(key, d, DurationFormatSecond)
	//e.Encoder.AddDuration(key, durationSecond(val))
	return e
}

// Minute encode time.Duration to an string minutes, format sample "10min".
func (e *Entry) Minute(key string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddDuration(key, d, DurationFormatMinute)
	//e.Encoder.AddDuration(key, durationMinute(val))
	return e
}

// Hour encode time.Duration to an string hours, format sample "2h".
func (e *Entry) Hour(key string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddDuration(key, d, DurationFormatHour)
	//e.Encoder.AddDuration(key, durationHour(val))
	return e
}

// Byte encode the value to an integer number.
func (e *Entry) Byte(key string, b byte) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddByte(key, b)
	return e
}

// Bytes encode the value to an integer array.
func (e *Entry) Bytes(key string, bb []byte) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, byteArray(bb)))
	return e
}

func (e *Entry) String(key string, s string) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddString(key, s)
	return e
}

func (e *Entry) Strings(key string, ss []string) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, stringArray(ss)))
	return e
}

func (e *Entry) Bool(key string, v bool) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddBool(key, v)
	return e
}

func (e *Entry) Bools(key string, vv []bool) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, bools(vv)))
	return e
}

func (e *Entry) Int(key string, i int) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(key, int64(i))
	return e
}

func (e *Entry) Ints(key string, ii []int) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, ints(ii)))
	return e
}

func (e *Entry) Int8(key string, i int8) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(key, int64(i))
	return e
}

func (e *Entry) Int8s(key string, ii []int8) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, int8s(ii)))
	return e
}

func (e *Entry) Int16(key string, i int16) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(key, int64(i))
	return e
}

func (e *Entry) Int16s(key string, ii []int16) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, int16s(ii)))
	return e
}

func (e *Entry) Int32(key string, i int32) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(key, int64(i))
	return e
}

func (e *Entry) Int32s(key string, ii []int32) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, int32s(ii)))
	return e
}

func (e *Entry) Int64(key string, i int64) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(key, int64(i))
	return e
}

func (e *Entry) Int64s(key string, ii []int64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, int64s(ii)))
	return e
}

func (e *Entry) Uint(key string, i uint) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(key, uint64(i))
	return e
}

func (e *Entry) Uints(key string, ii []uint) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, uints(ii)))
	return e
}

func (e *Entry) Uint8(key string, i uint8) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(key, uint64(i))
	return e
}

func (e *Entry) Uint8s(key string, ii []uint8) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, uint8s(ii)))
	return e
}

func (e *Entry) Uint16(key string, i uint16) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(key, uint64(i))
	return e
}

func (e *Entry) Uint16s(key string, ii []uint16) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, uint16s(ii)))
	return e
}

func (e *Entry) Uint32(key string, i uint32) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(key, uint64(i))
	return e
}

func (e *Entry) Uint32s(key string, ii []uint32) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, uint32s(ii)))
	return e
}

func (e *Entry) Uint64(key string, i uint64) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(key, uint64(i))
	return e
}

func (e *Entry) Uint64s(key string, ii []uint64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, uint64s(ii)))
	return e
}

func (e *Entry) Float32(key string, f float32) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddFloat64(key, float64(f))
	return e
}

func (e *Entry) Float32s(key string, ff []float32) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, float32s(ff)))
	return e
}

func (e *Entry) Float64(key string, f float64) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddFloat64(key, f)
	return e
}

func (e *Entry) Float64s(key string, ff []float64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, float64s(ff)))
	return e
}

func (e *Entry) Complex64(key string, c complex64) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddComplex128(key, complex128(c))
	return e
}

func (e *Entry) Complex64s(key string, cc []complex64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, complex64s(cc)))
	return e
}

func (e *Entry) Complex128(key string, c complex128) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddComplex128(key, c)
	return e
}

func (e *Entry) Complex128s(key string, cc []complex128) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, complex128s(cc)))
	return e
}

func (e *Entry) Error(key string, err error) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddString(key, err.Error())
	return e
}

func (e *Entry) Errors(key string, errs []error) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, errorArray(errs)))
	return e
}

func (e *Entry) Time(key string, t time.Time, layout string) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddTime(key, t, layout)
	return e
}

func (e *Entry) Array(key string, am ArrayMarshaler) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(key, am))
	return e
}

func (e *Entry) Object(key string, om ObjectMarshaler) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddObject(key, om))
	return e
}

// Any uses reflection to serialize arbitrary objects, so it can be
// slow and allocation-heavy.
func (e *Entry) Any(key string, i interface{}) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddInterface(key, i))
	return e
}
