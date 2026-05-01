package tracking

import (
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	ID          uuid.UUID
	Description string
	Project     string
	Tags        []string
	StartedAt   time.Time
	StoppedAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (e Entry) Active() bool {
	return e.StoppedAt == nil
}

func (e Entry) Duration(now time.Time) time.Duration {
	if e.StoppedAt != nil {
		return e.StoppedAt.Sub(e.StartedAt)
	}
	return now.Sub(e.StartedAt)
}
