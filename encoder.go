package glog

import "time"

// EncoderFunc used to get a new Encoder instances
type EncoderFunc func() Encoder

// FieldEncoder used to add single field
type FieldEncoder interface {
	AppendByte(val byte) // AppendByte encode val to the form of binary
	AppendString(val string)
	AppendBool(val bool)
	AppendInt64(val int64)
	AppendUnt64(val uint64)
	AppendFloat64(val float64)
	AppendTime(val time.Time, layout string)
	AppendComplex128(val complex128)

	AppendArray(arr ArrayMarshaler) error
	AppendObject(obj ObjectMarshaler) error

	// AppendInterface uses reflection to serialize arbitrary objects, so it's
	// slow and allocation-heavy.
	AppendInterface(val interface{}) error
}

// ArrayEncoder used to add array-type field
type ArrayEncoder interface {
	FieldEncoder
}

// ObjectEncoder used to add a complete k/v field
type ObjectEncoder interface {
	AddByte(key string, val byte) // AddByte encode val to the form of binary
	AddString(key string, val string)
	AddBool(key string, val bool)
	AddInt64(key string, val int64)
	AddUnt64(key string, val uint64)
	AddFloat64(key string, val float64)
	AddTime(key string, val time.Time, layout string)
	AddComplex128(key string, val complex128)

	AddArray(key string, arr ArrayMarshaler) error
	AddObject(key string, obj ObjectMarshaler) error

	// AddInterface uses reflection to serialize arbitrary objects, so it can be
	// slow and allocation-heavy.
	AddInterface(key string, val interface{}) error
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
	ArrayEncoder
	BuildEncoder

	// Bytes returns a mutable reference to the byte slice of the Encoder
	Bytes() []byte
	// Callers must not retain references to the Encoder after calling Close.
	Close() error
}
