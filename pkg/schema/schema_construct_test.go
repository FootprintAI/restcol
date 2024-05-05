package schema

import (
	"testing"

	appmodelcollections "github.com/footprintai/restcol/pkg/models/collections"
	dotnotation "github.com/footprintai/restcol/pkg/notation/dot"
	"github.com/stretchr/testify/assert"
)

func TestSchemaBuilder(t *testing.T) {
	fields := []*appmodelcollections.ModelFieldSchema{
		&appmodelcollections.ModelFieldSchema{
			FieldName:    dotnotation.New("a1"),
			FieldExample: appmodelcollections.Must(appmodelcollections.NewSwagValue(123)),
		},
		&appmodelcollections.ModelFieldSchema{
			FieldName:    dotnotation.New("a2", "b1"),
			FieldExample: appmodelcollections.Must(appmodelcollections.NewSwagValue("stringvalue")),
		},

		&appmodelcollections.ModelFieldSchema{
			FieldName:    dotnotation.New("a2", "b2"),
			FieldExample: appmodelcollections.Must(appmodelcollections.NewSwagValue(4345)),
		},

		&appmodelcollections.ModelFieldSchema{
			FieldName:    dotnotation.New("a3", "b1", "c1"),
			FieldExample: appmodelcollections.Must(appmodelcollections.NewSwagValue(123)),
		},

		&appmodelcollections.ModelFieldSchema{
			FieldName:    dotnotation.New("a3", "b1", "c2"),
			FieldExample: appmodelcollections.Must(appmodelcollections.NewSwagValue(456)),
		},
		&appmodelcollections.ModelFieldSchema{
			FieldName:    dotnotation.New("a3", "b1", "c3"),
			FieldExample: appmodelcollections.Must(appmodelcollections.NewSwagValue(789)),
		},
	}
	v, err := Build(fields)
	assert.NoError(t, err)
	assert.EqualValues(t, map[string]interface{}{
		"a1": float64(123),
		"a2": map[string]interface{}{
			"b1": "stringvalue",
			"b2": float64(4345),
		},
		"a3": map[string]interface{}{
			"b1": map[string]interface{}{
				"c1": float64(123),
				"c2": float64(456),
				"c3": float64(789),
			},
		},
	}, v)

}
