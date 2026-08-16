package storagedocuments

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	appmodelcollections "github.com/footprintai/restcol/pkg/models/collections"
	appmodeldocuments "github.com/footprintai/restcol/pkg/models/documents"
	appmodelprojects "github.com/footprintai/restcol/pkg/models/projects"
	storagecollectionstestutils "github.com/footprintai/restcol/pkg/storage/collections"
	storageprojects "github.com/footprintai/restcol/pkg/storage/projects"
	storagetestutils "github.com/footprintai/restcol/pkg/storage/testutil"
	"github.com/sdinsure/agent/pkg/logger"
)

func TestDocument(t *testing.T) {
	// launch postgres with the following command
	// docker run --rm --name postgres \
	// -e TZ=gmt+8 \
	// -e POSTGRES_USER=postgres \
	// -e POSTGRES_PASSWORD=password \
	// -e POSTGRES_DB=unittest \
	// -p 5432:5432 -d library/postgres:14.1
	//
	// or run ./run_postgre.sh

	if testing.Short() {
		t.Skip("skip this for now")
		return
	}
	ctx := context.Background()
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger(false))
	assert.NoError(t, err)

	regularProject, _, err := storageprojects.TestProjectSuite(postgrescli)
	assert.Nil(t, err)

	modelCollection, err := storagecollectionstestutils.TestCollectionSuite(postgrescli, regularProject)
	assert.NoError(t, err)

	dcrud := &DocumentCURD{postgrescli}
	assert.Nil(t, dcrud.AutoMigrate())

	record := &appmodeldocuments.ModelDocument{
		ID:                appmodeldocuments.NewDocumentID(),
		Data:              appmodeldocuments.NewModelDocumentData(map[string]interface{}{"foo": "bar"}),
		ModelCollectionID: modelCollection.ID,
		ModelProjectID:    regularProject.ID,
	}
	assert.Nil(t, dcrud.Write(ctx, "", record))

	found, err := dcrud.Get(ctx, "", regularProject.ID, modelCollection.ID, record.ID)
	assert.Nil(t, err)
	assert.EqualValues(t, found, record)

}

func TestDocumentQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("skip this for now")
		return
	}
	ctx := context.Background()
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger(false))
	assert.NoError(t, err)

	regularProject, _, err := storageprojects.TestProjectSuite(postgrescli)
	assert.Nil(t, err)

	modelCollection, err := storagecollectionstestutils.TestCollectionSuite(postgrescli, regularProject)
	assert.NoError(t, err)

	dcrud := &DocumentCURD{postgrescli}
	assert.Nil(t, dcrud.AutoMigrate())

	records := newDocs(regularProject.ID, modelCollection.ID, 100)
	assert.Nil(t, dcrud.BatchWrite(ctx, "", records))

	queryTime := time.Now()
	queryDocs, err := dcrud.Query(
		ctx,
		"",
		regularProject.ID,
		modelCollection.ID,
		WithEndedAt(queryTime),
		WithLimitCount(101),
	)
	assert.Nil(t, err)
	assert.Len(t, queryDocs, 100)

	// write 2nd batches
	records = newDocs(regularProject.ID, modelCollection.ID, 100)
	assert.Nil(t, dcrud.BatchWrite(ctx, "", records))

	// query should get the same results
	queryDocs, err = dcrud.Query(
		ctx,
		"",
		regularProject.ID,
		modelCollection.ID,
		WithEndedAt(queryTime),
		WithLimitCount(101),
	)
	assert.Nil(t, err)
	assert.Len(t, queryDocs, 100)
}

func newDocs(pid appmodelprojects.ProjectID, cid appmodelcollections.CollectionID, count int) []*appmodeldocuments.ModelDocument {
	docs := []*appmodeldocuments.ModelDocument{}

	for i := 0; i < count; i++ {
		did := appmodeldocuments.NewDocumentID()
		record := &appmodeldocuments.ModelDocument{
			ID:                did,
			Data:              appmodeldocuments.NewModelDocumentData(map[string]interface{}{"foo": "bar", "myid": did.String()}),
			ModelCollectionID: cid,
			ModelProjectID:    pid,
		}
		docs = append(docs, record)
	}
	return docs
}

