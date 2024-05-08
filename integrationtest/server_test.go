package integrationtest

import (
	"encoding/json"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"

	restcolopenapicollections "github.com/footprintai/restcol/api/go-openapiv2/client/collections"
	restcolopenapidocument "github.com/footprintai/restcol/api/go-openapiv2/client/document"
	restcolopenapimodel "github.com/footprintai/restcol/api/go-openapiv2/models"
)

func TestIntegrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skip now")
		return
	}

	suite := SetupTest(t)
	defer suite.Close()

	SetupCollection(t, suite)

	client := suite.NewClient()
	// post /api/newdoc
	createDocumentParam := &restcolopenapidocument.RestColServiceCreateDocument2Params{
		Body: &restcolopenapimodel.RestColServiceCreateDocumentBody{
			Data: []byte(jsonData),
		},
		CollectionID: cid,
		ProjectID:    projectId,
	}
	restcolCreateDocumentOk, err := client.Document.RestColServiceCreateDocument2(createDocumentParam, noAuthInfo())
	assert.NoError(t, err)

	createDocumentParam2 := &restcolopenapidocument.RestColServiceCreateDocument2Params{
		Body: &restcolopenapimodel.RestColServiceCreateDocumentBody{
			Data: []byte(jsonData),
		},
		CollectionID: cid,
		ProjectID:    projectId,
	}
	restcolCreateDocumentOk2, err := client.Document.RestColServiceCreateDocument2(createDocumentParam2, noAuthInfo())
	assert.NoError(t, err)

	// make sure the docs are not the same
	assert.True(t, restcolCreateDocumentOk2.Payload.Metadata.DocumentID != restcolCreateDocumentOk.Payload.Metadata.DocumentID)
	// make sure the schema id is the same as it can be reused
	assert.True(t, restcolCreateDocumentOk2.Payload.Metadata.SchemaID == restcolCreateDocumentOk.Payload.Metadata.SchemaID)
	assert.True(t, restcolCreateDocumentOk2.Payload.Metadata.SchemaID != "")

	// get /api/collections/{cid}
	getCollectionParams := &restcolopenapicollections.RestColServiceGetCollectionParams{
		CollectionID: cid,
		ProjectID:    projectId,
	}
	restcolGetCollectionOk, err := client.Collections.RestColServiceGetCollection(getCollectionParams, noAuthInfo())
	assert.NoError(t, err)
	expectedSchema := []*restcolopenapimodel.APISchemaField{
		&restcolopenapimodel.APISchemaField{
			Datatype: restcolopenapimodel.APISchemaFieldDataTypeSCHEMAFIELDDATATYPESTRING.Pointer(),
			Example:  "bar",
			Name:     "foo",
		},
		&restcolopenapimodel.APISchemaField{
			Datatype: restcolopenapimodel.APISchemaFieldDataTypeSCHEMAFIELDDATATYPESTRING.Pointer(),
			Example:  "foo2bar",
			Name:     "foo2.foo2foo",
		},
		&restcolopenapimodel.APISchemaField{
			Datatype: restcolopenapimodel.APISchemaFieldDataTypeSCHEMAFIELDDATATYPENUMBER.Pointer(),
			Example:  json.Number("123"),
			Name:     "foo3.foo3foo",
		},
	}
	assert.EqualValues(t, expectedSchema, restcolGetCollectionOk.Payload.Schemas)

	// get doc with field selector
	getDocumentParam := &restcolopenapidocument.RestColServiceGetDocumentParams{
		DocumentID:     restcolCreateDocumentOk.Payload.Metadata.DocumentID,
		CollectionID:   cid,
		ProjectID:      projectId,
		FieldSelectors: []string{"foo2.foo2foo"},
	}
	restcolGetDocumentOk, err := client.Document.RestColServiceGetDocument(getDocumentParam, noAuthInfo())
	assert.NoError(t, err)
	assert.EqualValues(t, map[string]interface{}{
		"foo2": map[string]interface{}{
			"foo2foo": "foo2bar",
		},
	}, restcolGetDocumentOk.Payload.Data)

	// get doc with empty field selector
	getDocumentParam = &restcolopenapidocument.RestColServiceGetDocumentParams{
		DocumentID:   restcolCreateDocumentOk.Payload.Metadata.DocumentID,
		CollectionID: cid,
		ProjectID:    projectId,
	}
	restcolGetDocumentOk, err = client.Document.RestColServiceGetDocument(getDocumentParam, noAuthInfo())
	assert.NoError(t, err)
	assert.EqualValues(t, map[string]interface{}{
		"foo": "bar",
		"foo2": map[string]interface{}{
			"foo2foo": "foo2bar",
		},
		"foo3": map[string]interface{}{
			"foo3foo": json.Number("123"),
		},
	}, restcolGetDocumentOk.Payload.Data)
}

func mustJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

var (
	jsonData = `{"foo":"bar", "foo2": {"foo2foo": "foo2bar"}, "foo3": {"foo3foo": 123}}`
)

func noAuthInfo() runtime.ClientAuthInfoWriterFunc {
	return runtime.ClientAuthInfoWriterFunc(func(clientRequest runtime.ClientRequest, registry strfmt.Registry) error {
		// TODO: should set header if we want to include auth in the requst
		return nil
	})
}
