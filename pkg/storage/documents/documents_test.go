package storagedocuments

import (
	"context"
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

	found, err := dcrud.Get(ctx, "", record.ID)
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