func TestDocumentSameID(t *testing.T) {
	if testing.Short() {
		t.Skip("skip this for now")
		return
	}
	ctx := context.Background()
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger(false))
	assert.NoError(t, err)

	regularProject, _, err := storageprojects.TestProjectSuite(postgrescli)
	assert.Nil(t, err)

	modelCollection1, err := storagecollectionstestutils.TestCollectionSuite(postgrescli, regularProject)
	assert.NoError(t, err)

	modelCollection2, err := storagecollectionstestutils.TestCollectionSuite(postgrescli, regularProject)
	assert.NoError(t, err)

	dcrud := &DocumentCURD{postgrescli}
	assert.Nil(t, dcrud.AutoMigrate())

	record := &appmodeldocuments.ModelDocument{
		ID:                appmodeldocuments.NewDocumentID(),
		Data:              appmodeldocuments.NewModelDocumentData(map[string]interface{}{"foo": "bar"}),
		ModelCollectionID: modelCollection1.ID,
		ModelProjectID:    regularProject.ID,
	}
	assert.Nil(t, dcrud.Write(ctx, "", record))

	// write again with same did but different cid
	record2 := &appmodeldocuments.ModelDocument{
		ID:                record.ID,
		Data:              appmodeldocuments.NewModelDocumentData(map[string]interface{}{"foo": "bar"}),
		ModelCollectionID: modelCollection2.ID,
		ModelProjectID:    regularProject.ID,
	}
	assert.Nil(t, dcrud.Write(ctx, "", record2))

	found1, err := dcrud.Get(ctx, "", regularProject.ID, modelCollection1.ID, record.ID)
	assert.Nil(t, err)
	found2, err := dcrud.Get(ctx, "", regularProject.ID, modelCollection2.ID, record.ID)
	assert.Nil(t, err)
	assert.True(t, found1.ID == found2.ID)
	assert.True(t, found1.ModelProjectID == found2.ModelProjectID)
	assert.True(t, found1.ModelCollectionID != found2.ModelCollectionID)
	assert.False(t, reflect.DeepEqual(found1, found2))

}

func TestDocumentCURD_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger(false))
	assert.NoError(t, err)

	regularProject, _, err := storageprojects.TestProjectSuite(postgrescli)
	assert.Nil(t, err)

	docCURD := NewDocumentCURD(postgrescli)
	assert.Nil(t, docCURD.AutoMigrate())

	collectionsCURD := storagecollectionstestutils.NewCollectionCURD(postgrescli)
	assert.Nil(t, collectionsCURD.AutoMigrate())

	// Create a test collection
	collection, err := storagecollectionstestutils.TestCollectionSuite(postgrescli, regularProject)
	assert.NoError(t, err)
	err = collectionsCURD.Write(context.Background(), "", collection)
	assert.Nil(t, err)

	// Create a test document
	testDoc := &appmodeldocuments.ModelDocument{
		ID:                appmodeldocuments.NewDocumentID(),
		Data:              appmodeldocuments.NewModelDocumentData(map[string]interface{}{"test": "data"}),
		ModelCollectionID: collection.ID,
		ModelProjectID:    regularProject.ID,
	}

	err = docCURD.Write(context.Background(), "", testDoc)
	assert.NoError(t, err)

	// Verify document exists
	foundDoc, err := docCURD.Get(context.Background(), "", regularProject.ID, collection.ID, testDoc.ID)
	assert.NoError(t, err)
	assert.Equal(t, testDoc.ID, foundDoc.ID)

	// Test successful deletion
	err = docCURD.Delete(context.Background(), "", regularProject.ID, collection.ID, testDoc.ID)
	assert.NoError(t, err)

	// Gone from the normal read path...
	deletedDoc, err := docCURD.Get(context.Background(), "", regularProject.ID, collection.ID, testDoc.ID)
	assert.NoError(t, err)                  // Get doesn't error, just returns empty record
	assert.Empty(t, deletedDoc.ID.String()) // ID is empty when the row is not there

	// ...and gone from the table, not merely hidden by gorm's deleted_at
	// filter. This is the assertion that separates a hard delete from a soft
	// one; without Unscoped here the query below would still find the row and
	// the test would pass on a tombstone (#136).
	var remaining int64
	assert.NoError(t, postgrescli.GormDB().Unscoped().
		Model(&appmodeldocuments.ModelDocument{}).
		Where("id = ?", testDoc.ID.String()).
		Count(&remaining).Error)
	assert.Zero(t, remaining, "the row must be deleted, not tombstoned")

	// Test delete non-existent document (should not error in GORM)
	nonExistentID := appmodeldocuments.NewDocumentID()
	err = docCURD.Delete(context.Background(), "", regularProject.ID, collection.ID, nonExistentID)
	assert.NoError(t, err) // GORM delete doesn't error for non-existent records
}

