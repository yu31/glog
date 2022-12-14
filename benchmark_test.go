package glog

import (
	"io/ioutil"
	"testing"
	"time"
)

var (
	fakeMessage = "Test logging, but use a somewhat realistic message length."
)

func BenchmarkNewDefault(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l := NewDefault()
			_ = l.Close()
		}
	})
}

func BenchmarkLogger_Clone(b *testing.B) {
	l := NewDefault()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			nl := l.Clone()
			_ = nl.Close()
		}
	})
}

func BenchmarkLogEmpty(b *testing.B) {
	l := NewDefault().WithExporter(StandardExporter(ioutil.Discard))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Debug().Msg("").Fire()
		}
	})
}

func BenchmarkLogDisabled(b *testing.B) {
	l := NewDefault().WithExporter(StandardExporter(ioutil.Discard)).WithLevel(InfoLevel)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Debug().Msg(fakeMessage).Fire()
		}
	})
}

func BenchmarkLogMsg(b *testing.B) {
	l := NewDefault().WithExporter(StandardExporter(ioutil.Discard))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Debug().Msg(fakeMessage).Fire()
		}
	})
}

func BenchmarkLogFields(b *testing.B) {
	l := NewDefault().WithExporter(StandardExporter(ioutil.Discard))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info().
				String("string", "four!").
				Time("time", time.Time{}, "").
				Int("int", 123).
				Float32("float", -2.203230293249593).
				Msg(fakeMessage).
				Fire()
		}
	})
}

func BenchmarkLogWithFields(b *testing.B) {
	l := NewDefault().WithExporter(StandardExporter(ioutil.Discard))
	l.WithFields().AddString("string", "four")
	l.WithFields().AddTime("time", time.Time{}, "")
	l.WithFields().AddInt64("int", 123)
	l.WithFields().AddFloat64("float", -2.203230293249593)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Debug().Msg(fakeMessage).Fire()
		}
	})
}

func BenchmarkLog10Fields(b *testing.B) {
	l := NewDefault().WithExporter(StandardExporter(ioutil.Discard))
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Debug().
				Msg("Test logger out").
				String("String1", "Value1").
				String("String2", "Value2").
				String("String2", "Value3").
				Byte("Byte", 'a').
				Bytes("Bytes", []byte("abc")).
				Int64("Int64", 64).
				Uint64("Uint64", 64).
				Float64("Float64", 99.99).
				Bool("Bool", true).
				Fire()
		}
	})
}

func BenchmarkLog10String(b *testing.B) {
	l := NewDefault().WithExporter(StandardExporter(ioutil.Discard))
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Debug().
				String("String", "HelloWorld").
				String("String", "HelloWorld").
				String("String", "HelloWorld").
				String("String", "HelloWorld").
				String("String", "HelloWorld").
				String("String", "HelloWorld").
				String("String", "HelloWorld").
				String("String", "HelloWorld").
				String("String", "HelloWorld").
				String("String", "HelloWorld").
				Fire()
		}
	})
}

func BenchmarkLog10Int64(b *testing.B) {
	l := NewDefault().WithExporter(StandardExporter(ioutil.Discard))
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Debug().
				Int("int64", 1234567890).
				Int("int64", 1234567890).
				Int("int64", 1234567890).
				Int("int64", 1234567890).
				Int("int64", 1234567890).
				Int("int64", 1234567890).
				Int("int64", 1234567890).
				Int("int64", 1234567890).
				Int("int64", 1234567890).
				Int("int64", 1234567890).
				Fire()
		}
	})
}
