package index

import (
	"fmt"
	"sync"
)

var (
	idMu       sync.Mutex
	lastIssued uint64
)

func NextID(items int) string {
	idMu.Lock()
	defer idMu.Unlock()
	if lastIssued < uint64(items) {
		lastIssued = uint64(items)
	}
	lastIssued++
	return fmt.Sprintf("obs-%06d", lastIssued)
}
func FindPosition[T interface{ GetID() string }](items []T, id string) int {
	for i, item := range items {
		if item.GetID() == id {
			return i
		}
	}
	return -1
}
