package migrator

//go:generate mockery -case=underscore -output=. -inpkg -name=Migrator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/riabininkf/go-migrator/pkg/generator"
	"github.com/riabininkf/go-migrator/pkg/migration"
	"github.com/riabininkf/go-migrator/pkg/registry"
	"github.com/riabininkf/go-migrator/pkg/scanner"
	"go.uber.org/zap"
)

var (
	defaultLogger  = zap.NewNop()
	defaultTimeout = time.Second * 10
)

type Migrator interface {
	Generate(name string, path string) error
	Up(path string) error
	Down(path string) error
	Redo(path string) error
	Status() error
	Version() error
}

type migrator struct {
	scanner   scanner.Scanner
	generator generator.Generator
	registry  registry.Registry
	logger    *zap.Logger
	timeout   time.Duration
}

func New(scanner scanner.Scanner, generator generator.Generator, registry registry.Registry, opts ...Option) Migrator {
	m := &migrator{
		scanner:   scanner,
		generator: generator,
		registry:  registry,
		logger:    defaultLogger,
		timeout:   defaultTimeout,
	}

	for _, opt := range opts {
		opt.apply(m)
	}

	return m
}

func (m *migrator) Up(path string) error {
	var err error

	var migrations []migration.Migration
	if migrations, err = m.getNewMigrations(path); err != nil {
		return fmt.Errorf("can't get new migrations: %w", err)
	}

	if len(migrations) == 0 {
		m.logger.Info("no new migrations found")
		return nil
	}

	for _, migration := range migrations {
		m.logger.Info("applying migration", zap.String("migration", migration.Name()))
		if err = m.upWithRegistration(migration); err != nil {
			return fmt.Errorf("can't apply migration: %w", err)
		}

		m.logger.Info("migration applied successfully", zap.String("migration", migration.Name()))
	}

	return nil
}
func (m *migrator) Generate(name string, path string) error {
	version := time.Now().UTC().Format("2006_01_02_15_04_05")
	fullName := fmt.Sprintf("%s_%s", version, strings.ReplaceAll(name, " ", "_"))

	return m.generator.Generate(fullName, path)
}

func (m *migrator) Down(path string) error {
	var err error

	var last migration.Migration
	if last, err = m.getLastAppliedMigration(path); err != nil {
		return fmt.Errorf("can't get last applied migration: %w", err)
	}

	if last == nil {
		m.logger.Info("no migrations to roll back")
		return nil
	}

	m.logger.Info("rolling back migration", zap.String("migration", last.Name()))
	if err = m.downWithRegistration(last); err != nil {
		return fmt.Errorf("can't rollback migration: %w", err)
	}

	m.logger.Info("migration rolled back successfully", zap.String("migration", last.Name()))
	return nil
}

func (m *migrator) Redo(path string) error {
	var err error

	var last migration.Migration
	if last, err = m.getLastAppliedMigration(path); err != nil {
		return fmt.Errorf("can't get last applied migration: %w", err)
	}

	if last == nil {
		m.logger.Info("no migrations to redo")
		return nil
	}

	if err = m.downWithRegistration(last); err != nil {
		return fmt.Errorf("can't roll back migration: %w", err)
	}

	if err = m.upWithRegistration(last); err != nil {
		return fmt.Errorf("can't apply migration: %w", err)
	}

	m.logger.Info("migration reapplied successfully", zap.String("migration", last.Name()))
	return nil
}

func (m *migrator) Status() error {
	var err error

	ctx, cancelFunc := context.WithTimeout(context.Background(), m.timeout)
	defer cancelFunc()

	var registrations []registry.Registration
	if registrations, err = m.registry.All(ctx); err != nil {
		return fmt.Errorf("can't get migrations from registry: %w", err)
	}

	if len(registrations) == 0 {
		m.logger.Info("no migrations found")
		return nil
	}

	for _, r := range registrations {
		if r.IsUp() {
			m.logger.Info(
				"migration status",
				zap.String("name", r.Name()),
				zap.String("status", r.Status()),
			)
		}
	}

	return nil
}

