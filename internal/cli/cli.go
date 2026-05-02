package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/alexander-addd/momentum/internal/tracking"
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

func Run(ctx context.Context, service tracking.Tracker, args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		printHelp(stdout)
		return exitOK
	}

	command := args[0]
	commandArgs := args[1:]
	r := runner{service: service, stdout: stdout, stderr: stderr}

	switch command {
	case "start":
		return r.runStart(ctx, commandArgs)
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

func serviceError(stderr io.Writer, format string, args ...any) int {
	fmt.Fprintf(stderr, "Error: "+format+"\n\n", args...)
	return exitError
}

func printHelp(w io.Writer) {
	fmt.Fprint(w, helpText)
}
