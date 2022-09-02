package registry

//go:generate mockery -case=underscore -output=. -inpkg -name=Registry

import "context"

const (
	StatusNew      = "new"
	StatusFailed   = "failed"
	StatusFinished = "finished"
)

const (
	TypeUp   = "up"
	TypeDown = "down"
)

type Registry interface {
	All(ctx context.Context) ([]Registration, error)
	Up(ctx context.Context, name string) (Registration, error)
	Down(ctx context.Context, name string) (Registration, error)
	Process(ctx context.Context, id uint) error
	Finish(ctx context.Context, id uint) error
	Fail(ctx context.Context, id uint) error
}
