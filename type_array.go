package glog

type bools []bool

func (vv bools) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendBool(vv[i])
	}
	return nil
}

type complex128s []complex128

func (vv complex128s) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendComplex128(vv[i])
	}
	return nil
}

type complex64s []complex64

func (vv complex64s) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendComplex128(complex128(vv[i]))
	}
	return nil
}

type float64s []float64

func (vv float64s) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendFloat64(vv[i])
	}
	return nil
}

type float32s []float32

func (vv float32s) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendFloat64(float64(vv[i]))
	}
	return nil
}

type ints []int

func (vv ints) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendInt64(int64(vv[i]))
	}
	return nil
}

type int64s []int64

func (vv int64s) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendInt64(vv[i])
	}
	return nil
}

type int32s []int32

func (vv int32s) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendInt64(int64(vv[i]))
	}
	return nil
}

type int16s []int16

func (vv int16s) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendInt64(int64(vv[i]))
	}
	return nil
}

type int8s []int8

func (vv int8s) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendInt64(int64(vv[i]))
	}
	return nil
}

type uints []uint

func (vv uints) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendUnt64(uint64(vv[i]))
	}
	return nil
}

type uint64s []uint64

func (vv uint64s) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendUnt64(vv[i])
	}
	return nil
}

type uint32s []uint32

func (vv uint32s) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendUnt64(uint64(vv[i]))
	}
	return nil
}

type uint16s []uint16

func (vv uint16s) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendUnt64(uint64(vv[i]))
	}
	return nil
}

type uint8s []uint8

func (vv uint8s) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendUnt64(uint64(vv[i]))
	}
	return nil
}

type byteArray []byte

func (vv byteArray) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendByte(vv[i])
	}
	return nil
}

type stringArray []string

func (vv stringArray) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendString(vv[i])
	}
	return nil
}

type errorArray []error

func (vv errorArray) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendString(vv[i].Error())
	}
	return nil
}
