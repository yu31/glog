package glog

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Logger define the logger.
type Logger struct {
	// level set the minimum accepted level,
	// less than the level will be ignored;
	level Level

	// timeLayout set the time format in log message.
	timeLayout string

	// caller set whether adds caller info in log message
	caller bool

	// encoderFunc used to get a new encoder in log entry
	// Notes: change the encoderFunc will cause the fields empty and rebuild
	encoderFunc EncoderFunc

	// fields add fixed field into every log entry
	fields Encoder

	// exporter used to export the log by every entry.Fire
	exporter Exporter

	// errorOutput is the error output writer of this logger
	// logger will write error message into this while failed to log message
	//
	// NOTE: logger will not check returning error of this writer
	errorOutput io.Writer
}

func NewDefault() *Logger {
	l := &Logger{
		level:       DebugLevel,
		timeLayout:  defaultTimeLayout,
		caller:      false,
		encoderFunc: TextEncoder,
		fields:      nil,
		exporter:    DefaultExporter,
		errorOutput: os.Stderr,
	}
	l.fields = l.encoderFunc()
	return l
}

// WithLevel will set logger's level.
func (l *Logger) WithLevel(level Level) *Logger {
	l.level = level
	return l
}

// WithTimeLayout will set logger's timeLayout.
func (l *Logger) WithTimeLayout(layout string) *Logger {
	l.timeLayout = layout
	return l
}

// WithCaller will set logger's caller.
func (l *Logger) WithCaller(ok bool) *Logger {
	l.caller = ok
	return l
}

// WithExporter will set logger's exporter.
func (l *Logger) WithExporter(exporter Exporter) *Logger {
	l.exporter = exporter
	return l
}

// WithEncoderFunc will set logger's encoderFunc.
func (l *Logger) WithEncoderFunc(f EncoderFunc) *Logger {
	l.encoderFunc = f
	if l.fields != nil {
		_ = l.fields.Close()
	}
	l.fields = f()
	return l
}

// WithErrorOutput will set logger's exporter.
func (l *Logger) WithErrorOutput(w io.Writer) *Logger {
	l.errorOutput = w
	return l
}

// WithFields for add fixed fields into the log entry
func (l *Logger) WithFields() Encoder {
	return l.fields
}

// ResetFields for clear the data in fields
func (l *Logger) ResetFields() Encoder {
	_ = l.fields.Close()
	l.fields = l.encoderFunc()
	return l.fields
}

// Clone will copy and return a new logger
func (l *Logger) Clone() *Logger {
	lc := &Logger{
		level:       l.level,
		timeLayout:  l.timeLayout,
		caller:      l.caller,
		exporter:    l.exporter,
		encoderFunc: l.encoderFunc,
		fields:      l.encoderFunc(),
		errorOutput: l.errorOutput,
	}
	err := lc.fields.WriteIn(l.fields.Bytes())
	if err != nil {
		_, _ = fmt.Fprintf(l.errorOutput, "%s [inner] write fields fail when clone: %v\n", time.Now().Format(l.timeLayout), err)
	}
	return lc
}

// Close close the logger for releasing resources
//
// Notes: Close don't close the Exporter because it may be
// shared by multiple Logger instances
func (l *Logger) Close() error {
	var errs []error

	if err := l.fields.Close(); err != nil {
		errs = append(errs, err)
	}

	l.exporter = nil
	l.encoderFunc = nil
	l.fields = nil
	l.errorOutput = nil

	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%v", errs)
}

func (l *Logger) newEntry(level Level) *Entry {
	if level >= l.level {
		return newEntry(l, level)
	}
	return nil
}

// Debug will return an Entry with DebugLevel.
func (l *Logger) Debug() *Entry {
	return l.newEntry(DebugLevel)
}

// Debug will return an Entry with InfoLevel.
func (l *Logger) Info() *Entry {
	return l.newEntry(InfoLevel)
}

// Debug will return an Entry with WarnLevel.
func (l *Logger) Warn() *Entry {
	return l.newEntry(WarnLevel)
}

// Debug will return an Entry with ErrorLevel.
func (l *Logger) Error() *Entry {
	return l.newEntry(ErrorLevel)
}

// Debug will return an Entry with FatalLevel.
func (l *Logger) Fatal() *Entry {
	return l.newEntry(FatalLevel)
}
