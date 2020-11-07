package glog

// ArrayMarshaler allows user-defined data types to efficiently add an array into to the log entry.
type ArrayMarshaler interface {
	MarshalArray(ae ArrayEncoder) error
}

// ArrayMarshalerFunc is a type adapter that turns a function into an ArrayMarshaler.
type ArrayMarshalerFunc func(ae ArrayEncoder) error

// MarshalLogArray calls the underlying function.
func (f ArrayMarshalerFunc) MarshalArray(ae ArrayEncoder) error {
	return f(ae)
}

// ObjectMarshaler allows user-defined data types to efficiently add an object into to the log entry.
type ObjectMarshaler interface {
	MarshalObject(oe ObjectEncoder) error
}

// ObjectMarshalerFunc is a type adapter that turns a function into an ObjectMarshaler.
type ObjectMarshalerFunc func(oe ObjectEncoder) error

// MarshalLogObject calls the underlying function.
func (f ObjectMarshalerFunc) MarshalObject(oe ObjectEncoder) error {
	return f(oe)
}
