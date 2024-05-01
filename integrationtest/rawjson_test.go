package integrationtest

import (
	"encoding/json"
	"testing"

	restcolopenapidocument "github.com/footprintai/restcol/api/go-openapiv2/client/document"
	restcolopenapimodel "github.com/footprintai/restcol/api/go-openapiv2/models"
	"github.com/stretchr/testify/assert"
)

func TestRawJSONData(t *testing.T) {
	if testing.Short() {
		t.Skip("skip now")
		return
	}

	suite := SetupTest(t)
	defer suite.Close()

	client := suite.NewClient()

	// post /api/newdoc with raw json
	createDocumentParam := &restcolopenapidocument.RestColServiceCreateDocumentParams{
		Body: &restcolopenapimodel.APICreateDocumentRequest{
			CollectionID: "", // empty cid, would create a new collection
			DocumentID:   "",
			ProjectID:    "",  // empty pid, would use default pid
			Dataformat:   nil, // use auto infer
			Data:         []byte(makeRawJson()),
		},
	}
	restcolCreateDocumentOk, err := client.Document.RestColServiceCreateDocument(createDocumentParam, noAuthInfo())
	assert.NoError(t, err)
	assert.True(t, restcolCreateDocumentOk.Payload.Metadata.DocumentID != "")
}

func makeRawJson() []byte {
	v := map[string]interface{}{
		"stringfield": "foo",
		"intfield":    12345,
		"floatfield":  3.1415,
		"bytefields":  []byte("this is bytes string"),
		"nullfields":  nil,
		"arrayfields": []struct {
			FieldA string
			FieldB string
		}{{FieldA: "fieldvaluea", FieldB: "fieldvalueb"}},
	}
	rawJson, _ := json.Marshal(v)
	return rawJson
}
