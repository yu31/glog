package glog

import (
	"context"
	"fmt"
	"io"
	"os"
)

var (
	DefaultExporter = MatchExporter(os.Stdout, nil)
)

var (
	_ Exporter = (*matcherExporter)(nil)
	_ Exporter = (*multipleExporter)(nil)
)

// Record represents the Entry's content.
type Record struct {
	ctx   context.Context
	level Level
	data  []byte
}

// Context returns context where in Logger.
func (r *Record) Context() context.Context {
	return r.ctx
}

// Level returns the log level of the entry.
func (r *Record) Level() Level {
	return r.level
}

// Bytes returns the Entry's content.
func (r *Record) Bytes() []byte {
	return r.data
}

// Copy returns an copy of Entry's content.
func (r *Record) Copy() []byte {
	bs := make([]byte, len(r.data))
	copy(bs, r.data)
	return bs
}

// Exporter used to handle the Entry.
type Exporter interface {
	// NOTICE: The `data` will be reuse by put back to sync.Pool.
	//
	// Thus the `data` should be disposed after the `Export` returns.
	// If the `data` is processed asynchronously, you should get data with Copy method.
	Export(record *Record) error

	// Close to close the exporter.
	Close() error
}

// MatchExporter return a Exporter implements by matcherExporter;
// This used to write only the specified level of Entry.
// The parameters Filter allowed to be nil.
func MatchExporter(w io.Writer, f Filter) Exporter {
	return &matcherExporter{w: w, f: f}
}

// matcherExporter creates an exporter that write log entry into an io.Writer.
type matcherExporter struct {
	w io.Writer
	f Filter
}

func (exp *matcherExporter) Export(record *Record) error {
	if exp.f != nil && !exp.f.Match(record.Level()) {
		return nil
	}
	_, err := exp.w.Write(record.Bytes())
	return err
}

// Close for close the Exporter.
func (exp *matcherExporter) Close() error {
	if c, ok := exp.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// MultipleExporter used for apply multiple Exporter to a Entry.
func MultipleExporter(exporters ...Exporter) Exporter {
	return &multipleExporter{exporters: exporters}
}

type multipleExporter struct {
	exporters []Exporter
}

func (exp *multipleExporter) Export(record *Record) error {
	var errs []error
	for i := range exp.exporters {
		err := exp.exporters[i].Export(record)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%v", errs)
}

func (exp *multipleExporter) Close() error {
	var errs []error

	for i := range exp.exporters {
		if e := exp.exporters[i].Close(); e != nil {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%v", errs)
}

// FileExporter is a MatchExporter wrapper that uses a file.
func FileExporter(name string, f Filter) (Exporter, error) {
	w, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return MatchExporter(w, f), nil
}
