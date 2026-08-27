package report

import (
	"encoding/json"
	"example.com/orbit-notebook/internal/model"
	"io"
)

func WriteJSON(w io.Writer, value any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(value)
}
func Created(w io.Writer, o model.Observation) error {
	return WriteJSON(w, map[string]any{"status": "created", "observation": o.Summary()})
}
func Annotated(w io.Writer, o model.Observation) error {
	return WriteJSON(w, map[string]any{"status": "annotated", "observation": o.Summary()})
}
func Reviewed(w io.Writer, o model.Observation) error {
	return WriteJSON(w, map[string]any{"status": "reviewed", "observation": o.Summary()})
}
func Released(w io.Writer, o model.Observation) error {
	return WriteJSON(w, map[string]any{"status": "released", "observation": o.Summary()})
}
