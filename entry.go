package glog

import (
	"fmt"
	"time"
)

// Entry used to build a log record.
type Entry struct {
	level   Level
	encoder Encoder

	l *Logger
}

// newEntry will create a new entry with level and fields.
func newEntry(l *Logger, level Level) *Entry {
	e := &Entry{
		level:   level,
		encoder: l.encoderFunc(),

		l: l,
	}
	e.encodeHeads()
	return e
}

// withError handle any error if happen in entry inside
func (e *Entry) withError(err error) {
	if err == nil {
		return
	}
	_, _ = fmt.Fprintf(e.l.errorOutput, "[glog] %s handle log entry error: %v\n", time.Now().Format(e.l.timeLayout), err)
}

func (e *Entry) encodeHeads() {
	e.encoder.AddBeginMarker()
	e.encoder.AddEntryTime(time.Now(), e.l.timeLayout)
	e.encoder.AddLevel(e.level)
}

func (e *Entry) encodeEnds() {
	e.withError(e.encoder.WriteIn(e.l.fields.Bytes()))
	if e.l.caller {
		e.encoder.AddCaller(2)
	}
	e.encoder.AddEndMarker()
	e.encoder.AddLineBreak()
}

func (e *Entry) free() {
	if e == nil {
		return
	}
	e.withError(e.encoder.Close())
	e.l = nil
	e.encoder = nil
}

// Fire sends the *Entry to Logger's exporter.
//
// NOTICE: once this method is called, the *Entry should be disposed.
// Calling Fire twice can have unexpected result.
func (e *Entry) Fire() {
	if e == nil {
		return
	}
	e.encodeEnds()

	// NOTICE: The `data` will be reuse by put back to sync.Pool.
	// Thus the `*Record` should be disposed after the `Export` returns.
	e.withError(e.l.exporter.Export(&Record{
		ctx:   e.l.ctx,
		level: e.level,
		data:  e.encoder.Bytes(),
	}))

	// Release resources
	e.free()
}

func (e *Entry) Msg(msg string) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddMsg(msg)
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
	e.encoder.AddRawBytes(k, bs)
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
	e.encoder.AddRawString(k, s)
	return e
}

// Byte encode the value to an integer format.
func (e *Entry) Byte(k string, b byte) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddByte(k, b)
	return e
}

// Bytes encode the value to an integer array.
func (e *Entry) Bytes(k string, bb []byte) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, byteArray(bb)))
	return e
}

func (e *Entry) String(k string, s string) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddString(k, s)
	return e
}

func (e *Entry) Strings(k string, ss []string) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, stringArray(ss)))
	return e
}

func (e *Entry) Bool(k string, v bool) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddBool(k, v)
	return e
}

func (e *Entry) Bools(k string, vv []bool) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, bools(vv)))
	return e
}

func (e *Entry) Int(k string, i int) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddInt64(k, int64(i))
	return e
}

func (e *Entry) Ints(k string, ii []int) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, ints(ii)))
	return e
}

func (e *Entry) Int8(k string, i int8) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddInt64(k, int64(i))
	return e
}

func (e *Entry) Int8s(k string, ii []int8) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, int8s(ii)))
	return e
}

func (e *Entry) Int16(k string, i int16) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddInt64(k, int64(i))
	return e
}

func (e *Entry) Int16s(k string, ii []int16) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, int16s(ii)))
	return e
}

func (e *Entry) Int32(k string, i int32) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddInt64(k, int64(i))
	return e
}

func (e *Entry) Int32s(k string, ii []int32) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, int32s(ii)))
	return e
}

func (e *Entry) Int64(k string, i int64) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddInt64(k, int64(i))
	return e
}

func (e *Entry) Int64s(k string, ii []int64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, int64s(ii)))
	return e
}

func (e *Entry) Uint(k string, i uint) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddUnt64(k, uint64(i))
	return e
}

func (e *Entry) Uints(k string, ii []uint) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, uints(ii)))
	return e
}

func (e *Entry) Uint8(k string, i uint8) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddUnt64(k, uint64(i))
	return e
}

