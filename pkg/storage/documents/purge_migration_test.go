package storagedocuments

import (
	"context"
	"testing"

	"github.com/sdinsure/agent/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/footprintai/restcol/migrations"
	appmodeldocuments "github.com/footprintai/restcol/pkg/models/documents"
	storagecollections "github.com/footprintai/restcol/pkg/storage/collections"
	storageprojects "github.com/footprintai/restcol/pkg/storage/projects"
	storagetestutils "github.com/footprintai/restcol/pkg/storage/testutil"
)

// TestMigration0002PurgesTombstonesAndKeepsLiveRows proves what
// 0002_purge_soft_deleted_documents claims, which until now nothing did.
//
// # WHY THIS WAS MISSING, AND WHY IT MATTERED
//
// migrations/embed_test.go says so in its own header: those tests are "about
// COVERAGE, not about the SQL" - they assert every file is reachable through
// the embed, not that any of them does anything. So `0002` shipped, was pinned
// by grandturks (#1004/#1014), and was applied to a deployed cluster, without a
// single assertion that it removes a tombstone.
//
// The gap surfaced when FootprintAI/grandturks#977's criterion 4 - "after
// #970's migration, count(*) WHERE deleted_at IS NOT NULL returns 0" - was run
// against ubm01. It returned 0. It had also returned 0 BEFORE the migration ran,
// because that database had never held a tombstone: `DELETE 0`. The criterion
// passed while demonstrating nothing, which is the exact failure sprint #827
// closed on, one level down.
//
// This test is the missing half. It cannot be satisfied vacuously: it CREATES a
// tombstone, asserts it is there, then asserts the migration removed it.
//
// # THE TOMBSTONE IS MADE THE WAY THE OLD CODE MADE IT
//
// DocumentCURD.Delete is Unscoped now (#136/#137), so it hard-deletes and
// cannot produce the rows this migration exists to purge. A scoped gorm delete
// is what the pre-#137 code did, and it is what wrote every tombstone in every
// deployed database. Reproducing it here rather than INSERTing a hand-made row
// keeps the fixture honest: gorm sets deleted_at exactly as it did in
// production, and if that behaviour ever changes this test changes with it.
//
// # THE SECOND ASSERTION IS THE ONE THAT PROTECTS DATA
//
// "Tombstones are gone" is half the contract. "Live rows are untouched" is the
// half that makes the migration safe to run, and a `DELETE FROM ... WHERE
// deleted_at IS NOT NULL` with a broken predicate would still satisfy the
// first. Both are asserted, and the keeper is read back through the normal path
// rather than counted, so a row that survived as a tombstone would not pass.
func TestMigration0002PurgesTombstonesAndKeepsLiveRows(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}

	ctx := context.Background()
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger(false))
	require.NoError(t, err)

	regularProject, _, err := storageprojects.TestProjectSuite(postgrescli)
	require.NoError(t, err)

	docCURD := NewDocumentCURD(postgrescli)
	require.NoError(t, docCURD.AutoMigrate())

	collectionsCURD := storagecollections.NewCollectionCURD(postgrescli)
	require.NoError(t, collectionsCURD.AutoMigrate())

	collection, err := storagecollections.TestCollectionSuite(postgrescli, regularProject)
	require.NoError(t, err)
	require.NoError(t, collectionsCURD.Write(ctx, "", collection))

	newDoc := func(payload map[string]interface{}) *appmodeldocuments.ModelDocument {
		return &appmodeldocuments.ModelDocument{
			ID:                appmodeldocuments.NewDocumentID(),
			Data:              appmodeldocuments.NewModelDocumentData(payload),
			ModelCollectionID: collection.ID,
			ModelProjectID:    regularProject.ID,
		}
	}

	keeper := newDoc(map[string]interface{}{"keep": "this row must survive"})
	require.NoError(t, docCURD.Write(ctx, "", keeper))

	// The payload names what is actually at stake: for taskWorker-sinked
	// documents the `data` jsonb holds base64 image bytes, which is why #936
	// raised retention over these rows rather than tidiness.
	doomed := newDoc(map[string]interface{}{"imageBytes": "iVBORw0KGgoAAAANSUhEUg=="})
	require.NoError(t, docCURD.Write(ctx, "", doomed))

	// A SCOPED delete - what the code did before #137 - leaves the row with
	// deleted_at set and its data intact.
	require.NoError(t, postgrescli.GormDB().
		Where("id = ?", doomed.ID.String()).
		Delete(&appmodeldocuments.ModelDocument{}).Error)

	countTombstones := func() int64 {
		var n int64
		require.NoError(t, postgrescli.GormDB().Unscoped().
			Model(&appmodeldocuments.ModelDocument{}).
			Where("deleted_at IS NOT NULL").
			Count(&n).Error)
		return n
	}

	// Guard the fixture itself. If this is 0 the migration below would "pass"
	// against nothing, which is precisely the vacuous result that prompted the
	// test - so it is a hard failure, not an assumption.
	require.NotZero(t, countTombstones(),
		"fixture did not produce a tombstone; the rest of this test would be vacuous")

	// The row is unreachable through the API and still on disk - the state
	// #936 objected to.
	var doomedData string
	require.NoError(t, postgrescli.GormDB().Unscoped().
		Model(&appmodeldocuments.ModelDocument{}).
		Where("id = ?", doomed.ID.String()).
		Select("data->>'imageBytes'").Scan(&doomedData).Error)
	assert.Equal(t, "iVBORw0KGgoAAAANSUhEUg==", doomedData,
		"the tombstoned row should still hold its payload before the purge")

	// Apply the migration from the EMBEDDED copy, not from disk - that is what
	// a consumer pinning this module actually runs (#1004).
	up, err := migrations.EmbeddedMigrations.ReadFile("0002_purge_soft_deleted_documents.up.sql")
	require.NoError(t, err, "the migration must be reachable through the embed")
	require.NoError(t, postgrescli.GormDB().Exec(string(up)).Error)

	assert.Zero(t, countTombstones(),
		"0002 left a soft-deleted row behind - the migration does not do what it says")

	// And the live row is still readable through the normal path. Counting
	// would not be enough: a keeper that survived AS A TOMBSTONE would still be
	// counted, and would still be broken.
	found, err := docCURD.Get(ctx, "", regularProject.ID, collection.ID, keeper.ID)
	require.NoError(t, err)
	assert.Equal(t, keeper.ID, found.ID,
		"0002 removed a LIVE document - its predicate is wrong and it is unsafe to run")
}

// TestMigration0002IsIdempotent - #970's runbook tells operators "re-running is
// harmless", and that claim is load-bearing: the deployed sites are applied by
// hand, so a second run is likely rather than hypothetical.
func TestMigration0002IsIdempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}

	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger(false))
	require.NoError(t, err)

	docCURD := NewDocumentCURD(postgrescli)
	require.NoError(t, docCURD.AutoMigrate())

	up, err := migrations.EmbeddedMigrations.ReadFile("0002_purge_soft_deleted_documents.up.sql")
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		require.NoErrorf(t, postgrescli.GormDB().Exec(string(up)).Error,
			"run %d of 0002 failed; the runbook promises re-running is harmless", i+1)
	}
}
