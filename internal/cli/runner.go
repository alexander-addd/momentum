package cli

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

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

	fmt.Fprintf(r.stdout, "New entry started: %s, project: %s\n", entry.Description, getTitleIfEmpty(entry.Project))
	return exitOK
}

func (r runner) log(ctx context.Context, args []string) int {
	input, err := parseLog(args)
	if err != nil {
		return usageError(r.stderr, err.Error())
	}

	entries, err := r.service.Log(ctx, input)
	if err != nil {
		return serviceError(r.stderr, "error getting log: %v", err)
	}

	printEntries(r.stdout, entries)
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

func parseLog(args []string) (tracking.LogInput, error) {
	if len(args) == 0 {
		return tracking.LogInput{Limit: defaultLogLimit}, nil
	}
	if len(args) > 1 {
		return tracking.LogInput{}, fmt.Errorf("log accepts at most one limit")
	}

	limit, err := strconv.Atoi(args[0])
	if err != nil || limit <= 0 {
		return tracking.LogInput{}, fmt.Errorf("log limit must be a positive integer")
	}

	return tracking.LogInput{Limit: limit}, nil
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

// Formatter
// Might make sense to extract it
func printEntries(w io.Writer, entries []tracking.Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No recent entries")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 4, ' ', 0)
	fmt.Fprintln(tw, "ID\tSTARTED\tDURATION\tPROJECT\tDESCRIPTION")

	for _, entry := range entries {
		fmt.Fprintf(
			tw,
			"%s\t%s\t%s\t%s\t%s\n",
			shortEntryID(entry),
			formatEntryTime(entry.StartedAt),
			formatEntryDuration(entry),
			formatProject(entry.Project),
			entry.Description,
		)
	}

	_ = tw.Flush()
}

func shortEntryID(entry tracking.Entry) string {
	id := entry.ID.String()
	if len(id) <= 8 {
		return id
	}
	return id[:8]
}

func formatEntryTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func formatEntryDuration(entry tracking.Entry) string {
	if entry.StoppedAt == nil {
		return "running"
	}

	return formatDuration(entry.StoppedAt.Sub(entry.StartedAt))
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		return "0s"
	}

	d = d.Round(time.Second)
	hours := int(d / time.Hour)
	d -= time.Duration(hours) * time.Hour
	minutes := int(d / time.Minute)
	d -= time.Duration(minutes) * time.Minute
	seconds := int(d / time.Second)

	switch {
	case hours > 0 && minutes > 0:
		return fmt.Sprintf("%dh %dm", hours, minutes)
	case hours > 0:
		return fmt.Sprintf("%dh", hours)
	case minutes > 0 && seconds > 0:
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	case minutes > 0:
		return fmt.Sprintf("%dm", minutes)
	default:
		return fmt.Sprintf("%ds", seconds)
	}
}

func formatProject(project string) string {
	if project == "" {
		return "-"
	}
	return project
}
