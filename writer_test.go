package paasio

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"testing"
)

// this test could be improved to test that error conditions are preserved.
func testWrite(t *testing.T, writer func(io.Writer) WriteCounter) {
	for i, test := range []struct {
		writes []string
	}{
		{nil},
		{[]string{""}},
		{[]string{"I", " ", "never met ", "", "a gohper"}},
	} {
		var buf bytes.Buffer
		buft := writer(&buf)
		for _, s := range test.writes {
			n, err := buft.Write([]byte(s))
			if err != nil {
				t.Errorf("test %d: Write(%q) unexpected error: %v", i, s, err)
				continue
			}
			if n != len(s) {
				t.Errorf("test %d: Write(%q) unexpected number of bytes written: %v", i, s, n)
				continue
			}
		}
		out := buf.String()
		if out != strings.Join(test.writes, "") {
			t.Errorf("test %d: unexpected content in underlying writer: %q", i, out)
		}
	}
}

func TestWriteWriter(t *testing.T) {
	testWrite(t, NewWriteCounter)
}

func TestWriteReadWriter(t *testing.T) {
	testWrite(t, func(w io.Writer) WriteCounter {
		var r nopReader
		return NewReadWriteCounter(readWriter{r, w})
	})
}

func testWriteTotal(t *testing.T, wt WriteCounter) {
	numGo := 8000
	numBytes := 50
	totalBytes := int64(numGo) * int64(numBytes)
	p := make([]byte, numBytes)

	t.Logf("Calling Write() with %d*%d=%d bytes", numGo, numBytes, totalBytes)
	wg := new(sync.WaitGroup)
	wg.Add(numGo)
	start := make(chan struct{})
	for i := 0; i < numGo; i++ {
		go func() {
			<-start
			wt.Write(p)
			wg.Done()
		}()
	}
	close(start)

	wg.Wait()
	n, nops := wt.WriteCount()
	if n != totalBytes {
		t.Errorf("expected %d bytes written; %d bytes reported", totalBytes, n)
	}
	if nops != numGo {
		t.Errorf("expected %d write operations; %d operations reported", numGo, nops)
	}
}

func TestWriteTotalWriter(t *testing.T) {
	var w nopWriter
	testWriteTotal(t, NewWriteCounter(w))
}

func TestWriteTotalReadWriter(t *testing.T) {
	var rw nopReaderWriter
	testWriteTotal(t, NewReadWriteCounter(rw))
}

func TestWriteCountConsistencyWriter(t *testing.T) {
	var w nopWriter
	testWriteCountConsistency(t, NewWriteCounter(w))
}

func TestWriteCountConsistencyReadWriter(t *testing.T) {
	var rw nopReaderWriter
	testWriteCountConsistency(t, NewReadWriteCounter(rw))
}

func testWriteCountConsistency(t *testing.T, wc WriteCounter) {
	const numGo = 4000
	const numBytes = 50
	p := make([]byte, numBytes)

	wg := new(sync.WaitGroup)
	wg.Add(2 * numGo)
	start := make(chan struct{})
	for i := 0; i < numGo; i++ {
		go func() {
			<-start
			wc.Write(p)
			wg.Done()
		}()
		go func() {
			<-start
			n, nops := wc.WriteCount()
			expectedOps := n / numBytes
			if int64(nops) != n/numBytes {
				t.Errorf("expected %d nops@%d bytes written; %d ops reported", expectedOps, n, nops)
			}
			wg.Done()
		}()
	}
	close(start)
	wg.Wait()
}
