// Package migrations carries restcol's SQL migrations and exposes them as an
// embed.FS so a consuming module can apply exactly the migrations that shipped
// with the version it pins.
//
// WHY THIS FILE EXISTS. These .sql files are plain files in the module, and
// Go's embed cannot reach into a dependency — a consumer that imports restcol
// can see the package but not the files beside it. So a downstream service
// wanting to run these migrations had only one option: copy them. That is what
// FootprintAI/grandturks#987 did, into `deploy/demo/restcol-migrations/` behind
// a drift test, and it was the right stopgap for a deployment that was broken.
//
// But a copy of a file whose original is versioned elsewhere drifts, and the
// drift is silent in the direction that matters: the copy keeps working while
// the pinned version moves underneath it. Exporting the FS makes the pinned
// version and the applied migrations the same fact by construction, which is
// what FootprintAI/grandturks#1004 asked for.
//
// Consumers apply these with golang-migrate's iofs source driver:
//
//	src, err := iofs.New(migrations.EmbeddedMigrations, ".")
//	m, err := migrate.NewWithSourceInstance("iofs", src, databaseURL)
//	err = m.Up()
//
// Files follow golang-migrate's naming convention,
// `<version>_<slug>.{up,down}.sql`; see README.md in this directory.
package migrations

import "embed"

// EmbeddedMigrations holds every migration in this directory, in both
// directions.
//
// The pattern is deliberately `*.sql` at this level and nothing deeper: these
// migrations are a flat, ordered sequence, and a subdirectory would be a
// second ordering that golang-migrate does not read. TestEmbeddedMigrations-
// CoverEveryFileOnDisk fails if a .sql file ever appears outside this pattern,
// so an accidental subdirectory is caught rather than silently skipped.
//
//go:embed *.sql
var EmbeddedMigrations embed.FS
