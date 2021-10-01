package paasio

import (
	"runtime"
	"testing"
)

func TestMultiThreaded(t *testing.T) {
	mincpu := 2
	minproc := 2
	ncpu := runtime.NumCPU()
	if ncpu < mincpu {
		t.Fatalf("at least %d cpu cores are required", mincpu)
	}
	nproc := runtime.GOMAXPROCS(0)
	if nproc < minproc {
		t.Errorf("at least %d threads are required; rerun the tests", minproc)
		t.Errorf("")
		t.Errorf("\tgo test -cpu %d ...", minproc)
	}
}
