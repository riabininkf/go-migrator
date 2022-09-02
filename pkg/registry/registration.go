package registry

//go:generate mockery -case=underscore -output=. -inpkg -name=Registration

import (
	"time"
)

type (
	Registration interface {
		ID() uint
		Name() string
		Type() string
		Status() string
		Updated() time.Time
		IsUp() bool
	}

	registration struct {
		id            uint
		name          string
		migrationType string
		lastState     string
		status        string
		updated       time.Time
	}
)

func (r *registration) ID() uint {
	return r.id
}

func (r *registration) Name() string {
	return r.name
}

func (r *registration) Type() string {
	return r.migrationType
}

func (r *registration) Status() string {
	return r.status
}

func (r *registration) Updated() time.Time {
	return r.updated
}

func (r *registration) IsUp() bool {
	return r.lastState == TypeUp
}
