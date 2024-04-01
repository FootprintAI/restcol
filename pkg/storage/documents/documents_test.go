package storagedocuments

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"

	appmodeldocuments "github.com/footprintai/restcol/pkg/models/documents"
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
		Data:              datatypes.JSON("{\"foo\": \"bar\"}"),
		ModelCollectionID: modelCollection.ID,
	}
	assert.Nil(t, dcrud.Write(ctx, "", record))

	found, err := dcrud.Get(ctx, "", record.ID)
	assert.Nil(t, err)
	assert.EqualValues(t, found, record)
}
