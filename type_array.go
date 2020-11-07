package glog

type bools []bool

func (eles bools) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendBool(eles[i])
	}
	return nil
}

type complex128s []complex128

func (eles complex128s) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendComplex128(eles[i])
	}
	return nil
}

type complex64s []complex64

func (eles complex64s) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendComplex128(complex128(eles[i]))
	}
	return nil
}

type float64s []float64

func (eles float64s) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendFloat64(eles[i])
	}
	return nil
}

type float32s []float32

func (eles float32s) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendFloat64(float64(eles[i]))
	}
	return nil
}

type ints []int

func (eles ints) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendInt64(int64(eles[i]))
	}
	return nil
}

type int64s []int64

func (eles int64s) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendInt64(eles[i])
	}
	return nil
}

type int32s []int32

func (eles int32s) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendInt64(int64(eles[i]))
	}
	return nil
}

type int16s []int16

func (eles int16s) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendInt64(int64(eles[i]))
	}
	return nil
}

type int8s []int8

func (eles int8s) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendInt64(int64(eles[i]))
	}
	return nil
}

type uints []uint

func (eles uints) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendUnt64(uint64(eles[i]))
	}
	return nil
}

type uint64s []uint64

func (eles uint64s) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendUnt64(eles[i])
	}
	return nil
}

type uint32s []uint32

func (eles uint32s) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendUnt64(uint64(eles[i]))
	}
	return nil
}

type uint16s []uint16

func (eles uint16s) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendUnt64(uint64(eles[i]))
	}
	return nil
}

type uint8s []uint8

func (eles uint8s) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendUnt64(uint64(eles[i]))
	}
	return nil
}

type runes []rune

func (eles runes) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendInt64(int64(eles[i]))
	}
	return nil
}

type stringArray []string

func (eles stringArray) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendString(eles[i])
	}
	return nil
}

type byteArray []byte

func (eles byteArray) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendByte(eles[i])
	}
	return nil
}

type errorArray []error

func (eles errorArray) MarshalArray(arr ArrayEncoder) error {
	for i := range eles {
		arr.AppendString(eles[i].Error())
	}
	return nil
}
