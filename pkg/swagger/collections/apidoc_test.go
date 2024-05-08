package collectiondoc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	modelcollections "github.com/footprintai/restcol/pkg/models/collections"
	modelprojects "github.com/footprintai/restcol/pkg/models/projects"
	dotnotation "github.com/footprintai/restcol/pkg/notation/dot"
)

func TestCollectionDoc(t *testing.T) {
	cid1 := modelcollections.NewCollectionID()

	swagVal := func(v any) *modelcollections.SwagValueValue {
		return modelcollections.Must(
			modelcollections.NewSwagValue(v),
		)
	}

	c1 := &modelcollections.ModelCollection{
		ID:             cid1,
		ModelProjectID: modelprojects.NewProjectID(1),
		Summary:        "test swagger doc generation",
		Schemas: []*modelcollections.ModelSchema{
			&modelcollections.ModelSchema{
				Fields: []*modelcollections.ModelFieldSchema{
					&modelcollections.ModelFieldSchema{
						FieldName:      dotnotation.New("foo"),
						FieldValueType: modelcollections.StringSwagValueType,
						FieldExample:   swagVal("fooval"),
					},
					&modelcollections.ModelFieldSchema{
						FieldName:      dotnotation.New("foostruct", "bar"),
						FieldValueType: modelcollections.StringSwagValueType,
						FieldExample:   swagVal("foostruct.barval"),
					},
					&modelcollections.ModelFieldSchema{
						FieldName:      dotnotation.New("bar"),
						FieldValueType: modelcollections.StringSwagValueType,
						FieldExample:   swagVal("barval"),
					},
				},
			},
		},
	}

	csd := NewCollectionSwaggerDoc(c1)
	doc, err := csd.RenderDoc()
	assert.NoError(t, err)

	fmt.Printf("doc:%s\n", doc)
}
