package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/alexander-addd/momentum/internal/tracking"
)

func TestRunSuccessfulCommands(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantStdout string
	}{
		{
			name:       "no args prints help",
			args:       nil,
			wantStdout: helpText,
		},
		{
			name:       "help prints help",
			args:       []string{"help"},
			wantStdout: helpText,
		},
		{
			name:       "start with defaults",
			args:       []string{"start", "Write README"},
			wantStdout: "Started project Empty with tag Empty\n",
		},
		{
			name:       "start with project and tag after description",
			args:       []string{"start", "Write README", "--project", "momentum", "--tag", "planning"},
			wantStdout: "Started project momentum with tag planning\n",
		},
		{
			name:       "start with project and tag before description",
			args:       []string{"start", "--project", "momentum", "--tag", "planning", "Write README"},
			wantStdout: "Started project momentum with tag planning\n",
		},
		{
			name:       "stop",
			args:       []string{"stop"},
			wantStdout: "Stopped active timer\n",
		},
		{
			name:       "status",
			args:       []string{"status"},
			wantStdout: "No active timer\n",
		},
		{
			name:       "today",
			args:       []string{"today"},
			wantStdout: "No entries tracked today\n",
		},
		{
			name:       "log",
			args:       []string{"log"},
			wantStdout: "No recent entries\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			gotCode := Run(context.Background(), fakeTracker{}, tt.args, &stdout, &stderr)

			if gotCode != 0 {
				t.Fatalf("Run() exit code = %d, want 0", gotCode)
			}
			if got := stdout.String(); got != tt.wantStdout {
				t.Fatalf("stdout = %q, want %q", got, tt.wantStdout)
			}
			if got := stderr.String(); got != "" {
				t.Fatalf("stderr = %q, want empty", got)
			}
		})
	}
}

func TestRunStartValidationErrors(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		wantErrOutput string
	}{
		{
			name:          "missing description",
			args:          []string{"start"},
			wantErrOutput: "start requires a description",
		},
		{
			name:          "unknown flag",
			args:          []string{"start", "Write README", "--client", "acme"},
			wantErrOutput: "unknown start flag \"--client\"",
		},
		{
			name:          "missing project value",
			args:          []string{"start", "Write README", "--project"},
			wantErrOutput: "--project requires a value",
		},
		{
			name:          "equals project flag",
			args:          []string{"start", "Write README", "--project="},
			wantErrOutput: "unknown start flag \"--project=\"",
		},
		{
			name:          "missing tag value",
			args:          []string{"start", "Write README", "--tag"},
			wantErrOutput: "--tag requires a value",
		},
		{
			name:          "equals tag flag",
			args:          []string{"start", "Write README", "--tag="},
			wantErrOutput: "unknown start flag \"--tag=\"",
		},
		{
			name:          "repeated project",
			args:          []string{"start", "Write README", "--project", "momentum", "--project", "other"},
			wantErrOutput: "--project can only be provided once",
		},
		{
			name:          "repeated tag",
			args:          []string{"start", "Write README", "--tag", "planning", "--tag", "writing"},
			wantErrOutput: "--tag can only be provided once",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertRunError(t, tt.args, tt.wantErrOutput)
		})
	}
}

func TestRunCommandValidationErrors(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		wantErrOutput string
	}{
		{
			name:          "unknown command",
			args:          []string{"dance"},
			wantErrOutput: "unknown command \"dance\"",
		},
		{
			name:          "help extra args",
			args:          []string{"help", "start"},
			wantErrOutput: "help does not accept arguments",
		},
		{
			name:          "stop extra args",
			args:          []string{"stop", "now"},
			wantErrOutput: "command does not accept arguments",
		},
		{
			name:          "status extra args",
			args:          []string{"status", "now"},
			wantErrOutput: "command does not accept arguments",
		},
		{
			name:          "today extra args",
			args:          []string{"today", "now"},
			wantErrOutput: "command does not accept arguments",
		},
		{
			name:          "log extra args",
			args:          []string{"log", "now"},
			wantErrOutput: "command does not accept arguments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertRunError(t, tt.args, tt.wantErrOutput)
		})
	}
}

func assertRunError(t *testing.T, args []string, wantErrOutput string) {
	t.Helper()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	gotCode := Run(context.Background(), fakeTracker{}, args, &stdout, &stderr)

	if gotCode == 0 {
		t.Fatalf("Run() exit code = 0, want non-zero")
	}
	if got := stdout.String(); got != "" {
		t.Fatalf("stdout = %q, want empty", got)
	}

	gotStderr := stderr.String()
	if !strings.Contains(gotStderr, "Error: "+wantErrOutput) {
		t.Fatalf("stderr = %q, want error containing %q", gotStderr, wantErrOutput)
	}
	if !strings.Contains(gotStderr, helpText) {
		t.Fatalf("stderr = %q, want help text", gotStderr)
	}
}

type fakeTracker struct{}

func (fakeTracker) Start(context.Context, tracking.StartInput) (tracking.Entry, error) {
	return tracking.Entry{}, nil
}

func (fakeTracker) Stop(context.Context) (tracking.Entry, error) {
	return tracking.Entry{}, nil
}

func (fakeTracker) Status(context.Context) (tracking.Status, error) {
	return tracking.Status{}, nil
}

func (fakeTracker) Today(context.Context) ([]tracking.Entry, error) {
	return nil, nil
}

func (fakeTracker) Log(context.Context, int) ([]tracking.Entry, error) {
	return nil, nil
}
