package tracking

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alexander-addd/momentum/internal/storage"
	"github.com/google/uuid"
)

type Service struct {
	clock Clock
	store *storage.Queries // maybe repository layer later
}

type Tracker interface {
	Start(context.Context, StartInput) (Entry, error)
	Stop(context.Context) (Status, error)
	Status(context.Context) (Status, error)
	Today(context.Context) ([]Entry, error)
	Log(context.Context, int) ([]Entry, error)
}

type StartInput struct {
	Description string
	Project     string
	Tags        []string // ignoring for now
}

type Status struct {
	Active  bool
	Entry   Entry
	Elapsed time.Duration
}

type resolvedProject struct {
	id   sql.NullString
	name string
}

func NewService(clock Clock, store *storage.Queries) *Service {
	return &Service{clock: clock, store: store}
}

func (s *Service) Start(ctx context.Context, input StartInput) (Entry, error) {
	if strings.TrimSpace(input.Description) == "" {
		return Entry{}, ErrEmptyDescription
	}

	_, err := s.store.GetActiveEntry(ctx)
	if err == nil {
		return Entry{}, ErrActiveEntryExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return Entry{}, fmt.Errorf("get active entry: %w", err)
	}

	now := s.clock.Now()

	project, err := s.resolveProjectID(ctx, input.Project, now)
	if err != nil {
		return Entry{}, fmt.Errorf("resolve project id: %w", err)
	}

	entryID := uuid.New().String()
	entryPayload := storage.CreateEntryParams{
		ID:          entryID,
		Description: input.Description,
		ProjectID:   project.id,
		StartedAt:   now.Unix(),
		CreatedAt:   now.Unix(),
		UpdatedAt:   now.Unix(),
	}

	err = s.store.CreateEntry(ctx, entryPayload)
	if err != nil {
		return Entry{}, fmt.Errorf("create entry: %w", err)
	}

	newEntry, err := s.store.GetEntryByID(ctx, entryID)
	if err != nil {
		return Entry{}, fmt.Errorf("get entry by id: %w", err)
	}

	mappedEntry, err := toEntryMapper(newEntry)
	if err != nil {
		return Entry{}, fmt.Errorf("map to entry: %w", err)
	}

	return mappedEntry, nil
}

func (s *Service) Stop(ctx context.Context) (Status, error) {
	entry, err := s.store.GetActiveEntry(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return Status{}, ErrNoActiveEntry
	}
	if err != nil {
		return Status{}, fmt.Errorf("get active entry: %w", err)
	}

	now := s.clock.Now()
	entryStoppedAt := sql.NullInt64{Int64: now.Unix(), Valid: true}
	entryUpdatedAt := now.Unix()

	rowsAffected, err := s.store.StopActiveEntry(ctx, storage.StopActiveEntryParams{
		ID:        entry.ID,
		StoppedAt: entryStoppedAt,
		UpdatedAt: entryUpdatedAt,
	})
	if err != nil {
		return Status{}, fmt.Errorf("stop active entry: %w", err)
	}
	if rowsAffected == 0 {
		return Status{}, ErrNoActiveEntry
	}

	entry.StoppedAt = entryStoppedAt
	entry.UpdatedAt = entryUpdatedAt

	mappedEntry, err := toEntryMapper(entry)
	if err != nil {
		return Status{}, fmt.Errorf("map to entry: %w", err)
	}

	status := Status{
		Entry:   mappedEntry,
		Active:  mappedEntry.Active(),
		Elapsed: mappedEntry.Duration(now),
	}

	return status, nil
}

func (s *Service) Status(ctx context.Context) (Status, error) {
	return Status{}, nil
}

func (s *Service) Today(ctx context.Context) ([]Entry, error) {
	return []Entry{}, nil
}

func (s *Service) Log(ctx context.Context, limit int) ([]Entry, error) {
	return []Entry{}, nil
}

func (s *Service) resolveProjectID(ctx context.Context, name string, now time.Time) (resolvedProject, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return resolvedProject{id: sql.NullString{}, name: ""}, nil
	}

	project, err := s.store.GetProjectByName(ctx, name)
	switch {
	case err == nil:
		return resolvedProject{
			id:   sql.NullString{String: project.ID, Valid: true},
			name: project.Name,
		}, nil

	case errors.Is(err, sql.ErrNoRows):
		id := uuid.New()

		err := s.store.CreateProject(ctx, storage.CreateProjectParams{
			ID:        id.String(),
			Name:      name,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
		})
		if err != nil {
			return resolvedProject{id: sql.NullString{}, name: ""}, fmt.Errorf("create project: %w", err)
		}

		return resolvedProject{
			id:   sql.NullString{String: id.String(), Valid: true},
			name: name}, nil

	default:
		return resolvedProject{id: sql.NullString{}, name: ""}, fmt.Errorf("get project by name: %w", err)
	}
}
