package report

import (
	"example.com/orbit-notebook/internal/model"
	"io"
)

func Summaries(w io.Writer, items []model.Summary) error {
	return WriteJSON(w, map[string]any{"status": "ok", "count": len(items), "items": items})
}
func Detail(w io.Writer, item model.Observation) error {
	return WriteJSON(w, map[string]any{"status": "ok", "observation": item})
}
