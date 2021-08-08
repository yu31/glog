package glog

import (
	"fmt"
	"io"
	"os"
)

var (
	_ Exporter = (*matcherExporter)(nil)
	_ Exporter = (*multipleExporter)(nil)
)

var (
	DefaultExporter = MatchExporter(os.Stdout, nil)
)

// Exporter used to handle the Entry.
type Exporter interface {
	Export(entry *Entry) error
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

func (exp *matcherExporter) Export(entry *Entry) error {
	if exp.f != nil && !exp.f.Match(entry.Level) {
		return nil
	}
	_, err := exp.w.Write(entry.Encoder.Bytes())
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

func (exp *multipleExporter) Export(entry *Entry) error {
	var errs []error
	for i := range exp.exporters {
		err := exp.exporters[i].Export(entry)
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
