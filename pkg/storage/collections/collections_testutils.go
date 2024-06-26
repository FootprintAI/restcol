package storagecollections

import (
	"context"

	apppb "github.com/footprintai/restcol/api/pb"
	appmodelcollections "github.com/footprintai/restcol/pkg/models/collections"
	appmodelprojects "github.com/footprintai/restcol/pkg/models/projects"
	dotnotation "github.com/footprintai/restcol/pkg/notation/dot"
	storagepostgres "github.com/sdinsure/agent/pkg/storage/postgres"
)

func TestCollectionSuite(
	postgrescli *storagepostgres.PostgresDb,
	modelProject *appmodelprojects.ModelProject,
) (*appmodelcollections.ModelCollection, error) {
	ctx := context.Background()

	tcrud := &CollectionCURD{postgrescli}
	if err := tcrud.AutoMigrate(); err != nil {
		return nil, err
	}

	cid := appmodelcollections.NewCollectionID()
	mc := appmodelcollections.NewModelCollection(
		modelProject.ID,
		cid,
		apppb.CollectionType_COLLECTION_TYPE_REGULAR_FILES,
		"from testsuite",
		[]*appmodelcollections.ModelSchema{
			&appmodelcollections.ModelSchema{
				Fields: []*appmodelcollections.ModelFieldSchema{
					&appmodelcollections.ModelFieldSchema{
						FieldName:      dotnotation.New("foo"),
						FieldValueType: "string",
					},
					&appmodelcollections.ModelFieldSchema{
						FieldName:      dotnotation.New("bar"),
						FieldValueType: "string",
					},
				},
			},
		},
	)
	if err := tcrud.Write(ctx, "", &mc); err != nil {
		return nil, err
	}
	return &mc, nil
}