func (m *migrator) Version() error {
	var err error

	var last registry.Registration
	if last, err = m.getLastAppliedRegistration(); err != nil {
		return fmt.Errorf("can't get last registrations: %w", err)
	}

	if last == nil {
		m.logger.Info("no migrations found")
		return nil
	}

	m.logger.Info("database version", zap.Uint("last_id", last.ID()), zap.String("last_name", last.Name()))
	return nil
}

func (m *migrator) getNewMigrations(path string) ([]migration.Migration, error) {
	var err error

	var migrations []migration.Migration
	if migrations, err = m.scanner.Scan(path); err != nil {
		return nil, fmt.Errorf("can't get migrations from provided directory: %w", err)
	}

	if len(migrations) == 0 {
		return migrations, nil
	}

	var last registry.Registration
	if last, err = m.getLastAppliedRegistration(); err != nil {
		return nil, fmt.Errorf("can't get last applied migration: %w", err)
	}

	if last == nil {
		return migrations, nil
	}

	if migrations[len(migrations)-1].Name() == last.Name() {
		return make([]migration.Migration, 0), nil
	}

	for i, migration := range migrations {
		if migration.Name() == last.Name() {
			return migrations[i+1:], nil
		}
	}

	return make([]migration.Migration, 0), nil
}

func (m *migrator) getLastAppliedMigration(path string) (migration.Migration, error) {
	var err error

	var last registry.Registration
	if last, err = m.getLastAppliedRegistration(); err != nil {
		return nil, fmt.Errorf("can't get last applied registration: %w", err)
	}

	if last == nil {
		return nil, nil
	}

	var migrations []migration.Migration
	if migrations, err = m.scanner.Scan(path); err != nil {
		return nil, fmt.Errorf("can't get migrations from %s: %w", path, err)
	}

	for _, migration := range migrations {
		if migration.Name() == last.Name() {
			return migration, nil
		}
	}

	return nil, nil
}

func (m *migrator) getLastAppliedRegistration() (registry.Registration, error) {
	var err error

	ctx, cancelFunc := context.WithTimeout(context.Background(), m.timeout)
	defer cancelFunc()

	var registrations []registry.Registration
	if registrations, err = m.registry.All(ctx); err != nil {
		return nil, fmt.Errorf("can't get registrations: %w", err)
	}

	for i := len(registrations) - 1; i >= 0; i-- {
		if registrations[i].IsUp() {
			return registrations[i], nil
		}
	}

	return nil, nil
}

func (m *migrator) upWithRegistration(migration migration.Migration) error {
	var err error

	ctx, cancelFunc := context.WithTimeout(context.Background(), m.timeout)
	defer cancelFunc()

	var registration registry.Registration
	if registration, err = m.registry.Up(ctx, migration.Name()); err != nil {
		return fmt.Errorf("can't register migration: %w", err)
	}

	if err := migration.Up(); err != nil {
		if registryError := m.registry.Fail(ctx, registration.ID()); registryError != nil {
			m.logger.Error("can't register migration failure", zap.Error(err))
		}

		return fmt.Errorf("can't apply migration: %w", err)
	}

	if err := m.registry.Finish(ctx, registration.ID()); err != nil {
		return fmt.Errorf("can't mark migration as finished: %w", err)
	}

	return nil
}

func (m *migrator) downWithRegistration(migration migration.Migration) error {
	var err error

	ctx, cancelFunc := context.WithTimeout(context.Background(), m.timeout)
	defer cancelFunc()

	var registration registry.Registration
	if registration, err = m.registry.Down(ctx, migration.Name()); err != nil {
		return fmt.Errorf("can't register migration: %w", err)
	}

	if err := migration.Down(); err != nil {
		if registryError := m.registry.Fail(ctx, registration.ID()); registryError != nil {
			m.logger.Error("can't register migration failure", zap.Error(err))
		}

		return fmt.Errorf("can't apply migration: %w", err)
	}

	if err := m.registry.Finish(ctx, registration.ID()); err != nil {
		return fmt.Errorf("can't mark migration as finished: %w", err)
	}

	return nil
}