func (e *Entry) Uint8s(k string, ii []uint8) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, uint8s(ii)))
	return e
}

func (e *Entry) Uint16(k string, i uint16) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddUnt64(k, uint64(i))
	return e
}

func (e *Entry) Uint16s(k string, ii []uint16) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, uint16s(ii)))
	return e
}

func (e *Entry) Uint32(k string, i uint32) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddUnt64(k, uint64(i))
	return e
}

func (e *Entry) Uint32s(k string, ii []uint32) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, uint32s(ii)))
	return e
}

func (e *Entry) Uint64(k string, i uint64) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddUnt64(k, i)
	return e
}

func (e *Entry) Uint64s(k string, ii []uint64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, uint64s(ii)))
	return e
}

func (e *Entry) Float32(k string, f float32) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddFloat64(k, float64(f))
	return e
}

func (e *Entry) Float32s(k string, ff []float32) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, float32s(ff)))
	return e
}

func (e *Entry) Float64(k string, f float64) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddFloat64(k, f)
	return e
}

func (e *Entry) Float64s(k string, ff []float64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, float64s(ff)))
	return e
}

func (e *Entry) Complex64(k string, c complex64) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddComplex128(k, complex128(c))
	return e
}

func (e *Entry) Complex64s(k string, cc []complex64) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, complex64s(cc)))
	return e
}

func (e *Entry) Complex128(k string, c complex128) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddComplex128(k, c)
	return e
}

func (e *Entry) Complex128s(k string, cc []complex128) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, complex128s(cc)))
	return e
}

// Nanosecond encode time.Duration to an string nanoseconds; format sample "1004854348ns".
func (e *Entry) Nanosecond(k string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddDuration(k, d, DurationFormatNano)
	//e.Encoder.AddDuration(k, durationNano(val))
	return e
}

// Microsecond encode time.Duration to an string microseconds, format sample "1004854us".
func (e *Entry) Microsecond(k string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddDuration(k, d, DurationFormatMicro)
	//e.Encoder.AddDuration(k, durationMicro(val))
	return e
}

// Millisecond encode time.Duration to an string milliseconds, format sample "1004ms".
func (e *Entry) Millisecond(k string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddDuration(k, d, DurationFormatMilli)
	//e.Encoder.AddDuration(k, durationMilli(val))
	return e
}

// Second encode time.Duration to an string seconds, format sample "1.004854348s".
func (e *Entry) Second(k string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddDuration(k, d, DurationFormatSecond)
	//e.Encoder.AddDuration(k, durationSecond(val))
	return e
}

// Minute encode time.Duration to an string minutes, format sample "10min".
func (e *Entry) Minute(k string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddDuration(k, d, DurationFormatMinute)
	//e.Encoder.AddDuration(k, durationMinute(val))
	return e
}

// Hour encode time.Duration to an string hours, format sample "2h".
func (e *Entry) Hour(k string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddDuration(k, d, DurationFormatHour)
	//e.Encoder.AddDuration(k, durationHour(val))
	return e
}

func (e *Entry) Error(k string, err error) *Entry {
	if e == nil {
		return nil
	}
	if err != nil {
		e.encoder.AddString(k, err.Error())
	} else {
		e.encoder.AddString(k, "<nil>")
	}
	return e
}

func (e *Entry) Errors(k string, errs []error) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, errorArray(errs)))
	return e
}

func (e *Entry) Time(k string, t time.Time, layout string) *Entry {
	if e == nil {
		return nil
	}
	e.encoder.AddTime(k, t, layout)
	return e
}

func (e *Entry) Array(k string, am ArrayMarshaler) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddArray(k, am))
	return e
}

func (e *Entry) Object(k string, om ObjectMarshaler) *Entry {
	if e == nil {
		return nil
	}
	e.withError(e.encoder.AddObject(k, om))
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
		e.withError(e.encoder.AddArray(k, m))
	case ObjectMarshaler:
		e.withError(e.encoder.AddObject(k, m))
	default:
		e.withError(e.encoder.AddInterface(k, i))
	}
	return e
}
