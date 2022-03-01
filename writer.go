package glog

import (
	"io"
)

type nopWriterCloser struct {
	io.Writer
}

func (nopWriterCloser) Close() error { return nil }

// NopWriterCloser returns a ReadCloser with a no-op Close method wrapping
// the provided Reader r. It same as ioutil.NopCloser.
func NopWriterCloser(r io.Writer) io.WriteCloser {
	return nopWriterCloser{r}
}
