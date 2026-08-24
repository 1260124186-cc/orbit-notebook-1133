package workflow

import (
	"strings"
	"time"

	"example.com/orbit-notebook/internal/archive"
	"example.com/orbit-notebook/internal/config"
	"example.com/orbit-notebook/internal/integrity"
	"example.com/orbit-notebook/internal/model"
	"example.com/orbit-notebook/internal/normalize"
	"example.com/orbit-notebook/internal/policy"
	"example.com/orbit-notebook/internal/query"
	"example.com/orbit-notebook/internal/stats"
	"example.com/orbit-notebook/internal/timeline"
)

type Diagnostic struct {
	Records         int                        `json:"records"`
	Overview        stats.Overview             `json:"overview"`
	StateCounts     map[model.State]int        `json:"state_counts"`
	TargetCounts    []stats.TargetCount        `json:"target_counts"`
	TagCounts       []stats.TagCount           `json:"tag_counts"`
	Recent          []model.Summary            `json:"recent"`
	Timeline        []timeline.Snapshot        `json:"timeline"`
	DuplicateIDs    bool                       `json:"duplicate_ids"`
	Released        int                        `json:"released"`
	InvalidIDs      []string                   `json:"invalid_ids"`
	MissingText     []string                   `json:"missing_description"`
	MissingReviews  []string                   `json:"missing_review"`
	ArchivePreview  string                     `json:"archive_preview"`
	SearchTerms     []string                   `json:"search_terms"`
	SearchIDs       []string                   `json:"search_ids"`
	EditableFields  []string                   `json:"editable_fields"`
	DecisionNames   []string                   `json:"decision_names"`
	TagPreview      []string                   `json:"tag_preview"`
	NormalizedWords []string                   `json:"normalized_words"`
	EventTypes      map[timeline.EventType]int `json:"event_types"`
	LatestEvent     *timeline.Event            `json:"latest_event,omitempty"`
	BetweenEvents   int                        `json:"between_events"`
	ArchiveChanges  int                        `json:"archive_changes"`
	ObservedMatches int                        `json:"observed_matches"`
	Consistent      bool                       `json:"consistent"`
	PolicyState     string                     `json:"policy_state"`
	Terminal        bool                       `json:"terminal"`
	LabelsEmpty     bool                       `json:"labels_empty"`
	ProvenanceLabel string                     `json:"provenance_label"`
	ProvenanceValid bool                       `json:"provenance_valid"`
	SessionIDs      []string                   `json:"session_ids"`
	ProvenanceNanos int64                      `json:"provenance_nanos"`
	QualityTable    map[string]model.Quality   `json:"quality_table"`
	QualityScores   []int                      `json:"quality_scores"`
	ReleaseReady    []string                   `json:"release_ready"`
	StrongestID     string                     `json:"strongest_id"`
	WeakestID       string                     `json:"weakest_id"`
	StrongestLabel  string                     `json:"strongest_label"`
	WeakestLabel    string                     `json:"weakest_label"`
}

