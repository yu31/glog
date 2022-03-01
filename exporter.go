package glog

import (
	"fmt"
	"io"
	"os"
)

var (
	DefaultExporter = StandardExporter(NopWriterCloser(os.Stdout))
)

// Exporter used to handle the Entry.
type Exporter interface {
	// Export NOTICE: The `data` will be reuse by put back to sync.Pool.
	//
	// Thus the `data` should be disposed after the `Export` returns.
	// If the `data` is processed asynchronously, you should get data with Copy method.
	Export(record *Record) error

	// Close to close the exporter.
	Close() error
}

var (
	_ Exporter = (*standardExporter)(nil)
	_ Exporter = (*matcherExporter)(nil)
	_ Exporter = (*multipleExporter)(nil)
)

// StandardExporter return a Exporter implements by standardExporter.
// This used to write log record to writer.
func StandardExporter(w io.Writer) Exporter {
	return &standardExporter{w: w}
}

// standardExporter creates an exporter that write log entry into an io.Writer.
type standardExporter struct {
	w io.Writer
}

func (exp *standardExporter) Export(record *Record) error {
	_, err := exp.w.Write(record.Bytes())
	return err
}

// Close for close the Exporter.
func (exp *standardExporter) Close() error {
	if c, ok := exp.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// FilterExporter return a Exporter implements by matcherExporter;
// This used to write only the specified level of Entry.
func FilterExporter(w io.Writer, f Filter) Exporter {
	return &matcherExporter{w: w, f: f}
}

// matcherExporter creates an exporter that write log entry into an io.Writer.
type matcherExporter struct {
	w io.Writer
	f Filter
}

func (exp *matcherExporter) Export(record *Record) error {
	if !exp.f.Match(record.Level()) {
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

//// FileExporter is a FilterExporter wrapper that uses a file.
//func FileExporter(name string, f Filter) (Exporter, error) {
//	w, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
//	if err != nil {
//		return nil, err
//	}
//	return FilterExporter(w, f), nil
//}
