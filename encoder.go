package glog

import "time"

// EncoderFunc used to return a new Encoder instances
type EncoderFunc func() Encoder

// FieldEncoder used to add single elements
type FieldEncoder interface {
	AppendByte(b byte) // AppendByte encode b to the form of binary
	AppendString(s string)
	AppendBool(v bool)
	AppendInt64(i int64)
	AppendUnt64(i uint64)
	AppendFloat64(f float64)
	AppendComplex128(c complex128)

	AppendTime(t time.Time, layout string)
	AppendDuration(d time.Duration, layout int8)

	AppendArray(am ArrayMarshaler) error
	AppendObject(om ObjectMarshaler) error

	// AppendInterface uses reflection to serialize arbitrary objects, so it's
	// slow and allocation-heavy.
	AppendInterface(i interface{}) error
}

// ArrayEncoder used to add array-type field
type ArrayEncoder interface {
	FieldEncoder
}

// ObjectEncoder used to add an k/v field
type ObjectEncoder interface {
	AddByte(key string, b byte) // AddByte encode b to the form of binary
	AddString(key string, s string)
	AddBool(key string, v bool)
	AddInt64(key string, i int64)
	AddUnt64(key string, i uint64)
	AddFloat64(key string, f float64)
	AddComplex128(key string, c complex128)

	AddTime(key string, t time.Time, layout string)
	AddDuration(key string, d time.Duration, layout int8)

	AddArray(key string, am ArrayMarshaler) error
	AddObject(key string, om ObjectMarshaler) error

	// AddInterface uses reflection to serialize arbitrary objects, so it can be
	// slow and allocation-heavy.
	AddInterface(key string, i interface{}) error
}

// BuildEncoder used to add some specific fields
type BuildEncoder interface {
	AddMsg(msg string)
	AddEntryTime(t time.Time, layout string)
	AddLevel(level Level)
	AddCaller(skip int)

	// AddBeginMarker add the begin marker
	AddBeginMarker()
	// AppendEndMarker add the end marker
	AddEndMarker()
	// AppendLineBreak add the line break
	AddLineBreak()

	// WriteIn used to write encoded data
	WriteIn(p []byte) error
}

type Encoder interface {
	ObjectEncoder
	BuildEncoder

	// Bytes returns a mutable reference to the byte slice of the Encoder
	Bytes() []byte
	// Callers must not retain references to the Encoder after calling Close.
	Close() error
}
