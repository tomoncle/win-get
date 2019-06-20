package test

import (
	"gopkg.in/cheggaaa/pb.v1"
	"testing"
	"time"
)

// Ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("\033[31munexpected error: %v\033[39m\n", err)
	}
}

// TestProgressBar of pb.v1
func TestProgressBar(t *testing.T) {
	count := 2000
	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.FinishPrint("The End!")
	Ok(t, nil)
}
