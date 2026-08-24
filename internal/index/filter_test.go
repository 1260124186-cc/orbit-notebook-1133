package index

import (
	"example.com/orbit-notebook/internal/model"
	"testing"
	"time"
)

func TestApplyFiltersAndSort(t *testing.T) {
	items := []model.Observation{{ID: "obs-2", Target: "M31", ObservedAt: time.Unix(2, 0), State: model.StateDraft, Tags: []string{"galaxy"}}, {ID: "obs-1", Target: "M31", ObservedAt: time.Unix(1, 0), State: model.StateReviewed, Tags: []string{"galaxy"}}}
	got := Apply(items, Filter{Target: "m31", Tag: "galaxy"})
	if len(got) != 2 || got[0].ID != "obs-1" {
		t.Fatalf("unexpected %#v", got)
	}
}
