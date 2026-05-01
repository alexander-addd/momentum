package cli

import (
	"fmt"
	"io"
	"strings"
)

const (
	exitOK    = 0
	exitError = 1
)

const helpText = `Usage:
  momentum start <description> [--project <name>] [--tag <tag>]
  momentum stop
  momentum status
  momentum today
  momentum log
  momentum help
`

type startOptions struct {
	description string
	project     string
	tag         string
}

func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		printHelp(stdout)
		return exitOK
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "start":
		return runStart(commandArgs, stdout, stderr)
	case "stop":
		return runNoArgCommand(commandArgs, stdout, stderr, "Stopped active timer\n")
	case "status":
		return runNoArgCommand(commandArgs, stdout, stderr, "No active timer\n")
	case "today":
		return runNoArgCommand(commandArgs, stdout, stderr, "No entries tracked today\n")
	case "log":
		return runNoArgCommand(commandArgs, stdout, stderr, "No recent entries\n")
	case "help":
		if len(commandArgs) > 0 {
			return usageError(stderr, "help does not accept arguments")
		}
		printHelp(stdout)
		return exitOK
	default:
		return usageError(stderr, "unknown command %q", command)
	}
}

func runStart(args []string, stdout io.Writer, stderr io.Writer) int {
	options, err := parseStart(args)
	if err != nil {
		return usageError(stderr, err.Error())
	}

	fmt.Fprintf(stdout, "Started project %s with tag %s\n", getTitleIfEmpty(options.project), getTitleIfEmpty(options.tag))
	return exitOK
}

func parseStart(args []string) (startOptions, error) {
	options := startOptions{
		project: "",
		tag:     "",
	}

	var description []string
	seenProject := false
	seenTag := false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "--project":
			value, seen, err := parseFlagValue(args, i, "--project", seenProject)
			if err != nil {
				return startOptions{}, err
			}
			seenProject = seen
			options.project = value
			i++

		case arg == "--tag":
			value, seen, err := parseFlagValue(args, i, "--tag", seenTag)
			if err != nil {
				return startOptions{}, err
			}
			seenTag = seen
			options.tag = value
			i++

		case strings.HasPrefix(arg, "--"):
			return startOptions{}, fmt.Errorf("unknown start flag %q", arg)

		default:
			description = append(description, arg)
		}
	}

	if len(description) == 0 {
		return startOptions{}, fmt.Errorf("start requires a description")
	}

	options.description = strings.Join(description, " ")
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

func runNoArgCommand(args []string, stdout io.Writer, stderr io.Writer, message string) int {
	if len(args) > 0 {
		return usageError(stderr, "command does not accept arguments")
	}

	fmt.Fprint(stdout, message)
	return exitOK
}

func usageError(stderr io.Writer, format string, args ...any) int {
	fmt.Fprintf(stderr, "Error: "+format+"\n\n", args...)
	printHelp(stderr)
	return exitError
}

func printHelp(w io.Writer) {
	fmt.Fprint(w, helpText)
}

func getTitleIfEmpty(title string) string {
	if title == "" {
		return "Empty"
	}
	return title
}
