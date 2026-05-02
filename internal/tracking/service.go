package tracking

import (
	"context"
	"time"

	"github.com/alexander-addd/momentum/internal/storage"
)

type Service struct {
	clock Clock
	store *storage.Queries // maybe repository layer later
}

type Tracker interface {
	Start(context.Context, StartInput) (Entry, error)
	Stop(context.Context) (Entry, error)
	Status(context.Context) (Status, error)
	Today(context.Context) ([]Entry, error)
	Log(context.Context, int) ([]Entry, error)
}

type StartInput struct {
	Description string
	Project     string
	Tags        []string
}

type Status struct {
	Active  bool
	Entry   Entry
	Elapsed time.Duration
}

func NewService(clock Clock, store *storage.Queries) *Service {
	return &Service{clock: clock, store: store}
}

func (s *Service) Start(ctx context.Context, input StartInput) (Entry, error) {
	return Entry{}, nil
}

func (s *Service) Stop(ctx context.Context) (Entry, error) {
	return Entry{}, nil
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
