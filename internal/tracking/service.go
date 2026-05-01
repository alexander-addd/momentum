package tracking

import (
	"context"
	"time"
)

type Service struct {
	clock Clock
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

func NewService(clock Clock) *Service {
	return &Service{clock: clock}
}

func (t *Service) Start(ctx context.Context, input StartInput) (Entry, error)
func (t *Service) Stop(ctx context.Context) (Entry, error)
func (t *Service) Status(ctx context.Context) (Status, error)
func (t *Service) Today(ctx context.Context) ([]Entry, error)
func (t *Service) Log(ctx context.Context, limit int) ([]Entry, error)
