// Package storagetestutil provides the postgres handle restcol's tests run
// against.
//
// It wraps sdinsure's NewTestPostgresCli and adjusts one thing: gorm's clock.
package storagetestutil

import (
	"time"

	"github.com/sdinsure/agent/pkg/logger"
	storagepostgres "github.com/sdinsure/agent/pkg/storage/postgres"
	storagetestutils "github.com/sdinsure/agent/pkg/storage/testutils"
)

// NewTestPostgresCli returns a postgres handle whose gorm clock is truncated to
// microseconds.
//
// gorm stamps CreatedAt/UpdatedAt from Go's clock, which carries nanoseconds.
// Postgres timestamp columns store microseconds. Without this, the struct left
// in memory after a write and the struct read back from the database differ in
// the last three digits:
//
//	written:  2026-08-13 11:58:56.333178631
//	read back: 2026-08-13 11:58:56.333178000
//
// Every whole-struct comparison then fails on a field the test never set, and
// the failure looks like a storage bug rather than a clock-resolution
// mismatch. Truncating here makes the in-memory value equal to the stored one
// for every model carrying CreatedAt/UpdatedAt, rather than fixing one
// assertion at a time.
//
// Tests that need the raw sdinsure handle can still call it directly; this is
// the default for anything that compares records it wrote.
func NewTestPostgresCli(log logger.Logger) (*storagepostgres.PostgresDb, error) {
	db, err := storagetestutils.NewTestPostgresCli(log)
	if err != nil {
		return nil, err
	}
	db.GormDB().NowFunc = func() time.Time {
		return time.Now().Truncate(time.Microsecond)
	}
	return db, nil
}
