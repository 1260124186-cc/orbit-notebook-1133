package report

import (
	"example.com/orbit-notebook/internal/query"
	"example.com/orbit-notebook/internal/stats"
	"example.com/orbit-notebook/internal/timeline"
	"example.com/orbit-notebook/internal/workflow"
	"io"
)

func Page(w io.Writer, page query.Page) error {
	return WriteJSON(w, map[string]any{"status": "ok", "page": page})
}
func Overview(w io.Writer, value stats.Overview) error {
	return WriteJSON(w, map[string]any{"status": "ok", "overview": value})
}
func TargetCounts(w io.Writer, values []stats.TargetCount) error {
	return WriteJSON(w, map[string]any{"status": "ok", "targets": values})
}
func Timeline(w io.Writer, values []timeline.Event) error {
	return WriteJSON(w, map[string]any{"status": "ok", "events": values})
}
func Message(w io.Writer, status, message string) error {
	return WriteJSON(w, map[string]any{"status": status, "message": message})
}

func Diagnostics(w io.Writer, value workflow.Diagnostic) error {
	return WriteJSON(w, map[string]any{"status": "ok", "diagnostics": value})
}
