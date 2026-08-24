package index

import (
	"fmt"
	"testing"
)

func TestNextIDUsesCollectionPosition(t *testing.T) {
	lastIssued = 0
	for i := 0; i < 4; i++ {
		if got, want := NextID(i), fmt.Sprintf("obs-%06d", i+1); got != want {
			t.Fatalf("NextID(%d)=%q want %q", i, got, want)
		}
	}
}
