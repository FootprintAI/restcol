package schema

import (
	"strings"
	"testing"

	"github.com/sdinsure/agent/pkg/logger"
	"github.com/stretchr/testify/assert"

	appmodelcollections "github.com/footprintai/restcol/pkg/models/collections"
	dotnotation "github.com/footprintai/restcol/pkg/notation/dot"
)

func TestTraverseMap(t *testing.T) {

	traversedMap := make(map[string]interface{})
	TraverseMap(logger.NewLogger(), m1, []string{}, func(prefixes []string, current string, val any) error {
		path := strings.Join(append(prefixes, current), ".")
		traversedMap[path] = val
		return nil
	})

	expectedMap := map[string]interface{}{
		"foo":                           "bar",
		"fooInt":                        64,
		"foostruct.foostructslice":      []string{"val1", "val2", "val3"},
		"foostruct.foostructslicefloat": []float64{6.14, 12.245},
		"foostruct.foostructstring":     "bar",
		"foostruct.foostructstruct.foostructstructstring": "bar",
	}
	assert.EqualValues(t, expectedMap, traversedMap)

}

var (
	m1 = map[string]interface{}{
		"foo":    "bar",
		"fooInt": 64,
		"foostruct": map[string]interface{}{
			"foostructstring": "bar",
			"foostructslice": []string{
				"val1",
				"val2",
				"val3",
			},
			"foostructslicefloat": []float64{
				6.14,
				12.245,
			},
			"foostructstruct": map[string]interface{}{
				"foostructstructstring": "bar",
			},
		},
	}
)

func TestSchemaFlatten(t *testing.T) {
	schemaBuilder := NewSchemaBuilder(logger.NewLogger())

	m1Schema, err := schemaBuilder.Flatten(m1)
	assert.NoError(t, err)

	m1Schema2 := &appmodelcollections.ModelSchema{
		Fields: []*appmodelcollections.ModelFieldSchema{
			&appmodelcollections.ModelFieldSchema{
				FieldName:      dotnotation.New("foo"),
				FieldValueType: appmodelcollections.StringSwagValueType,
			},
			&appmodelcollections.ModelFieldSchema{
				FieldName:      dotnotation.New("fooInt"),
				FieldValueType: appmodelcollections.NumberSwagValueType,
			},
			&appmodelcollections.ModelFieldSchema{
				FieldName:      dotnotation.New("foostruct", "foostructslice"),
				FieldValueType: appmodelcollections.ArraySwagValueType,
			},
			&appmodelcollections.ModelFieldSchema{
				FieldName:      dotnotation.New("foostruct", "foostructslicefloat"),
				FieldValueType: appmodelcollections.ArraySwagValueType,
			},
			&appmodelcollections.ModelFieldSchema{
				FieldName:      dotnotation.New("foostruct", "foostructstring"),
				FieldValueType: appmodelcollections.StringSwagValueType,
			},
			&appmodelcollections.ModelFieldSchema{
				FieldName:      dotnotation.New("foostruct", "foostructstruct", "foostructstructstring"),
				FieldValueType: appmodelcollections.StringSwagValueType,
			},
		},
	}
	equal := schemaBuilder.Equals(m1Schema, m1Schema2)
	assert.True(t, equal)
}