func TestDocumentCURD_Delete_WithWrongScope(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger(false))
	assert.NoError(t, err)

	regularProject, proxyProject, err := storageprojects.TestProjectSuite(postgrescli)
	assert.Nil(t, err)

	docCURD := NewDocumentCURD(postgrescli)
	assert.Nil(t, docCURD.AutoMigrate())

	collectionsCURD := storagecollectionstestutils.NewCollectionCURD(postgrescli)
	assert.Nil(t, collectionsCURD.AutoMigrate())

	// Create collections in different projects
	regularCollection, err := storagecollectionstestutils.TestCollectionSuite(postgrescli, regularProject)
	assert.NoError(t, err)
	proxyCollection, err := storagecollectionstestutils.TestCollectionSuite(postgrescli, proxyProject)
	assert.NoError(t, err)

	err = collectionsCURD.Write(context.Background(), "", regularCollection)
	assert.Nil(t, err)
	err = collectionsCURD.Write(context.Background(), "", proxyCollection)
	assert.Nil(t, err)

	// Create document in regular project
	testDoc := &appmodeldocuments.ModelDocument{
		ID:                appmodeldocuments.NewDocumentID(),
		Data:              appmodeldocuments.NewModelDocumentData(map[string]interface{}{"test": "data"}),
		ModelCollectionID: regularCollection.ID,
		ModelProjectID:    regularProject.ID,
	}

	err = docCURD.Write(context.Background(), "", testDoc)
	assert.NoError(t, err)

	// Try to delete with wrong project scope (should not delete anything)
	err = docCURD.Delete(context.Background(), "", proxyProject.ID, regularCollection.ID, testDoc.ID)
	assert.NoError(t, err) // No error, but nothing deleted

	// Verify document still exists in correct scope
	foundDoc, err := docCURD.Get(context.Background(), "", regularProject.ID, regularCollection.ID, testDoc.ID)
	assert.NoError(t, err)
	assert.Equal(t, testDoc.ID, foundDoc.ID)

	// Try to delete with wrong collection scope (should not delete anything)
	err = docCURD.Delete(context.Background(), "", regularProject.ID, proxyCollection.ID, testDoc.ID)
	assert.NoError(t, err) // No error, but nothing deleted

	// Verify document still exists in correct scope
	foundDoc, err = docCURD.Get(context.Background(), "", regularProject.ID, regularCollection.ID, testDoc.ID)
	assert.NoError(t, err)
	assert.Equal(t, testDoc.ID, foundDoc.ID)
}

