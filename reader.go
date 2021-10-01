package paasio

import (
	"io"
	"sync"
	"time"
)

type nopReader struct {
	error
}

func (r nopReader) Read(p []byte) (int, error) {
	time.Sleep(sleepTime)
	if r.error != nil {
		return 0, r.error
	}
	return len(p), nil
}

func NewReadCounter(r io.Reader) ReadCounter {
	rw := readerWriter{r, nopWriter{}}
	return &readerWriterCounter{rw: rw, Mutex: new(sync.Mutex)}
}
