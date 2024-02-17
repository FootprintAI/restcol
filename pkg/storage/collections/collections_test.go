package storagecollections

import (
	"context"
	"testing"

	appmodelcollections "github.com/footprintai/restcol/pkg/models/collections"
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

	ctx := context.Background()
	postgrescli, err := storagetestutils.NewTestPostgresCli(logger.NewLogger())
	assert.NoError(t, err)
	tcrud := &CollectionCURD{postgrescli}
	assert.Nil(t, tcrud.AutoMigrate())

	cid := appmodelcollections.NewCollectionID()

	mc := appmodelcollections.ModelCollection{
		ID:      cid,
		Summary: "test description",
		Schemas: []appmodelcollections.ModelSchema{
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
	}

	assert.Nil(t, tcrud.Write(ctx, "", &mc))

	m, err := tcrud.GetLatestSchema(ctx, "", cid)
	assert.Nil(t, err)

	assert.EqualValues(t, m.Summary, mc.Summary)
	assert.EqualValues(t, len(m.Schemas), 1)
	assert.EqualValues(t, len(m.Schemas[0].Fields), 2)
	assert.EqualValues(t, m.Schemas[0].Fields[0].FieldName, "foo")
	assert.EqualValues(t, m.Schemas[0].Fields[0].FieldValueType, "string")
	assert.EqualValues(t, m.Schemas[0].Fields[1].FieldName, "bar")
	assert.EqualValues(t, m.Schemas[0].Fields[1].FieldValueType, "string")

	// change schema and desc
	mc2 := appmodelcollections.ModelCollection{
		ID:      cid,
		Summary: "test description - part 2",
		Schemas: []appmodelcollections.ModelSchema{
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
	}
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
}
