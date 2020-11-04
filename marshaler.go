package glog

// ArrayMarshaler allows user-defined data types to efficiently add themselves into to the log entry.
type ArrayMarshaler interface {
	MarshalLogArray(arr ArrayEncoder) error
}

// ArrayMarshalerFunc is a type adapter that turns a function into an ArrayMarshaler.
type ArrayMarshalerFunc func(arr ArrayEncoder) error

// MarshalLogArray calls the underlying function.
func (f ArrayMarshalerFunc) MarshalLogArray(arr ArrayEncoder) error {
	return f(arr)
}

// ObjectMarshaler allows user-defined data types to efficiently add themselves into to the log entry.
type ObjectMarshaler interface {
	MarshalLogObject(ObjectEncoder) error
}

// ObjectMarshalerFunc is a type adapter that turns a function into an ObjectMarshaler.
type ObjectMarshalerFunc func(ObjectEncoder) error

// MarshalLogObject calls the underlying function.
func (f ObjectMarshalerFunc) MarshalLogObject(enc ObjectEncoder) error {
	return f(enc)
}
