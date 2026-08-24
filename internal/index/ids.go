package index

import "fmt"

func NextID(items int) string { return fmt.Sprintf("obs-%06d", items+1) }
func FindPosition[T interface{ GetID() string }](items []T, id string) int {
	for i, item := range items {
		if item.GetID() == id {
			return i
		}
	}
	return -1
}
