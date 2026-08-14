package migrations

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These migrations are consumed by OTHER modules — grandturks pins restcol and
// applies exactly the migrations that shipped with the version it pinned
// (FootprintAI/grandturks#1004). Go's embed cannot reach into a dependency, so
// without EmbeddedMigrations a consumer's only option is to copy the files,
// which is what grandturks#987 had to do. A copy of a file versioned elsewhere
// drifts; that is the whole reason this exists.
//
// The tests below are therefore about COVERAGE, not about the SQL: every file
// in this directory has to be reachable through the embed, or a consumer
// silently applies a subset and believes it is current.

// readDiskSQL returns every .sql file in this directory, keyed by path
// relative to it — including any in subdirectories, which is exactly the case
// a bare `//go:embed *.sql` would miss.
func readDiskSQL(t *testing.T) map[string][]byte {
	t.Helper()
	found := map[string][]byte{}
	require.NoError(t, filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		found[filepath.ToSlash(path)] = raw
		return nil
	}))
	return found
}

func readEmbeddedSQL(t *testing.T) map[string][]byte {
	t.Helper()
	found := map[string][]byte{}
	require.NoError(t, fs.WalkDir(EmbeddedMigrations, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}
		raw, readErr := EmbeddedMigrations.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		found[path] = raw
		return nil
	}))
	return found
}

// TestEmbeddedMigrationsCoverEveryFileOnDisk is the one that matters. A
// migration that exists on disk but not in the embed is invisible to every
// consumer, and nothing else in the build would say so.
func TestEmbeddedMigrationsCoverEveryFileOnDisk(t *testing.T) {
	disk := readDiskSQL(t)
	embedded := readEmbeddedSQL(t)

	require.NotEmpty(t, disk, "no .sql files found — this test is not looking where it thinks it is")

	for name, want := range disk {
		got, ok := embedded[name]
		if !ok {
			t.Errorf("%s is on disk but not in EmbeddedMigrations.\n"+
				"A consumer pinning this module would apply the other migrations and believe "+
				"its schema is current. Check the //go:embed pattern covers it.", name)
			continue
		}
		assert.Equal(t, string(want), string(got), "%s differs between disk and the embed", name)
	}

	for name := range embedded {
		assert.Contains(t, disk, name, "%s is embedded but not on disk", name)
	}
}

// golang-migrate resolves a version by finding both directions. An up without
// a down is not a partial feature — it is a migration that cannot be rolled
// back, discovered at the moment someone needs to roll it back.
func TestEveryMigrationHasBothDirections(t *testing.T) {
	embedded := readEmbeddedSQL(t)

	for name := range embedded {
		var counterpart string
		switch {
		case strings.HasSuffix(name, ".up.sql"):
			counterpart = strings.TrimSuffix(name, ".up.sql") + ".down.sql"
		case strings.HasSuffix(name, ".down.sql"):
			counterpart = strings.TrimSuffix(name, ".down.sql") + ".up.sql"
		default:
			t.Errorf("%s is neither .up.sql nor .down.sql — golang-migrate will not see it as a migration", name)
			continue
		}
		assert.Contains(t, embedded, counterpart, "%s has no %s", name, counterpart)
	}
}

// An empty migration file applies cleanly and does nothing, so the version
// table advances and the schema does not. That is worse than a failure.
func TestNoEmbeddedMigrationIsEmpty(t *testing.T) {
	for name, content := range readEmbeddedSQL(t) {
		assert.NotEmpty(t, strings.TrimSpace(string(content)),
			"%s is empty: it would advance the schema version without changing the schema", name)
	}
}
