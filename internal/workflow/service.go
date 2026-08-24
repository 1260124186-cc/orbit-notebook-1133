package workflow

import (
	"example.com/orbit-notebook/internal/archive"
	"example.com/orbit-notebook/internal/clock"
	"example.com/orbit-notebook/internal/index"
	"example.com/orbit-notebook/internal/integrity"
	"example.com/orbit-notebook/internal/model"
	"example.com/orbit-notebook/internal/policy"
	"example.com/orbit-notebook/internal/query"
	"example.com/orbit-notebook/internal/stats"
	"example.com/orbit-notebook/internal/storage"
	"example.com/orbit-notebook/internal/timeline"
	"example.com/orbit-notebook/internal/validation"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	store *storage.Store
	clock clock.Clock
}

func New(store *storage.Store, c clock.Clock) *Service { return &Service{store: store, clock: c} }
func (s *Service) Capture(in validation.CaptureInput) (model.Observation, error) {
	if err := validation.Capture(in); err != nil {
		return model.Observation{}, err
	}
	var created model.Observation
	err := s.store.Replace(func(items []model.Observation) ([]model.Observation, error) {
		allocatedID := index.NextID(len(items))
		now := s.clock.Now()
		created = model.Observation{ID: allocatedID, Target: strings.TrimSpace(in.Target), Instrument: strings.TrimSpace(in.Instrument), ObservedAt: in.ObservedAt.UTC(), Description: strings.TrimSpace(in.Description), State: model.StateDraft, CreatedAt: now, UpdatedAt: now}
		created = integrity.Normalize(created)
		return append(items, created), nil
	})
	return created, err
}
func (s *Service) Annotate(id, description string, tags []string) (model.Observation, error) {
	if err := validation.Annotation(description); err != nil {
		return model.Observation{}, err
	}
	var result model.Observation
	err := s.store.Replace(func(items []model.Observation) ([]model.Observation, error) {
		for i := range items {
			if items[i].ID != id {
				continue
			}
			if !policy.CanAnnotate(items[i]) {
				return nil, fmt.Errorf(policy.ErrImmutable)
			}
			items[i].Description = strings.TrimSpace(items[i].Description + " " + strings.TrimSpace(description))
			items[i].Tags = normalizeTags(append(items[i].Tags, tags...))
			items[i].Description = strings.TrimSpace(items[i].Description)
			items[i].UpdatedAt = s.clock.Now()
			result = items[i]
			return items, nil
		}
		return nil, fmt.Errorf(policy.ErrNotFound)
	})
	return result, err
}
func normalizeTags(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		v := strings.ToLower(strings.TrimSpace(value))
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
func (s *Service) Review(id, reviewer, decision, note string) (model.Observation, error) {
	d, err := validation.Decision(decision)
	if err != nil {
		return model.Observation{}, err
	}
	if err = validation.Reviewer(reviewer, note); err != nil {
		return model.Observation{}, err
	}
	var result model.Observation
	err = s.store.Replace(func(items []model.Observation) ([]model.Observation, error) {
		for i := range items {
			if items[i].ID != id {
				continue
			}
			if err := validation.ReadyForReview(items[i]); err != nil {
				return nil, err
			}
			items[i].Review = &model.Review{Reviewer: strings.TrimSpace(reviewer), Decision: d, Note: strings.TrimSpace(note), ReviewedAt: s.clock.Now()}
			if d == model.DecisionApprove {
				items[i].State = model.StateReviewed
			}
			items[i].UpdatedAt = s.clock.Now()
			result = items[i]
			return items, nil
		}
		return nil, fmt.Errorf(policy.ErrNotFound)
	})
	return result, err
}
func (s *Service) Release(id, publisher string) (model.Observation, error) {
	if strings.TrimSpace(publisher) == "" {
		return model.Observation{}, fmt.Errorf("publisher is required")
	}
	var result model.Observation
	err := s.store.Replace(func(items []model.Observation) ([]model.Observation, error) {
		for i := range items {
			if items[i].ID != id {
				continue
			}
			if err := validation.ReadyForRelease(items[i]); err != nil {
				return nil, err
			}
			now := s.clock.Now()
			bulletin := model.NewBulletin(items[i], publisher, now)
			items[i].Bulletin = &bulletin
			items[i].State = model.StateReleased
			items[i].UpdatedAt = now
			result = items[i]
			return items, nil
		}
		return nil, fmt.Errorf(policy.ErrNotFound)
	})
	return result, err
}
func (s *Service) List(f index.Filter) ([]model.Summary, error) {
	items, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	return index.Apply(items, f), nil
}
func (s *Service) Show(id string) (model.Observation, error) {
	items, err := s.store.Load()
	if err != nil {
		return model.Observation{}, err
	}
	for _, item := range items {
		if item.ID == id {
			return item, nil
		}
	}
	return model.Observation{}, fmt.Errorf(policy.ErrNotFound)
}

// Search applies the richer query model while preserving the workflow boundary.
func (s *Service) Search(request query.Request) (query.Page, error) {
	items, err := s.store.Load()
	if err != nil {
		return query.Page{}, err
	}
	return query.Search(items, request, query.SortObserved), nil
}

func (s *Service) Overview() (stats.Overview, error) {
	items, err := s.store.Load()
	if err != nil {
		return stats.Overview{}, err
	}
	return stats.Build(items), nil
}

func (s *Service) TargetCounts() ([]stats.TargetCount, error) {
	items, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	return stats.ByTarget(items), nil
}

func (s *Service) Timeline(id string) ([]timeline.Event, error) {
	item, err := s.Show(id)
	if err != nil {
		return nil, err
	}
	return timeline.ForObservation(item), nil
}

func (s *Service) Export(path string) error {
	return s.store.Export(path, s.clock.Now())
}

func (s *Service) Import(path string) error {
	return s.store.Import(path)
}

func (s *Service) ArchivePreview() (string, error) {
	items, err := s.store.Load()
	if err != nil {
		return "", err
	}
	return archive.Markdown(items), nil
}

func (s *Service) ObservedBetween(start, end time.Time) ([]model.Summary, error) {
	items, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	result := make([]model.Summary, 0, len(items))
	for _, item := range items {
		if !item.ObservedAt.Before(start) && item.ObservedAt.Before(end) {
			result = append(result, item.Summary())
		}
	}
	model.SortSummaries(result)
	return result, nil
}
