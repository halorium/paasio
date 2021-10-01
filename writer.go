package paasio

import (
	"io"
	"sync"
	"time"
)

type nopWriter struct {
	error
}

func (w nopWriter) Write(p []byte) (int, error) {
	time.Sleep(sleepTime)
	if w.error != nil {
		return 0, w.error
	}
	return len(p), nil
}

func NewWriteCounter(w io.Writer) WriteCounter {
	rw := readerWriter{nopReader{}, w}
	return &readerWriterCounter{rw: rw, Mutex: new(sync.Mutex)}
}
