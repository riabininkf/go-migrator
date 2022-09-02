package migrator_test

import (
	"testing"

	"github.com/riabininkf/go-migrator/pkg/generator"
	"github.com/riabininkf/go-migrator/pkg/migration"
	"github.com/riabininkf/go-migrator/pkg/migrator"
	"github.com/riabininkf/go-migrator/pkg/registry"
	"github.com/riabininkf/go-migrator/pkg/scanner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMigrator_Generate(t *testing.T) {
	var (
		name = "test_name"
		path = "./"
	)

	g := new(generator.MockGenerator)
	g.On("Generate", mock.AnythingOfType("string"), path).Return(nil)

	m := migrator.New(new(scanner.MockScanner), g, new(registry.MockRegistry))
	assert.NoError(t, m.Generate(name, path))
}

func TestMigrator_Up(t *testing.T) {
	t.Run("positive case", func(t *testing.T) {
		path := "./"

		m1 := new(migration.MockMigration)
		m1.On("Name").Return("test_migration_name")
		m1.On("Up").Return(nil)

		s := new(scanner.MockScanner)
		s.On("Scan", path).Return([]migration.Migration{m1}, nil)

		r := new(registry.MockRegistry)
		r.On("All", mock.AnythingOfType("*context.timerCtx")).Return(nil, nil)

		registration := new(registry.MockRegistration)
		registration.On("ID").Return(uint(1))

		r.On("Up", mock.AnythingOfType("*context.timerCtx"), "test_migration_name").Return(registration, nil)
		r.On("Finish", mock.AnythingOfType("*context.timerCtx"), uint(1)).Return(nil)

		m := migrator.New(s, new(generator.MockGenerator), r)

		assert.NoError(t, m.Up(path))
	})

	t.Run("no migrations found", func(t *testing.T) {
		path := "./"

		m1 := new(migration.MockMigration)
		m1.On("Name").Return("test_migration_name")
		m1.On("Up").Return(nil)

		s := new(scanner.MockScanner)
		s.On("Scan", path).Return([]migration.Migration{}, nil)

		m := migrator.New(s, new(generator.MockGenerator), new(registry.MockRegistry))
		assert.NoError(t, m.Up(path))
	})
}

func TestMigrator_Down(t *testing.T) {
	t.Run("positive case", func(t *testing.T) {
		var (
			path = "./"
			name = "test_migration"
		)

		m1 := new(migration.MockMigration)
		m1.On("Name").Return(name)
		m1.On("Down").Return(nil)
		s := new(scanner.MockScanner)
		s.On("Scan", path).Return([]migration.Migration{m1}, nil)

		regUp := new(registry.MockRegistration)
		regUp.On("IsUp").Return(true)
		regUp.On("Name").Return(name)
		r := new(registry.MockRegistry)
		r.On("All", mock.AnythingOfType("*context.timerCtx")).
			Return([]registry.Registration{regUp}, nil)

		regDown := new(registry.MockRegistration)
		regDown.On("ID").Return(uint(2))

		r.On("Down", mock.AnythingOfType("*context.timerCtx"), name).Return(regDown, nil)
		r.On("Finish", mock.AnythingOfType("*context.timerCtx"), uint(2)).Return(nil)
		m := migrator.New(s, new(generator.MockGenerator), r)

		assert.NoError(t, m.Down(path))
	})

	t.Run("no registrations to rollback", func(t *testing.T) {
		path := "./"

		r := new(registry.MockRegistry)
		r.On("All", mock.AnythingOfType("*context.timerCtx")).Return(nil, nil)

		m := migrator.New(new(scanner.MockScanner), new(generator.MockGenerator), r)
		assert.NoError(t, m.Down(path))
	})
}
