package tracking

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alexander-addd/momentum/internal/storage"
	"github.com/google/uuid"
)

type sqlcEntryRow interface {
	storage.GetActiveEntryRow | storage.GetEntryByIDRow | storage.GetEntriesRow
}

// Add this for convenience, to avoid mapping db structs
type entryRow struct {
	ID          string
	Description string
	ProjectID   sql.NullString
	ProjectName sql.NullString
	StartedAt   int64
	StoppedAt   sql.NullInt64
	CreatedAt   int64
	UpdatedAt   int64
}

func toEntryMapper[T sqlcEntryRow](row T) (Entry, error) {
	e := entryRow(row)

	id, err := uuid.Parse(e.ID)
	if err != nil {
		return Entry{}, fmt.Errorf("parse entry id: %w", err)
	}

	var stoppedAt *time.Time
	if e.StoppedAt.Valid {
		t := time.Unix(e.StoppedAt.Int64, 0)
		stoppedAt = &t
	}

	entry := Entry{
		ID:          id,
		Description: e.Description,
		Project:     e.ProjectName.String,
		Tags:        []string{}, // ignoring for now
		StartedAt:   time.Unix(e.StartedAt, 0),
		StoppedAt:   stoppedAt,
		CreatedAt:   time.Unix(e.CreatedAt, 0),
		UpdatedAt:   time.Unix(e.UpdatedAt, 0),
	}

	return entry, nil
}

func toEntriesMapper(rows []storage.GetEntriesRow) ([]Entry, error) {
	entries := make([]Entry, len(rows))

	for i, row := range rows {
		entry, err := toEntryMapper(row)
		if err != nil {
			return nil, fmt.Errorf("map entry: %w", err)
		}

		entries[i] = entry
	}

	return entries, nil
}