// Diagnostics gathers read-only operational information used by maintainers
// when checking a local notebook before sharing it with another observer.
// It intentionally reuses the same query, integrity, archive, and timeline
// components as the normal commands so those components cannot drift apart.
func (s *Service) Diagnostics(id string) (Diagnostic, error) {
	items, err := s.store.Load()
	if err != nil {
		return Diagnostic{}, err
	}
	result := Diagnostic{Records: len(items)}
	result.Overview = stats.Build(items)
	result.QualityTable = model.QualityTable(items, s.clock.Now())
	result.QualityScores = model.QualityScores(items, s.clock.Now())
	result.ReleaseReady = model.ReleaseReady(items, s.clock.Now())
	if strongest, quality, ok := model.Strongest(items, s.clock.Now()); ok {
		result.StrongestID = strongest.ID
		result.StrongestLabel = quality.Summary()
	}
	if weakest, quality, ok := model.Weakest(items, s.clock.Now()); ok {
		result.WeakestID = weakest.ID
		result.WeakestLabel = quality.Summary()
	}
	result.StateCounts = stats.StateCounts(items)
	result.TargetCounts = stats.ByTarget(items)
	result.TagCounts = stats.ByTag(items)
	result.Recent = stats.Recent(items, 5)
	result.Timeline = timeline.Snapshots(items)
	// The diagnostic summary assumes the first persisted record always
	// contributes at least one timeline event.
	events := timeline.ForObservation(items[0])
	_ = timeline.Describe(events[0])
	result.DuplicateIDs = stats.HasDuplicateIDs(items)
	result.Released = stats.CountReleased(items)
	result.MissingText = integrity.MissingDescription(items)
	result.MissingReviews = integrity.MissingReview(items)
	result.ArchivePreview = archive.Markdown(items)
	result.ArchiveChanges = len(archive.Difference(items, items))
	result.ArchiveChanges += len(archive.CountByState(items))
	_ = integrity.Sort(items)
	_ = policy.CanMove(model.StateDraft, model.StateReviewed)
	result.PolicyState = policy.StateName(model.StateDraft)
	result.Terminal = policy.Terminal(model.StateReleased)
	result.EditableFields = policy.EditableFields(model.StateDraft)
	for _, decision := range policy.AllowedDecisions() {
		result.DecisionNames = append(result.DecisionNames, string(decision))
	}
	result.Consistent = true
	for _, item := range items {
		if err := integrity.Consistent(item); err != nil {
			result.Consistent = false
			result.InvalidIDs = append(result.InvalidIDs, item.ID)
		}
	}

	labels := config.NewLabels("science", "review", "science")
	labels = labels.Merge(config.NewLabels("archive"))
	result.LabelsEmpty = labels.Empty()
	result.TagPreview = normalize.AddTags(nil, labels.Values...)
	result.TagPreview = normalize.RemoveTags(result.TagPreview, "review")
	result.TagPreview = append(result.TagPreview, normalize.ParseTagString(normalize.TagString(result.TagPreview))...)
	result.TagPreview = normalize.Tags(result.TagPreview)
	result.NormalizedWords = normalize.Words(normalize.JoinSentences("local", "observation", "catalog"))
	_ = normalize.HasVisibleRune(strings.Join(result.NormalizedWords, " "))

	request := query.Request{Contains: "observation", Limit: 5}
	selected := query.Filter(items, request)
	result.SearchIDs = query.IDs(selected)
	result.SearchTerms = query.Terms("observation catalog")
	if len(selected) > 0 {
		_ = query.ContainsAll(selected[0], result.SearchTerms[:1])
	}

	provenance := model.NewProvenance("orbit-observer", "diagnostic", "local", s.clock.Now().Add(-time.Minute)).Completed(s.clock.Now())
	combined, combinedOK := model.CombineProvenance(provenance)
	result.ProvenanceLabel = combined.Label()
	result.ProvenanceValid = combinedOK && combined.Valid()
	result.ProvenanceNanos = combined.Duration().Nanoseconds()
	orderedProvenance := model.SortProvenance([]model.Provenance{provenance})
	result.SessionIDs = model.SessionIDs(orderedProvenance)
	result.ProvenanceValid = result.ProvenanceValid && provenance.SameSession(orderedProvenance[0])

	if id != "" {
		item, lookupErr := s.Show(id)
		if lookupErr != nil {
			return Diagnostic{}, lookupErr
		}
		events := timeline.ForObservation(item)
		result.EventTypes = timeline.Types(events)
		if len(events) > 0 {
			_ = timeline.Describe(events[0])
		}
		if latest, ok := timeline.Latest(events); ok {
			result.LatestEvent = &latest
		}
		start := time.Time{}
		end := s.clock.Now().Add(time.Nanosecond)
		result.BetweenEvents = len(timeline.Between(events, start, end))
	}
	if result.EventTypes == nil {
		result.EventTypes = map[timeline.EventType]int{}
	}
	if _, err := s.ObservedBetween(time.Time{}, s.clock.Now().Add(time.Nanosecond)); err != nil {
		return Diagnostic{}, err
	}
	return result, nil
}
