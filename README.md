# gomigrator
Design work for golang developer certification
Technical requirements - https://github.com/OtusGolang/final_project/blob/master/04-sql-migrator.md
## Installation
```bash
go get github.com/riabininkf/go-migrator
```
## Base usage
### Generate new migration
```bash
gomigrator create --name="create_table_test" --config="./config.json"
```

### Apply migrations
```bash
gomigrator up --config="./config.json"
```

### Roll back last migration
```bash
gomigrator down --config="./config.json"
```

### Reapply last migration (up + down)
```bash
gomigrator redo --config="./config.json"
```

### Show migrations status
```bash
gomigrator status --config="./config.json"
```

### Show database version (last migration)
```bash
gomigrator dbversion --config="./config.json"
```
## Configuration
### Flags
```bash
-- config - path to config file
-- path - path to directory with migrations
-- db_dsn - database dsn string
```
### Config file
Use `--config` flag or `GOMIGRATOR_CONFIG` env to provide a path to config file.
#### Example
```json
{
  "db": {
    "dsn": "postgres://user:password@localhost:5432/postgres"
  },
  "path": "./migrations"
}
```
#### Environment variables in config
```bash
export DB_DSN="postgres://user:password@localhost:5432/postgres"
export MIGRATIONS_PATH="./migrations"
```
```json
{
  "db": {
    "dsn": "${DB_DSN}"
  },
  "path": "${MIGRATIONS_PATH}"
}
```
