package glog

import (
	"time"

	"github.com/DataWorkbench/glog/pkg/buffer"
)

const (
	defaultTimeLayout = time.RFC3339Nano
)

// Defines the time format type.
const (
	// TimeFormatUnixSecond defines a time format that makes time fields to be
	// serialized as Unix timestamp integers in seconds.
	TimeFormatUnixSecond = "UnixSecond"

	// TimeFormatUnixMs defines a time format that makes time fields to be
	// serialized as Unix timestamp integers in milliseconds.
	TimeFormatUnixMilli = "UnixMilli"

	// TimeFormatUnixMicro defines a time format that makes time fields to be
	// serialized as Unix timestamp integers in microseconds.
	TimeFormatUnixMicro = "UnixMicro"

	// TimeFormatUnixNs defines a time format that makes time fields to be
	// serialized as Unix timestamp integers in nanoseconds.
	TimeFormatUnixNano = "UnixNano"
)

// Defines the format type for time.Duration.
const (
	DurationFormatNano int8 = iota
	DurationFormatMicro
	DurationFormatMilli
	DurationFormatSecond
	DurationFormatMinute
	DurationFormatHour
)

const (
	nanoSuffix   = "ns"
	microSuffix  = "µs"
	milliSuffix  = "ms"
	secondSuffix = "s"
	minuteSuffix = "min"
	hourSuffix   = "h"
)

// AppendDuration encode the time.Duration to a strings by specified layout.
func AppendDuration(buf *buffer.Buffer, d time.Duration, layout int8) {
	switch layout {
	case DurationFormatNano:
		buf.AppendInt(int64(d))
		buf.AppendString(nanoSuffix)
	case DurationFormatMicro:
		sec := d / time.Microsecond
		nsec := d % time.Microsecond
		buf.AppendFloat(float64(sec)+float64(nsec)/1e3, 64)
		buf.AppendString(microSuffix)
	case DurationFormatMilli:
		sec := d / time.Millisecond
		nsec := d % time.Millisecond
		buf.AppendFloat(float64(sec)+float64(nsec)/1e6, 64)
		buf.AppendString(milliSuffix)
	case DurationFormatSecond:
		buf.AppendFloat(d.Seconds(), 64)
		buf.AppendString(secondSuffix)
	case DurationFormatMinute:
		buf.AppendFloat(d.Minutes(), 64)
		buf.AppendString(minuteSuffix)
	case DurationFormatHour:
		buf.AppendFloat(d.Hours(), 64)
		buf.AppendString(hourSuffix)
	default:
		buf.AppendString(d.String())
	}
}