// TestDocumentCURD_DeleteByCollection_IsPermanent covers the cascade path that
// DeleteCollection(force=true) uses. It is the one most likely to be missed:
// deleting a whole collection is exactly when a caller expects the payloads to
// be gone, and a scoped delete there would retain every one of them.
func TestDocumentCURD_DeleteByCollection_IsPermanent(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	ctx := context.Background()
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger(false))
	assert.NoError(t, err)

	regularProject, _, err := storageprojects.TestProjectSuite(postgrescli)
	assert.Nil(t, err)

	collection, err := storagecollectionstestutils.TestCollectionSuite(postgrescli, regularProject)
	assert.NoError(t, err)

	docCURD := NewDocumentCURD(postgrescli)
	assert.Nil(t, docCURD.AutoMigrate())

	for i := 0; i < 3; i++ {
		assert.NoError(t, docCURD.Write(ctx, "", &appmodeldocuments.ModelDocument{
			ID:                appmodeldocuments.NewDocumentID(),
			Data:              appmodeldocuments.NewModelDocumentData(map[string]interface{}{"n": i}),
			ModelCollectionID: collection.ID,
			ModelProjectID:    regularProject.ID,
		}))
	}

	assert.NoError(t, docCURD.DeleteByCollection(ctx, "", regularProject.ID, collection.ID))

	// Unscoped: counts tombstones as well as live rows, so this fails if the
	// cascade only marked them deleted.
	var remaining int64
	assert.NoError(t, postgrescli.GormDB().Unscoped().
		Model(&appmodeldocuments.ModelDocument{}).
		Where("model_project_id = ? AND model_collection_id = ?",
			regularProject.ID.String(), collection.ID.String()).
		Count(&remaining).Error)
	assert.Zero(t, remaining, "the cascade must remove the rows, not tombstone them")
}

// Writing an existing documentId is an UPSERT - there is no Update rpc - and
// the metadata has to record that it happened. Before FootprintAI/restcol#141
// the API exposed only createdAt and deletedAt, so an overwritten document was
// indistinguishable from one written once.
//
// The storage layer was already correct: Write's OnConflict lists updated_at in
// DoUpdates and deliberately omits created_at. This pins that behaviour, because
// the new API field is only meaningful if these two move independently - if
// created_at ever joined the DoUpdates set, _updatedAt would still be populated
// and still tell a caller nothing.
func TestUpsertAdvancesUpdatedAtAndPreservesCreatedAt(t *testing.T) {
	if testing.Short() {
		t.Skip("needs postgres")
		return
	}
	ctx := context.Background()
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger(false))
	assert.NoError(t, err)

	regularProject, _, err := storageprojects.TestProjectSuite(postgrescli)
	assert.Nil(t, err)
	collection, err := storagecollectionstestutils.TestCollectionSuite(postgrescli, regularProject)
	assert.NoError(t, err)

	dcrud := NewDocumentCURD(postgrescli)
	assert.Nil(t, dcrud.AutoMigrate())

	docID := appmodeldocuments.NewDocumentID()
	write := func(v string) {
		assert.NoError(t, dcrud.Write(ctx, "", &appmodeldocuments.ModelDocument{
			ID:                docID,
			Data:              appmodeldocuments.NewModelDocumentData(map[string]interface{}{"v": v}),
			ModelCollectionID: collection.ID,
			ModelProjectID:    regularProject.ID,
		}))
	}

	write("first")
	first, err := dcrud.Get(ctx, "", regularProject.ID, collection.ID, docID)
	assert.NoError(t, err)

	// A fresh row: gorm stamps both columns on insert, which is what makes
	// _updatedAt safe to expose as non-optional.
	assert.False(t, first.UpdatedAt.IsZero(), "updated_at must be set on insert")
	assert.WithinDuration(t, first.CreatedAt, first.UpdatedAt, time.Second,
		"on a document written once, updated_at should match created_at")

	// Postgres stores microseconds, so without a pause the two writes can land
	// on the same instant and the assertion below would pass or fail on timing
	// rather than on behaviour.
	time.Sleep(10 * time.Millisecond)

	write("second")
	second, err := dcrud.Get(ctx, "", regularProject.ID, collection.ID, docID)
	assert.NoError(t, err)

	assert.True(t, second.UpdatedAt.After(first.UpdatedAt),
		"the upsert must advance updated_at; without it an overwrite leaves no trace")
	assert.WithinDuration(t, first.CreatedAt, second.CreatedAt, time.Millisecond,
		"created_at must survive the upsert - it is what _updatedAt is compared against")

	// And the upsert really did replace the payload rather than insert a second
	// row, which is the premise the timestamps are describing.
	assert.Equal(t, "second", second.Data.MapValue["v"])
	count, err := dcrud.CountByCollection(ctx, "", regularProject.ID, collection.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, count, "upsert must not append a second row")
}

