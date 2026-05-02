package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/alexander-addd/momentum/internal/tracking"
)

type runner struct {
	service tracking.Tracker
	stdout  io.Writer
	stderr  io.Writer
}

func (r runner) runStart(ctx context.Context, args []string) int {
	input, err := parseStart(args)
	if err != nil {
		return usageError(r.stderr, err.Error())
	}

	entry, err := r.service.Start(ctx, input)
	if err != nil {
		return serviceError(r.stderr, "error from service")
	}

	fmt.Fprintln(r.stdout, entry)

	fmt.Fprintf(r.stdout, "Started project %s\n", getTitleIfEmpty(input.Project))
	return exitOK
}

func parseStart(args []string) (tracking.StartInput, error) {
	options := tracking.StartInput{
		Project: "",
		Tags:    []string{},
	}

	var description []string
	seenProject := false
	// seenTag := false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		// TODO: parse tags
		switch {
		case arg == "--project":
			value, seen, err := parseFlagValue(args, i, "--project", seenProject)
			if err != nil {
				return tracking.StartInput{}, err
			}
			seenProject = seen
			options.Project = value
			i++

		case strings.HasPrefix(arg, "--"):
			return tracking.StartInput{}, fmt.Errorf("unknown start flag %q", arg)

		default:
			description = append(description, arg)
		}
	}

	if len(description) == 0 {
		return tracking.StartInput{}, fmt.Errorf("start requires a description")
	}

	options.Description = strings.Join(description, " ")
	return options, nil
}

func parseFlagValue(args []string, i int, flag string, seen bool) (string, bool, error) {
	if seen {
		return "", false, fmt.Errorf("%s can only be provided once", flag)
	}
	if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
		return "", false, fmt.Errorf("%s requires a value", flag)
	}

	return args[i+1], true, nil
}

func getTitleIfEmpty(title string) string {
	if title == "" {
		return "Empty"
	}
	return title
}
