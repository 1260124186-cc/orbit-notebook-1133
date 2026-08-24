package archive

import (
	"encoding/json"
	"example.com/orbit-notebook/internal/model"
	"fmt"
	"time"
)

type Document struct {
	Format      string              `json:"format"`
	GeneratedAt time.Time           `json:"generated_at"`
	Items       []model.Observation `json:"items"`
}

func New(items []model.Observation, now time.Time) Document {
	copyItems := append([]model.Observation{}, items...)
	return Document{Format: "orbit-notebook/v1", GeneratedAt: now.UTC(), Items: copyItems}
}
func Encode(document Document) ([]byte, error) { return json.MarshalIndent(document, "", "  ") }
func Decode(data []byte) (Document, error) {
	var result Document
	if err := json.Unmarshal(data, &result); err != nil {
		return Document{}, fmt.Errorf("decode archive: %w", err)
	}
	if result.Format != "orbit-notebook/v1" {
		return Document{}, fmt.Errorf("unsupported archive format %q", result.Format)
	}
	return result, nil
}
func Validate(document Document) error {
	if document.Format != "orbit-notebook/v1" {
		return fmt.Errorf("archive format is required")
	}
	if document.GeneratedAt.IsZero() {
		return fmt.Errorf("archive generated time is required")
	}
	for _, item := range document.Items {
		if item.ID == "" {
			return fmt.Errorf("archive contains item without id")
		}
	}
	return nil
}
