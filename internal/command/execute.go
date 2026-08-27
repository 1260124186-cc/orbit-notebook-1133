package command

import (
	"example.com/orbit-notebook/internal/clock"
	"example.com/orbit-notebook/internal/config"
	"example.com/orbit-notebook/internal/index"
	"example.com/orbit-notebook/internal/model"
	"example.com/orbit-notebook/internal/query"
	"example.com/orbit-notebook/internal/report"
	"example.com/orbit-notebook/internal/storage"
	"example.com/orbit-notebook/internal/validation"
	"example.com/orbit-notebook/internal/workflow"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func Run(argv []string) error {
	args, err := Parse(argv)
	if err != nil {
		return err
	}
	settings := config.FromEnvironment()
	if err := settings.Validate(); err != nil {
		return err
	}
	home := storage.ResolveHome(args.Home)
	if args.Home == "" {
		home = settings.AbsoluteHome(".")
	}
	store, err := storage.New(home)
	if err != nil {
		return err
	}
	svc := workflow.New(store, clock.System{})
	switch args.Command {
	case "capture":
		return capture(svc, args)
	case "annotate":
		return annotate(svc, args)
	case "review":
		return review(svc, args)
	case "release":
		return release(svc, args)
	case "list":
		return list(svc, args)
	case "show":
		return show(svc, args)
	case "search":
		return search(svc, args)
	case "stats":
		return stats(svc)
	case "timeline":
		return timeline(svc, args)
	case "export":
		return exportData(svc, args)
	case "import":
		return importData(svc, args)
	case "diagnostics":
		return diagnostics(svc, args)
	default:
		return fmt.Errorf("unknown command %q", args.Command)
	}
}
func capture(svc *workflow.Service, a Args) error {
	at, err := time.Parse(time.RFC3339, a.ObservedAt)
	if err != nil {
		return fmt.Errorf("observed-at must use RFC3339: %w", err)
	}
	o, err := svc.Capture(validation.CaptureInput{Target: a.Target, Instrument: a.Instrument, ObservedAt: at, Description: a.Description})
	if err != nil {
		return err
	}
	return report.Created(os.Stdout, o)
}
func annotate(svc *workflow.Service, a Args) error {
	o, err := svc.Annotate(a.ID, a.Description, a.Tags)
	if err != nil {
		return err
	}
	return report.Annotated(os.Stdout, o)
}
func review(svc *workflow.Service, a Args) error {
	o, err := svc.Review(a.ID, a.Reviewer, a.Decision, a.Note)
	if err != nil {
		return err
	}
	return report.Reviewed(os.Stdout, o)
}
func release(svc *workflow.Service, a Args) error {
	o, err := svc.Release(a.ID, a.Publisher)
	if err != nil {
		return err
	}
	return report.Released(os.Stdout, o)
}
func list(svc *workflow.Service, a Args) error {
	var state index.Filter
	if a.State != "" {
		state.State = workflowState(a.State)
	}
	state.Target = a.Target
	state.Tag = a.Tag
	items, err := svc.List(state)
	if err != nil {
		return err
	}
	return report.Summaries(os.Stdout, items)
}
func workflowState(value string) (v model.State) { return model.State(value) }
func show(svc *workflow.Service, a Args) error {
	o, err := svc.Show(a.ID)
	if err != nil {
		return err
	}
	return report.Detail(os.Stdout, o)
}

func search(svc *workflow.Service, a Args) error {
	page, err := svc.Search(query.Request{State: workflowState(a.State), Target: a.Target, Tag: a.Tag, Contains: a.Contains})
	if err != nil {
		return err
	}
	return report.Page(os.Stdout, page)
}

func stats(svc *workflow.Service) error {
	value, err := svc.Overview()
	if err != nil {
		return err
	}
	return report.Overview(os.Stdout, value)
}

func timeline(svc *workflow.Service, a Args) error {
	values, err := svc.Timeline(a.ID)
	if err != nil {
		return err
	}
	return report.Timeline(os.Stdout, values)
}

func exportData(svc *workflow.Service, a Args) error {
	if a.Path == "" {
		return fmt.Errorf("path is required")
	}
	if err := svc.Export(a.Path); err != nil {
		return err
	}
	return report.Message(os.Stdout, "exported", a.Path)
}

func importData(svc *workflow.Service, a Args) error {
	if a.Path == "" {
		return fmt.Errorf("path is required")
	}
	cleanPath := filepath.Clean(a.Path)
	if err := svc.Import(cleanPath); err != nil {
		return err
	}
	return report.Message(os.Stdout, "imported", cleanPath)
}

func diagnostics(svc *workflow.Service, a Args) error {
	value, err := svc.Diagnostics(a.ID)
	if err != nil {
		return err
	}
	return report.Diagnostics(os.Stdout, value)
}
