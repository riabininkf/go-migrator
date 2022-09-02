package registry

import (
	"context"
	"database/sql"
	"fmt"
)

const DefaultTableName = "migrations"

type postgres struct {
	tableName string
	db        *sql.DB
}

func NewPostgres(tableName string, db *sql.DB) (Registry, error) {
	p := &postgres{
		db:        db,
		tableName: tableName,
	}

	if err := p.init(); err != nil {
		return nil, fmt.Errorf("can't init registry: %w", err)
	}

	return p, nil
}

func (p *postgres) init() error {
	query := fmt.Sprintf(`
create table if not exists %s
(
	id 		serial,
    name    text                    not null,
 	type    text                    not null,
    status  text                    not null,
    updated timestamp default now() not null
);
`, p.tableName)

	if _, err := p.db.Exec(query); err != nil {
		return fmt.Errorf("can't create table for migrations: %w", err)
	}

	return nil
}

func (p *postgres) All(ctx context.Context) ([]Registration, error) {
	var err error

	query := fmt.Sprintf(`
select distinct on (name) id,
       name,
       type,
       status,
       updated,
       (
           select type
           from %s
           where name = m.name
           ORDER BY id DESC
           LIMIT 1
       ) as last_state
from %s m
order by name;
`, p.tableName, p.tableName)

	var rows *sql.Rows
	if rows, err = p.db.QueryContext(ctx, query); err != nil {
		return nil, fmt.Errorf("can't select registrations: %w", err)
	}
	defer rows.Close()

	registrations := make([]Registration, 0)

	for rows.Next() {
		r := &registration{}
		if err := rows.Scan(&r.id, &r.name, &r.migrationType, &r.status, &r.updated, &r.lastState); err != nil {
			return nil, fmt.Errorf("can't scan database row: %w", err)
		}

		registrations = append(registrations, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error due processing rows: %w", err)
	}

	return registrations, nil
}

func (p *postgres) Up(ctx context.Context, name string) (Registration, error) {
	return p.add(ctx, name, TypeUp)
}

func (p *postgres) Down(ctx context.Context, name string) (Registration, error) {
	return p.add(ctx, name, TypeDown)
}

func (p *postgres) Process(ctx context.Context, id uint) error {
	return p.updateStatus(ctx, id, StatusFinished)
}

func (p *postgres) Finish(ctx context.Context, id uint) error {
	return p.updateStatus(ctx, id, StatusFinished)
}

func (p *postgres) Fail(ctx context.Context, id uint) error {
	return p.updateStatus(ctx, id, StatusFailed)
}

func (p *postgres) add(ctx context.Context, name string, migrationType string) (Registration, error) {
	query := fmt.Sprintf(`
insert into %s (name, type, status)
values ($1, $2, $3)
RETURNING id, name, type, status, updated;
`, p.tableName)

	r := &registration{lastState: migrationType}
	if err := p.db.QueryRowContext(ctx, query, name, migrationType, StatusNew).Scan(
		&r.id, &r.name, &r.migrationType, &r.status, &r.updated,
	); err != nil {
		return nil, fmt.Errorf("can't insert new registration: %w", err)
	}

	return r, nil
}

func (p *postgres) updateStatus(ctx context.Context, id uint, status string) error {
	query := fmt.Sprintf("update %s set status = $1, updated = NOW() where id = $2", p.tableName)

	if _, err := p.db.ExecContext(ctx, query, status, id); err != nil {
		return fmt.Errorf("can't insert new registration: %w", err)
	}

	return nil
}
