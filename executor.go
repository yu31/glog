package glog

import (
	"fmt"
	"io"
	"os"
)

var (
	DefaultExecutor = MatchExecutor(os.Stdout, nil)
)

// Executor used to handle the Entry
type Executor interface {
	Execute(entry *Entry) error
	Close() error
}

// MatchExecutor return a Executor implements by matcherExecutor;
// This used to write only the specified level of Entry.
// The parameters Filter allowed to be nil
func MatchExecutor(w io.Writer, f Filter) Executor {
	return &matcherExecutor{w: w, f: f}
}

// matcherExecutor creates an executor that write log entry into an io.Writer.
type matcherExecutor struct {
	w io.Writer
	f Filter
}

func (exe *matcherExecutor) Execute(entry *Entry) error {
	if exe.f != nil && !exe.f.Match(entry.Level) {
		return nil
	}
	_, err := exe.w.Write(entry.Encoder.Bytes())
	return err
}

// Close for close the Executor
func (exe *matcherExecutor) Close() error {
	if c, ok := exe.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// MultipleExecutor used for apply multiple Executor to a Entry
func MultipleExecutor(executors ...Executor) Executor {
	return &multipleExecutor{executors: executors}
}

type multipleExecutor struct {
	executors []Executor
}

func (exe *multipleExecutor) Execute(entry *Entry) error {
	var errs []error
	for i := range exe.executors {
		err := exe.executors[i].Execute(entry)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%v", errs)
}

func (exe *multipleExecutor) Close() error {
	var errs []error

	for i := range exe.executors {
		if e := exe.executors[i].Close(); e != nil {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%v", errs)
}

// FileExecutor is a MatchExecutor wrapper that uses a file
func FileExecutor(name string, f Filter) (Executor, error) {
	w, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return MatchExecutor(w, f), nil
}
