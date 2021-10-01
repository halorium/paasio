package paasio

import (
	"bytes"
	"crypto/rand"
	"io"
	"sync"
	"testing"
)

// this test could be improved to test exact number of operations as well as
// ensure that error conditions are preserved.
func testRead(t *testing.T, reader func(io.Reader) ReadCounter) {
	chunkLen := 10 << 20 // 10MB
	orig := make([]byte, 10<<20)
	_, err := rand.Read(orig)
	if err != nil {
		t.Fatalf("error reading random data")
	}
	buf := bytes.NewBuffer(orig)
	rc := reader(buf)
	var obuf bytes.Buffer
	ncopy, err := io.Copy(&obuf, rc)
	if err != nil {
		t.Fatalf("error reading: %v", err)
	}
	if ncopy != int64(chunkLen) {
		t.Fatalf("copied %d bytes instead of %d", ncopy, chunkLen)
	}
	if string(orig) != obuf.String() {
		t.Fatalf("unexpected output from Read()")
	}
	n, nops := rc.ReadCount()
	if n != int64(chunkLen) {
		t.Fatalf("reported %d bytes read instead of %d", n, chunkLen)
	}
	if nops < 2 {
		t.Fatalf("unexpected number of reads: %v", nops)
	}
}

func TestReadReader(t *testing.T) {
	testRead(t, NewReadCounter)
}

func TestReadReadWriter(t *testing.T) {
	testRead(t, func(r io.Reader) ReadCounter {
		var w nopWriter
		return NewReadWriteCounter(readerWriter{r, w})
	})
}

func testReadTotal(t *testing.T, rc ReadCounter) {
	numGo := 8000
	numBytes := 50
	totalBytes := int64(numGo) * int64(numBytes)
	p := make([]byte, numBytes)

	t.Logf("Calling Read() for %d*%d=%d bytes", numGo, numBytes, totalBytes)
	wg := new(sync.WaitGroup)
	wg.Add(numGo)
	start := make(chan struct{})
	for i := 0; i < numGo; i++ {
		go func() {
			<-start
			rc.Read(p)
			wg.Done()
		}()
	}
	close(start)

	wg.Wait()
	n, nops := rc.ReadCount()
	if n != totalBytes {
		t.Errorf("expected %d bytes read; %d bytes reported", totalBytes, n)
	}
	if nops != numGo {
		t.Errorf("expected %d read operations; %d operations reported", numGo, nops)
	}
}

func TestReadTotalReader(t *testing.T) {
	var r nopReader
	testReadTotal(t, NewReadCounter(r))
}

func TestReadTotalReadWriter(t *testing.T) {
	var rw nopReaderWriter
	testReadTotal(t, NewReadWriteCounter(rw))
}

func TestReadCountConsistencyReader(t *testing.T) {
	var r nopReader
	testReadCountConsistency(t, NewReadCounter(r))
}

func TestReadCountConsistencyReadWriter(t *testing.T) {
	var rw nopReaderWriter
	testReadCountConsistency(t, NewReadWriteCounter(rw))
}

func testReadCountConsistency(t *testing.T, rc ReadCounter) {
	const numGo = 4000
	const numBytes = 50

	p := make([]byte, numBytes)

	wg := new(sync.WaitGroup)
	wg.Add(2 * numGo)
	start := make(chan struct{})
	for i := 0; i < numGo; i++ {
		go func() {
			<-start
			rc.Read(p)
			wg.Done()
		}()
		go func() {
			<-start
			n, nops := rc.ReadCount()
			expectedOps := n / numBytes
			if int64(nops) != expectedOps {
				t.Errorf("expected %d ops@%d bytes read; %d ops reported", expectedOps, n, nops)
			}
			wg.Done()
		}()
	}
	close(start)
	wg.Wait()
}
