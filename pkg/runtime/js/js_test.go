package jsruntime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwagConvert(t *testing.T) {
	rt := NewJSRuntime("foo")
	defer rt.Close()
	assert.NoError(t, rt.Load("swagdef.js"))

	scriptInRaw := `doconvert('{"foo": "bar", "foofloat": 3.14}')`
	v, err := rt.Run(scriptInRaw)
	assert.NoError(t, err)
	assert.EqualValues(t, `"definitions": {
	"foo": {
		"type": "string",
		"example": "bar"
	},
	"foofloat": {
		"type": "number",
		"example": "3.14"
	}
}`, v)
}

func TestSwagNestedConvert(t *testing.T) {
	rt := NewJSRuntime("foo")
	defer rt.Close()
	assert.NoError(t, rt.Load("swagdef.js"))

	scriptInRaw := `doconvert('{"foo": "bar", "foofloat": 3.14, "nested": {"nestedfoo": 3.15, "nestedbar": "hello"}}')`
	v, err := rt.Run(scriptInRaw)
	assert.NoError(t, err)
	assert.EqualValues(t, `"definitions": {
	"foo": {
		"type": "string",
		"example": "bar"
	},
	"foofloat": {
		"type": "number",
		"example": "3.14"
	},
	"nested": {
		"type": "object",
		"properties": {
			"nestedfoo": {
				"type": "number",
				"example": "3.15"
			},
			"nestedbar": {
				"type": "string",
				"example": "hello"
			}
		}
	}
}`, v)
}