// The scenario this pair of columns exists for: one principal creates a
// document, a DIFFERENT one overwrites it, and both are recoverable.
//
// Writing an existing documentId is an upsert and there is no Update rpc, so
// before created_by/updated_by there was no way to tell that a second party had
// replaced someone else's data - not from the content, not from the metadata.
//
// It pins the asymmetry that makes the pair meaningful: updated_by is in the
// upsert's DoUpdates and created_by is not. If created_by ever joined that
// list, BOTH fields would report the most recent writer, the pair would answer
// nothing, and every other assertion here would still pass.
func TestUpsertByADifferentPrincipalKeepsCreatorAndMovesUpdater(t *testing.T) {
	if testing.Short() {
		t.Skip("needs postgres")
		return
	}
	ctx := context.Background()
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger(false))
	assert.NoError(t, err)

	regularProject, _, err := storageprojects.TestProjectSuite(postgrescli)
	assert.Nil(t, err)
	collection, err := storagecollectionstestutils.TestCollectionSuite(postgrescli, regularProject)
	assert.NoError(t, err)

	dcrud := NewDocumentCURD(postgrescli)
	assert.Nil(t, dcrud.AutoMigrate())

	docID := appmodeldocuments.NewDocumentID()
	write := func(principal, value string) {
		assert.NoError(t, dcrud.Write(ctx, "", &appmodeldocuments.ModelDocument{
			ID:                docID,
			CreatedBy:         principal,
			UpdatedBy:         principal,
			Data:              appmodeldocuments.NewModelDocumentData(map[string]interface{}{"v": value}),
			ModelCollectionID: collection.ID,
			ModelProjectID:    regularProject.ID,
		}))
	}

	const alice = "grandturk:apikey:42:alice-key"
	const bob = "grandturk:apikey:42:bob-key"

	write(alice, "first")
	first, err := dcrud.Get(ctx, "", regularProject.ID, collection.ID, docID)
	assert.NoError(t, err)
	assert.Equal(t, alice, first.CreatedBy)
	assert.Equal(t, alice, first.UpdatedBy,
		"on a document written once, the creator is also the last writer")

	write(bob, "second")
	second, err := dcrud.Get(ctx, "", regularProject.ID, collection.ID, docID)
	assert.NoError(t, err)

	assert.Equal(t, alice, second.CreatedBy,
		"the creator must survive an overwrite by someone else - this is the "+
			"half that fails if created_by joins the upsert's DoUpdates")
	assert.Equal(t, bob, second.UpdatedBy,
		"the last writer must be the principal that actually overwrote it")
	assert.Equal(t, "second", second.Data.MapValue["v"],
		"and the overwrite must really have replaced the payload")

	count, err := dcrud.CountByCollection(ctx, "", regularProject.ID, collection.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, count, "an upsert must not append a second row")
}

// An unattributed write records an empty writer rather than a placeholder.
// A deployment with no caller resolver must be distinguishable from one whose
// documents were genuinely written by a principal called "system".
func TestUnattributedWritesRecordNoPrincipal(t *testing.T) {
	if testing.Short() {
		t.Skip("needs postgres")
		return
	}
	ctx := context.Background()
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger(false))
	assert.NoError(t, err)

	regularProject, _, err := storageprojects.TestProjectSuite(postgrescli)
	assert.Nil(t, err)
	collection, err := storagecollectionstestutils.TestCollectionSuite(postgrescli, regularProject)
	assert.NoError(t, err)

	dcrud := NewDocumentCURD(postgrescli)
	assert.Nil(t, dcrud.AutoMigrate())

	docID := appmodeldocuments.NewDocumentID()
	assert.NoError(t, dcrud.Write(ctx, "", &appmodeldocuments.ModelDocument{
		ID:                docID,
		Data:              appmodeldocuments.NewModelDocumentData(map[string]interface{}{"v": "x"}),
		ModelCollectionID: collection.ID,
		ModelProjectID:    regularProject.ID,
	}))

	got, err := dcrud.Get(ctx, "", regularProject.ID, collection.ID, docID)
	assert.NoError(t, err)
	assert.Empty(t, got.CreatedBy, "an unattributed document must not invent a writer")
	assert.Empty(t, got.UpdatedBy)
}
