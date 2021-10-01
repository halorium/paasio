package paasio

import (
	"io"
	"sync"
)

type nopReaderWriter struct {
	nopReader
	nopWriter
}

type readWriter readerWriter

type readerWriter struct {
	io.Reader
	io.Writer
}

type readerWriterCounter struct {
	rw         ReaderWriter
	TotalBytes int64
	TotalOps   int
	Mutex      *sync.Mutex
}

func (o *readerWriterCounter) Read(p []byte) (n int, err error) {
	o.Mutex.Lock()
	defer o.Mutex.Unlock()

	n, err = o.rw.Read(p)
	o.TotalBytes += int64(n)
	o.TotalOps += 1

	return n, err
}

func (o *readerWriterCounter) Write(p []byte) (int, error) {
	o.Mutex.Lock()
	defer o.Mutex.Unlock()
	n, err := o.rw.Write(p)
	o.TotalBytes += int64(n)
	o.TotalOps += 1

	return n, err
}

func (o *readerWriterCounter) ReadCount() (n int64, nops int) {
	o.Mutex.Lock()
	defer o.Mutex.Unlock()
	return o.TotalBytes, o.TotalOps
}

func (o *readerWriterCounter) WriteCount() (n int64, nops int) {
	o.Mutex.Lock()
	defer o.Mutex.Unlock()
	return o.TotalBytes, o.TotalOps
}

func NewReadWriteCounter(o ReaderWriter) ReadWriteCounter {
	return &readerWriterCounter{rw: o, Mutex: new(sync.Mutex)}
}

type ReaderWriter interface {
	io.Reader
	io.Writer
}
