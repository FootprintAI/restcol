package storagecollections

import (
	"context"
	"testing"

	apppb "github.com/footprintai/restcol/api/pb/proto"
	appmodelcollections "github.com/footprintai/restcol/pkg/models/collections"
	storageprojects "github.com/footprintai/restcol/pkg/storage/projects"
	"github.com/sdinsure/agent/pkg/logger"
	storagetestutils "github.com/sdinsure/agent/pkg/storage/testutils"
	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
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

	regularProject, proxyProject, err := storageprojects.TestProjectSuite(postgrescli)
	assert.NoError(t, err)

	tcrud := &CollectionCURD{postgrescli}
	assert.Nil(t, tcrud.AutoMigrate())

	cid := appmodelcollections.NewCollectionID()
	_, err = TestCollectionSuite(postgrescli, proxyProject)
	assert.Nil(t, err)

	mc := appmodelcollections.NewModelCollection(
		regularProject,
		cid,
		apppb.CollectionType_COLLECTION_TYPE_REGULAR_FILES,
		"test description",
		[]appmodelcollections.ModelSchema{
			appmodelcollections.ModelSchema{
				Fields: []appmodelcollections.ModelFieldSchema{
					appmodelcollections.ModelFieldSchema{
						FieldName:      "foo",
						FieldValueType: "string",
					},
					appmodelcollections.ModelFieldSchema{
						FieldName:      "bar",
						FieldValueType: "string",
					},
				},
			},
		},
	)

	assert.Nil(t, tcrud.Write(ctx, "", &mc))

	m, err := tcrud.GetLatestSchema(ctx, "", cid)
	assert.Nil(t, err)
	assert.NotNil(t, m)

	assert.EqualValues(t, m.Summary, mc.Summary)
	assert.EqualValues(t, len(m.Schemas), 1)
	assert.EqualValues(t, len(m.Schemas[0].Fields), 2)
	assert.True(t, m.Schemas[0].ID > 0)
	assert.EqualValues(t, m.Schemas[0].Fields[0].FieldName, "foo")
	assert.EqualValues(t, m.Schemas[0].Fields[0].FieldValueType, "string")
	assert.EqualValues(t, m.Schemas[0].Fields[1].FieldName, "bar")
	assert.EqualValues(t, m.Schemas[0].Fields[1].FieldValueType, "string")

	schemaId1 := m.Schemas[0].ID

	// change schema and desc
	mc2 := appmodelcollections.NewModelCollection(
		regularProject,
		cid,
		apppb.CollectionType_COLLECTION_TYPE_REGULAR_FILES,
		"test description - part 2",
		[]appmodelcollections.ModelSchema{
			appmodelcollections.ModelSchema{
				Fields: []appmodelcollections.ModelFieldSchema{
					appmodelcollections.ModelFieldSchema{
						FieldName:      "foo",
						FieldValueType: "string",
					},
					appmodelcollections.ModelFieldSchema{
						FieldName:      "bar2",
						FieldValueType: "string",
					},
					appmodelcollections.ModelFieldSchema{
						FieldName:      "bar3",
						FieldValueType: "float64",
					},
				},
			},
		},
	)

	assert.Nil(t, tcrud.Update(ctx, "", &mc2))

	m, err = tcrud.GetLatestSchema(ctx, "", cid)
	assert.Nil(t, err)

	assert.EqualValues(t, m.Summary, mc2.Summary)
	assert.EqualValues(t, len(m.Schemas), 1)
	assert.EqualValues(t, len(m.Schemas[0].Fields), 3)
	assert.EqualValues(t, m.Schemas[0].Fields[0].FieldName, "foo")
	assert.EqualValues(t, m.Schemas[0].Fields[0].FieldValueType, "string")
	assert.EqualValues(t, m.Schemas[0].Fields[1].FieldName, "bar2")
	assert.EqualValues(t, m.Schemas[0].Fields[1].FieldValueType, "string")
	assert.EqualValues(t, m.Schemas[0].Fields[2].FieldName, "bar3")
	assert.EqualValues(t, m.Schemas[0].Fields[2].FieldValueType, "float64")

	schemaId2 := m.Schemas[0].ID
	assert.True(t, schemaId2 > schemaId1)
}
