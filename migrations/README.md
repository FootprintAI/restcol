# Migrations

SQL migrations for restcol, named following the [golang-migrate] convention:
`<version>_<slug>.{up,down}.sql`.

## How restcol manages its schema today

The server uses GORM's `AutoMigrate`, and it is **already opt-in**: the sdinsure
storage layer checks the `--restcol_auto_migrate` flag and skips migration when
the flag is false (the default). Production deployments therefore manage schema
changes externally — applying the SQL files in this directory is how.

```
dev          → --restcol_auto_migrate=true, GORM tags drive the schema.
staging/prod → --restcol_auto_migrate=false, apply these files explicitly.
```

## Applying migrations

Any tool that can run ordered SQL files works. Example with [golang-migrate]:

```bash
# one-off install
brew install golang-migrate

# up
migrate -path migrations \
        -database "postgres://postgres:password@localhost:5432/unittest?sslmode=disable" \
        up

# down one step
migrate -path migrations \
        -database "postgres://postgres:password@localhost:5432/unittest?sslmode=disable" \
        down 1
```

`psql` also works for quick one-off application:

```bash
psql "$DATABASE_URL" -f migrations/0002_purge_soft_deleted_documents.up.sql
```

## Writing a new migration

1. Pick the next zero-padded version number.
2. Create both `<N>_<slug>.up.sql` and `<N>_<slug>.down.sql`. A no-op down is
   fine when the migration is strictly additive — write `-- no-op`.
3. Prefer idempotent statements (`IF NOT EXISTS`, `CREATE OR REPLACE`) so
   re-running is safe.
4. If you add a GORM struct tag that AutoMigrate understands, mirror the change
   in a migration file so prod environments pick it up too.

## Future work

Wiring `golang-migrate` into server startup (or a dedicated CLI entry point)
remains a larger refactor: it requires coordinating deploys, proving out a
baseline/snapshot migration matching the current live schema, and deciding how
to handle the existing AutoMigrate state. Track that as a separate project.

[golang-migrate]: https://github.com/golang-migrate/migrate

## What is here, and what is deliberately not

| | |
|---|---|
| `0002_purge_soft_deleted_documents` | erases document rows tombstoned before DeleteDocument became a hard delete (#137). Irreversible; its down migration is a documented no-op |
| `0003_drop_redundant_query_indexes` | removes the indexes `0001` used to create |

**There is no `0001`.** It added query indexes that GORM's `AutoMigrate` already
creates from `index:` struct tags on the models — the same columns under
different names, so its `CREATE INDEX IF NOT EXISTS` never matched and applying
it produced a second index over each. It was deleted rather than fixed: a
migration that re-declares what the models already declare is a second source of
truth for the same schema, and the models win.

Do not re-add index migrations for indexes struct tags can express. Add a
migration when the change is one struct tags *cannot* describe — a backfill, a
data purge, an index built concurrently, a column rewrite.

Numbering is not reused: `0003` follows `0002` even though `0001` is gone.
