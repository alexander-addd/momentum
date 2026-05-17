package cli

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alexander-addd/momentum/internal/tracking"
)

type runner struct {
	service tracking.Tracker
	stdout  io.Writer
	stderr  io.Writer
}

const defaultLogLimit = 10

func (r runner) start(ctx context.Context, args []string) int {
	input, err := parseStart(args)
	if err != nil {
		return usageError(r.stderr, err.Error())
	}

	entry, err := r.service.Start(ctx, input)
	if err != nil {
		return serviceError(r.stderr, "couldn't start timer: %v", err)
	}

	fmt.Fprintln(r.stdout, entry)

	fmt.Fprintf(r.stdout, "New entry started: %s, project: %s\n", entry.Description, getTitleIfEmpty(entry.Project))
	return exitOK
}

func (r runner) log(ctx context.Context, args []string) int {
	input := parseLog(args)

	entries, err := r.service.Log(ctx, input)
	if err != nil {
		return serviceError(r.stderr, "error getting log: %v", err)
	}

	return exitOK
}

func (r runner) stop(ctx context.Context) int {
	status, err := r.service.Stop(ctx)
	if err != nil {
		return serviceError(r.stderr, "couldn't stop timer: %v", err)
	}

	fmt.Fprintf(r.stdout, "Entry %q stopped\n", status.Entry.Description)
	fmt.Fprintf(r.stdout, "Time elapsed: %s\n", status.Elapsed)

	return exitOK
}

func (r runner) status(ctx context.Context) int {
	status, err := r.service.Status(ctx)
	if err != nil {
		return serviceError(r.stderr, "couldn't check status: %v", err)
	}

	if !status.Active {
		fmt.Fprintln(r.stdout, "No active entry")
		return exitOK
	}

	fmt.Fprintf(r.stdout, "Entry %q is active\n", status.Entry.Description)
	fmt.Fprintf(r.stdout, "Time elapsed: %s\n", status.Elapsed)

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

func parseLog(args []string) tracking.LogInput {
	if len(args) != 1 {
		return tracking.LogInput{Limit: defaultLogLimit}
	}

	limit, err := strconv.Atoi(args[0])
	if err != nil {
		return tracking.LogInput{Limit: defaultLogLimit}
	}

	return tracking.LogInput{Limit: limit}
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
