package command

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

type Args struct {
	Command                                                                                                                string
	Home, ID, Target, Instrument, ObservedAt, Description, Reviewer, Decision, Note, Publisher, State, Tag, Path, Contains string
	Tags                                                                                                                   []string
}

func Parse(argv []string) (Args, error) {
	if len(argv) < 1 {
		return Args{}, fmt.Errorf("command is required")
	}
	a := Args{Command: argv[0]}
	fs := flag.NewFlagSet(a.Command, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&a.Home, "home", "", "data home")
	fs.StringVar(&a.ID, "id", "", "observation id")
	fs.StringVar(&a.Target, "target", "", "target")
	fs.StringVar(&a.Instrument, "instrument", "", "instrument")
	fs.StringVar(&a.ObservedAt, "observed-at", "", "time")
	fs.StringVar(&a.Description, "description", "", "description")
	fs.StringVar(&a.Reviewer, "reviewer", "", "reviewer")
	fs.StringVar(&a.Decision, "decision", "", "decision")
	fs.StringVar(&a.Note, "note", "", "note")
	fs.StringVar(&a.Publisher, "publisher", "", "publisher")
	fs.StringVar(&a.State, "state", "", "state")
	fs.StringVar(&a.Tag, "tag", "", "tag")
	fs.StringVar(&a.Path, "path", "", "archive path")
	fs.StringVar(&a.Contains, "contains", "", "description text")
	tags := fs.String("tags", "", "tags")
	if err := fs.Parse(argv[1:]); err != nil {
		return Args{}, err
	}
	a.Tags = splitCSV(*tags)
	return a, nil
}

func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if v := strings.TrimSpace(part); v != "" {
			out = append(out, v)
		}
	}
	return out
}
