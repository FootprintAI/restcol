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
	"github.com/sdinsure/agent/pkg/logger"
	storagetestutils "github.com/sdinsure/agent/pkg/storage/testutils"
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
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger())
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
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger())
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
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger())
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
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger())
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
	err = collectionsCURD.Write(context.TODO(), "", collection)
	assert.Nil(t, err)

	// Create a test document
	testDoc := &appmodeldocuments.ModelDocument{
		ID:                appmodeldocuments.NewDocumentID(),
		Data:              appmodeldocuments.NewModelDocumentData(map[string]interface{}{"test": "data"}),
		ModelCollectionID: collection.ID,
		ModelProjectID:    regularProject.ID,
	}

	err = docCURD.Write(context.TODO(), "", testDoc)
	assert.NoError(t, err)

	// Verify document exists
	foundDoc, err := docCURD.Get(context.TODO(), "", regularProject.ID, collection.ID, testDoc.ID)
	assert.NoError(t, err)
	assert.Equal(t, testDoc.ID, foundDoc.ID)

	// Test successful deletion
	err = docCURD.Delete(context.TODO(), "", regularProject.ID, collection.ID, testDoc.ID)
	assert.NoError(t, err)

	// Verify document is soft deleted (should return empty record)
	deletedDoc, err := docCURD.Get(context.TODO(), "", regularProject.ID, collection.ID, testDoc.ID)
	assert.NoError(t, err) // Get doesn't error, just returns empty record
	assert.Empty(t, deletedDoc.ID.String()) // ID should be empty for soft-deleted record

	// Test delete non-existent document (should not error in GORM)
	nonExistentID := appmodeldocuments.NewDocumentID()
	err = docCURD.Delete(context.TODO(), "", regularProject.ID, collection.ID, nonExistentID)
	assert.NoError(t, err) // GORM delete doesn't error for non-existent records
}

func TestDocumentCURD_Delete_WithWrongScope(t *testing.T) {
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger())
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
	
	err = collectionsCURD.Write(context.TODO(), "", regularCollection)
	assert.Nil(t, err)
	err = collectionsCURD.Write(context.TODO(), "", proxyCollection)
	assert.Nil(t, err)

	// Create document in regular project
	testDoc := &appmodeldocuments.ModelDocument{
		ID:                appmodeldocuments.NewDocumentID(),
		Data:              appmodeldocuments.NewModelDocumentData(map[string]interface{}{"test": "data"}),
		ModelCollectionID: regularCollection.ID,
		ModelProjectID:    regularProject.ID,
	}

	err = docCURD.Write(context.TODO(), "", testDoc)
	assert.NoError(t, err)

	// Try to delete with wrong project scope (should not delete anything)
	err = docCURD.Delete(context.TODO(), "", proxyProject.ID, regularCollection.ID, testDoc.ID)
	assert.NoError(t, err) // No error, but nothing deleted

	// Verify document still exists in correct scope
	foundDoc, err := docCURD.Get(context.TODO(), "", regularProject.ID, regularCollection.ID, testDoc.ID)
	assert.NoError(t, err)
	assert.Equal(t, testDoc.ID, foundDoc.ID)

	// Try to delete with wrong collection scope (should not delete anything)
	err = docCURD.Delete(context.TODO(), "", regularProject.ID, proxyCollection.ID, testDoc.ID)
	assert.NoError(t, err) // No error, but nothing deleted

	// Verify document still exists in correct scope
	foundDoc, err = docCURD.Get(context.TODO(), "", regularProject.ID, regularCollection.ID, testDoc.ID)
	assert.NoError(t, err)
	assert.Equal(t, testDoc.ID, foundDoc.ID)
}
