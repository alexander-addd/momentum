package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alexander-addd/momentum/internal/tracking"
	"github.com/google/uuid"
)

func TestRunSuccessfulCommands(t *testing.T) {
	stopStatus := tracking.Status{
		Entry:   tracking.Entry{Description: "Write README"},
		Elapsed: 90 * time.Minute,
	}
	activeStatus := tracking.Status{
		Active:  true,
		Entry:   tracking.Entry{Description: "Write README"},
		Elapsed: 90 * time.Minute,
	}

	tests := []struct {
		name       string
		args       []string
		service    fakeTracker
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
			wantStdout: "New entry started: Write README, project: Empty\n",
		},
		{
			name:       "start with project",
			args:       []string{"start", "Write README", "--project", "momentum"},
			wantStdout: "New entry started: Write README, project: momentum\n",
		},
		{
			name:       "stop",
			args:       []string{"stop"},
			service:    fakeTracker{stopStatus: stopStatus},
			wantStdout: "Entry \"Write README\" stopped\nTime elapsed: 1h30m0s\n",
		},
		{
			name:       "status inactive",
			args:       []string{"status"},
			wantStdout: "No active entry\n",
		},
		{
			name:       "status active",
			args:       []string{"status"},
			service:    fakeTracker{status: activeStatus},
			wantStdout: "Entry \"Write README\" is active\nTime elapsed: 1h30m0s\n",
		},
		{
			name:       "today",
			args:       []string{"today"},
			wantStdout: "No entries tracked today\n",
		},
		{
			name:       "log empty",
			args:       []string{"log"},
			wantStdout: "No recent entries\n",
		},
		{
			name: "log entries",
			args: []string{"log", "2"},
			service: fakeTracker{logEntries: []tracking.Entry{
				stoppedEntry(
					"11111111-1111-1111-1111-111111111111",
					"Write README",
					"momentum",
					time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC),
					90*time.Minute,
				),
				{
					ID:          uuid.MustParse("22222222-2222-2222-2222-222222222222"),
					Description: "Keep going",
					StartedAt:   time.Date(2026, 5, 10, 10, 45, 0, 0, time.UTC),
				},
			}},
			wantStdout: "" +
				"ID        STARTED           DURATION  PROJECT   DESCRIPTION\n" +
				"11111111  2026-05-10 09:00  1h 30m    momentum  Write README\n" +
				"22222222  2026-05-10 10:45  running   -         Keep going\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			gotCode := Run(context.Background(), tt.service, tt.args, &stdout, &stderr)

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
			name:          "tag flag unsupported",
			args:          []string{"start", "Write README", "--tag", "planning"},
			wantErrOutput: "unknown start flag \"--tag\"",
		},
		{
			name:          "repeated project",
			args:          []string{"start", "Write README", "--project", "momentum", "--project", "other"},
			wantErrOutput: "--project can only be provided once",
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
			name:          "log invalid limit",
			args:          []string{"log", "now"},
			wantErrOutput: "log limit must be a positive integer",
		},
		{
			name:          "log zero limit",
			args:          []string{"log", "0"},
			wantErrOutput: "log limit must be a positive integer",
		},
		{
			name:          "log too many args",
			args:          []string{"log", "1", "2"},
			wantErrOutput: "log accepts at most one limit",
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

func stoppedEntry(id, description, project string, startedAt time.Time, duration time.Duration) tracking.Entry {
	stoppedAt := startedAt.Add(duration)
	return tracking.Entry{
		ID:          uuid.MustParse(id),
		Description: description,
		Project:     project,
		StartedAt:   startedAt,
		StoppedAt:   &stoppedAt,
	}
}

type fakeTracker struct {
	stopStatus tracking.Status
	status     tracking.Status
	logEntries []tracking.Entry
}

func (fakeTracker) Start(_ context.Context, input tracking.StartInput) (tracking.Entry, error) {
	return tracking.Entry{
		Description: input.Description,
		Project:     input.Project,
	}, nil
}

func (f fakeTracker) Stop(context.Context) (tracking.Status, error) {
	return f.stopStatus, nil
}

func (f fakeTracker) Status(context.Context) (tracking.Status, error) {
	return f.status, nil
}

func (fakeTracker) Today(context.Context) ([]tracking.Entry, error) {
	return nil, nil
}

func (f fakeTracker) Log(context.Context, tracking.LogInput) ([]tracking.Entry, error) {
	return f.logEntries, nil
}
