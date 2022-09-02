package migrator

import (
	"time"

	"go.uber.org/zap"
)

type Option interface {
	apply(*migrator)
}

// WithLogger option set configured logger.
func WithLogger(logger *zap.Logger) Option {
	return &withLoggerOption{logger: logger}
}

type withLoggerOption struct {
	logger *zap.Logger
}

func (o *withLoggerOption) apply(m *migrator) {
	m.logger = o.logger
}

// WithTimeout timeout option for each migration.
func WithTimeout(timeout time.Duration) Option {
	return &withTimeout{timeout: timeout}
}

type withTimeout struct {
	timeout time.Duration
}

func (o *withTimeout) apply(m *migrator) {
	m.timeout = o.timeout
}
