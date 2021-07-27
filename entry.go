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
		Level:   level,
		Encoder: l.encoderFunc(),
		l:       l,
	}
	e.encodeHeads()
	return e
}

// withError handle any error if happen in entry inside
func (e *Entry) withError(err error) {
	if err == nil {
		return
	}
	_, _ = fmt.Fprintf(e.l.errorOutput, "%s [inner] handle log entry error: %v\n", time.Now().Format(e.l.timeLayout), err)
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
	e.Encoder = nil
	e.l = nil
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

// RawBytes adds already serialized data to the log entry under key.
//
// No sanity check is performed on bs; it must not contains carriage returns
// or line break.
func (e *Entry) RawBytes(k string, bs []byte) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddRawBytes(k, bs)
	return e
}

// RawString adds already serialized data to the log entry under key.
//
// No sanity check is performed on s; it must not contains carriage returns
// or line break.
func (e *Entry) RawString(k string, s string) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddRawString(k, s)
	return e
}

// Byte encode the value to an integer format.
func (e *Entry) Byte(k string, b byte) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddByte(k, b)
	return e
}

// Bytes encode the value to an integer array.
func (e *Entry) Bytes(k string, bb []byte) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, byteArray(bb)))
	return e
}

func (e *Entry) String(k string, s string) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddString(k, s)
	return e
}

func (e *Entry) Strings(k string, ss []string) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, stringArray(ss)))
	return e
}

func (e *Entry) Bool(k string, v bool) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddBool(k, v)
	return e
}

func (e *Entry) Bools(k string, vv []bool) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, bools(vv)))
	return e
}

func (e *Entry) Int(k string, i int) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(k, int64(i))
	return e
}

func (e *Entry) Ints(k string, ii []int) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, ints(ii)))
	return e
}

func (e *Entry) Int8(k string, i int8) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(k, int64(i))
	return e
}

func (e *Entry) Int8s(k string, ii []int8) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, int8s(ii)))
	return e
}

func (e *Entry) Int16(k string, i int16) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(k, int64(i))
	return e
}

func (e *Entry) Int16s(k string, ii []int16) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, int16s(ii)))
	return e
}

func (e *Entry) Int32(k string, i int32) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(k, int64(i))
	return e
}

func (e *Entry) Int32s(k string, ii []int32) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, int32s(ii)))
	return e
}

func (e *Entry) Int64(k string, i int64) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddInt64(k, int64(i))
	return e
}

func (e *Entry) Int64s(k string, ii []int64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, int64s(ii)))
	return e
}

func (e *Entry) Uint(k string, i uint) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(k, uint64(i))
	return e
}

func (e *Entry) Uints(k string, ii []uint) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, uints(ii)))
	return e
}

func (e *Entry) Uint8(k string, i uint8) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(k, uint64(i))
	return e
}

func (e *Entry) Uint8s(k string, ii []uint8) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, uint8s(ii)))
	return e
}

func (e *Entry) Uint16(k string, i uint16) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(k, uint64(i))
	return e
}

func (e *Entry) Uint16s(k string, ii []uint16) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, uint16s(ii)))
	return e
}

func (e *Entry) Uint32(k string, i uint32) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(k, uint64(i))
	return e
}

func (e *Entry) Uint32s(k string, ii []uint32) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, uint32s(ii)))
	return e
}

func (e *Entry) Uint64(k string, i uint64) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddUnt64(k, uint64(i))
	return e
}

func (e *Entry) Uint64s(k string, ii []uint64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, uint64s(ii)))
	return e
}

func (e *Entry) Float32(k string, f float32) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddFloat64(k, float64(f))
	return e
}

func (e *Entry) Float32s(k string, ff []float32) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, float32s(ff)))
	return e
}

func (e *Entry) Float64(k string, f float64) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddFloat64(k, f)
	return e
}

func (e *Entry) Float64s(k string, ff []float64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, float64s(ff)))
	return e
}

func (e *Entry) Complex64(k string, c complex64) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddComplex128(k, complex128(c))
	return e
}

func (e *Entry) Complex64s(k string, cc []complex64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, complex64s(cc)))
	return e
}

func (e *Entry) Complex128(k string, c complex128) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddComplex128(k, c)
	return e
}

func (e *Entry) Complex128s(k string, cc []complex128) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, complex128s(cc)))
	return e
}

// Nanosecond encode time.Duration to an string nanoseconds; format sample "1004854348ns".
func (e *Entry) Nanosecond(k string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddDuration(k, d, DurationFormatNano)
	//e.Encoder.AddDuration(k, durationNano(val))
	return e
}

// Microsecond encode time.Duration to an string microseconds, format sample "1004854us".
func (e *Entry) Microsecond(k string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddDuration(k, d, DurationFormatMicro)
	//e.Encoder.AddDuration(k, durationMicro(val))
	return e
}

// Millisecond encode time.Duration to an string milliseconds, format sample "1004ms".
func (e *Entry) Millisecond(k string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddDuration(k, d, DurationFormatMilli)
	//e.Encoder.AddDuration(k, durationMilli(val))
	return e
}

// Second encode time.Duration to an string seconds, format sample "1.004854348s".
func (e *Entry) Second(k string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddDuration(k, d, DurationFormatSecond)
	//e.Encoder.AddDuration(k, durationSecond(val))
	return e
}

// Minute encode time.Duration to an string minutes, format sample "10min".
func (e *Entry) Minute(k string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddDuration(k, d, DurationFormatMinute)
	//e.Encoder.AddDuration(k, durationMinute(val))
	return e
}

// Hour encode time.Duration to an string hours, format sample "2h".
func (e *Entry) Hour(k string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddDuration(k, d, DurationFormatHour)
	//e.Encoder.AddDuration(k, durationHour(val))
	return e
}

func (e *Entry) Error(k string, err error) *Entry {
	if e == nil {
		return nil
	}
	if err != nil {
		e.Encoder.AddString(k, err.Error())
	} else {
		e.Encoder.AddString(k, "<nil>")
	}
	return e
}

func (e *Entry) Errors(k string, errs []error) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, errorArray(errs)))
	return e
}

func (e *Entry) Time(k string, t time.Time, layout string) *Entry {
	if e == nil {
		return nil
	}
	e.Encoder.AddTime(k, t, layout)
	return e
}

func (e *Entry) Array(k string, am ArrayMarshaler) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddArray(k, am))
	return e
}

func (e *Entry) Object(k string, om ObjectMarshaler) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.Encoder.AddObject(k, om))
	return e
}

// Any uses reflection to serialize arbitrary objects, so it can be
// slow and allocation-heavy.
func (e *Entry) Any(k string, i interface{}) *Entry {
	if e == nil {
		return nil
	}
	switch m := i.(type) {
	case ArrayMarshaler:
		e.withError(e.Encoder.AddArray(k, m))
	case ObjectMarshaler:
		e.withError(e.Encoder.AddObject(k, m))
	default:
		e.withError(e.Encoder.AddInterface(k, i))
	}
	return e
}
