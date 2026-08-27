package model

import (
	"sort"
	"strings"
	"time"
)

type QualityBand string

const (
	QualityWeak    QualityBand = "weak"
	QualityUsable  QualityBand = "usable"
	QualityStrong  QualityBand = "strong"
	QualityCurated QualityBand = "curated"
)

type Quality struct {
	Band      QualityBand `json:"band"`
	Score     int         `json:"score"`
	Signals   []string    `json:"signals"`
	Reviewed  bool        `json:"reviewed"`
	Published bool        `json:"published"`
	WordCount int         `json:"word_count"`
	TagCount  int         `json:"tag_count"`
	AgeHours  int         `json:"age_hours"`
}

func Evaluate(item Observation, now time.Time) Quality {
	quality := Quality{Band: QualityWeak, Signals: []string{}}
	if item.HasDescription() {
		quality.Score += 30
		quality.Signals = append(quality.Signals, "description")
	}
	if len(item.Tags) > 0 {
		quality.Score += 20
		quality.Signals = append(quality.Signals, "tags")
	}
	if item.Review != nil && item.Review.Approved() {
		quality.Score += 25
		quality.Reviewed = true
		quality.Signals = append(quality.Signals, "approval")
	}
	if item.Bulletin != nil && item.Bulletin.Valid() {
		quality.Score += 25
		quality.Published = true
		quality.Signals = append(quality.Signals, "bulletin")
	}
	quality.WordCount = len(strings.Fields(item.Description))
	quality.TagCount = len(item.Tags)
	if !item.ObservedAt.IsZero() && !now.Before(item.ObservedAt) {
		quality.AgeHours = int(now.Sub(item.ObservedAt) / time.Hour)
	}
	quality.Band = BandForScore(quality.Score)
	return quality
}

func BandForScore(score int) QualityBand {
	switch {
	case score >= 90:
		return QualityCurated
	case score >= 60:
		return QualityStrong
	case score >= 30:
		return QualityUsable
	default:
		return QualityWeak
	}
}

func (q Quality) ReadyForRelease() bool { return q.Band == QualityStrong || q.Band == QualityCurated }
func (q Quality) Complete() bool        { return q.Reviewed && q.Published }
func (q Quality) SignalCount() int      { return len(q.Signals) }
func (q Quality) HasSignal(value string) bool {
	for _, signal := range q.Signals {
		if signal == value {
			return true
		}
	}
	return false
}
func (q Quality) Summary() string { return string(q.Band) + ":" + strings.Join(q.Signals, ",") }

func QualityTable(items []Observation, now time.Time) map[string]Quality {
	result := make(map[string]Quality, len(items))
	for _, item := range items {
		result[item.ID] = Evaluate(item, now)
	}
	return result
}

func QualityOrder(items []Observation, now time.Time) []Observation {
	result := append([]Observation{}, items...)
	sort.SliceStable(result, func(i, j int) bool {
		left := Evaluate(result[i], now)
		right := Evaluate(result[j], now)
		if left.Score == right.Score {
			return result[i].ID < result[j].ID
		}
		return left.Score > right.Score
	})
	return result
}

func QualityScores(items []Observation, now time.Time) []int {
	ordered := QualityOrder(items, now)
	result := make([]int, 0, len(ordered))
	for _, item := range ordered {
		result = append(result, Evaluate(item, now).Score)
	}
	return result
}

func Strongest(items []Observation, now time.Time) (Observation, Quality, bool) {
	ordered := QualityOrder(items, now)
	if len(ordered) == 0 {
		return Observation{}, Quality{}, false
	}
	return ordered[0], Evaluate(ordered[0], now), true
}

func Weakest(items []Observation, now time.Time) (Observation, Quality, bool) {
	ordered := QualityOrder(items, now)
	if len(ordered) == 0 {
		return Observation{}, Quality{}, false
	}
	last := ordered[len(ordered)-1]
	return last, Evaluate(last, now), true
}

func ReleaseReady(items []Observation, now time.Time) []string {
	result := []string{}
	for _, item := range items {
		if Evaluate(item, now).ReadyForRelease() {
			result = append(result, item.ID)
		}
	}
	sort.Strings(result)
	return result
}
