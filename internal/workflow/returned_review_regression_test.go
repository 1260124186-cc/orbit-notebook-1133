package workflow

import (
	"example.com/orbit-notebook/internal/clock"
	"example.com/orbit-notebook/internal/integrity"
	"example.com/orbit-notebook/internal/model"
	"example.com/orbit-notebook/internal/storage"
	"example.com/orbit-notebook/internal/validation"
	"testing"
	"time"
)

func TestReturnedReviewCanBeAmendedAndReviewedAgain(t *testing.T) {
	store, err := storage.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	svc := New(store, clock.Fixed{Value: time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)})
	o, err := svc.Capture(validation.CaptureInput{Target: "M42", Instrument: "scope-a", ObservedAt: time.Date(2026, 8, 24, 9, 0, 0, 0, time.UTC), Description: "初始观测"})
	if err != nil {
		t.Fatal(err)
	}
	returned, err := svc.Review(o.ID, "alice", "return", "请补充细节")
	if err != nil {
		t.Fatal(err)
	}
	if returned.State != model.StateDraft {
		t.Fatalf("returned state=%q want draft", returned.State)
	}
	if err := integrity.Consistent(returned); err != nil {
		t.Fatalf("returned observation inconsistent: %v", err)
	}
	if _, err := svc.Annotate(o.ID, "补充结构", []string{"deep-sky"}); err != nil {
		t.Fatalf("returned observation cannot be amended: %v", err)
	}
	if _, err := svc.Review(o.ID, "alice", "approve", "信息完整"); err != nil {
		t.Fatalf("returned observation cannot be reviewed again: %v", err)
	}
	if _, err := svc.Release(o.ID, "bob"); err != nil {
		t.Fatalf("approved observation cannot be released: %v", err)
	}
}
