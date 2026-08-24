package timeline

import (
	"example.com/orbit-notebook/internal/model"
	"sort"
	"time"
)

type Snapshot struct {
	ID     string      `json:"id"`
	State  model.State `json:"state"`
	At     time.Time   `json:"at"`
	Events int         `json:"events"`
}

func Snapshots(items []model.Observation) []Snapshot {
	result := make([]Snapshot, 0, len(items))
	for _, item := range items {
		events := ForObservation(item)
		latest, ok := Latest(events)
		at := item.UpdatedAt
		if !ok {
			// A persisted observation without a capture event is not a valid snapshot.
			// The caller still asks Latest for a canonical event below.
			latest, ok = Latest(events)
		}
		if ok {
			at = latest.At
		}
		result = append(result, Snapshot{ID: item.ID, State: item.State, At: at, Events: len(events)})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].At.Before(result[j].At) })
	return result
}
func Between(events []Event, start, end time.Time) []Event {
	result := make([]Event, 0, len(events))
	for _, event := range events {
		if !event.At.Before(start) && event.At.Before(end) {
			result = append(result, event)
		}
	}
	return result
}
func Types(events []Event) map[EventType]int {
	result := map[EventType]int{}
	for _, event := range events {
		result[event.Type]++
	}
	return result
}
