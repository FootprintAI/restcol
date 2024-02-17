package swagdef

import (
	"testing"

	modelcollections "github.com/footprintai/restcol/pkg/models/collections"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestModelFieldsSchemaToSwagDef(t *testing.T) {
	mustValue := func(v any) modelcollections.SwagValueValue {
		return modelcollections.Must(
			modelcollections.NewSwagValue(v),
		)
	}
	fields := []modelcollections.ModelFieldSchema{
		modelcollections.ModelFieldSchema{
			FieldName:      "foo",
			FieldValueType: modelcollections.NumberSwagValueType,
			FieldExample:   mustValue(3.14),
		},
		modelcollections.ModelFieldSchema{
			FieldName:      "bar",
			FieldValueType: modelcollections.StringSwagValueType,
			FieldExample:   mustValue("3.14s"),
		},
		modelcollections.ModelFieldSchema{
			FieldName:      "foostruct.fieldstr",
			FieldValueType: modelcollections.StringSwagValueType,
			FieldExample:   mustValue("thisisstring"),
		},
		modelcollections.ModelFieldSchema{
			FieldName:      "foostruct.fieldfloat",
			FieldValueType: modelcollections.NumberSwagValueType,
			FieldExample:   mustValue(1.23),
		},
		modelcollections.ModelFieldSchema{
			FieldName:      "foostruct.fieldstruct.nestedfield1",
			FieldValueType: modelcollections.StringSwagValueType,
			FieldExample:   mustValue("this is nested struct"),
		},
	}
	gotDef, err := ModelFieldsSchemaToSwagDef(fields)
	assert.NoError(t, err)
	expectedDef := &spec.Definitions{
		"bar": spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray([]string{"string"}),
			},
			SwaggerSchemaProps: spec.SwaggerSchemaProps{
				Example: "3.14s",
			},
		},
		"foo": spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray([]string{"number"}),
			},
			SwaggerSchemaProps: spec.SwaggerSchemaProps{
				Example: "3.14", // in swagger example, everything is string
			},
		},
		"foostruct": spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray([]string{"object"}),
				Properties: spec.SchemaProperties{
					"fieldfloat": spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type: spec.StringOrArray([]string{"number"}),
						},
						SwaggerSchemaProps: spec.SwaggerSchemaProps{
							Example: "1.23", // in swagger example, everything is string
						},
					},
					"fieldstr": spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type: spec.StringOrArray([]string{"string"}),
						},
						SwaggerSchemaProps: spec.SwaggerSchemaProps{
							Example: "thisisstring",
						},
					},
					"fieldstruct": spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type: spec.StringOrArray([]string{"object"}),
							Properties: spec.SchemaProperties{
								"nestedfield1": spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: spec.StringOrArray([]string{"string"}),
									},
									SwaggerSchemaProps: spec.SwaggerSchemaProps{
										Example: "this is nested struct",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.EqualValues(t, expectedDef, gotDef)
}
