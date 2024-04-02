package integrationtest

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/sdinsure/agent/pkg/logger"
	storagetestutils "github.com/sdinsure/agent/pkg/storage/testutils"
	"github.com/stretchr/testify/assert"

	restcolopenapicollections "github.com/footprintai/restcol/api/go-openapiv2/client/collections"
	restcolopenapidocument "github.com/footprintai/restcol/api/go-openapiv2/client/document"
	restcolopenapiswagger "github.com/footprintai/restcol/api/go-openapiv2/client/swagger"
	restcolopenapimodel "github.com/footprintai/restcol/api/go-openapiv2/models"
	integrationtestclient "github.com/footprintai/restcol/integrationtest/client"
	integrationtestserver "github.com/footprintai/restcol/integrationtest/server"
)

func TestIntegrationTest(t *testing.T) {
	log := logger.NewLogger()
	postgresDb, err := storagetestutils.NewTestPostgresCli(log)
	if err != nil {
		assert.NoError(t, err)
	}
	svr, err := integrationtestserver.NewServer(50050, 50051, postgresDb, log)
	if err != nil {
		log.Fatal("%+v", err)
	}

	fmt.Print("integrationtest about to start\n")
	svr.Start(1 * time.Second)
	defer svr.Stop(3 * time.Second)

	time.Sleep(1 * time.Second)

	client := integrationtestclient.MustNewClient("localhost:50051")

	// post /api/newdoc
	createDocumentParam := &restcolopenapidocument.RestColServiceCreateDocumentParams{
		Body: &restcolopenapimodel.APICreateDocumentRequest{
			Cid:        "", // empty cid, would create a new collection
			Did:        "",
			Pid:        "",  // empty pid, would use default pid
			Dataformat: nil, // use auto infer
			Data:       []byte(jsonData),
		},
	}
	restcolCreateDocumentOk, err := client.Document.RestColServiceCreateDocument(createDocumentParam, noAuthInfo())
	assert.NoError(t, err)
	//	createdDid := restcolCreateDocumentOk.Payload.Metadata.Did
	//	createdPid := restcolCreateDocumentOk.Payload.Metadata.Pid
	createdCid := restcolCreateDocumentOk.Payload.Metadata.Cid

	// get /api/collections/{cid}
	getCollectionParams := &restcolopenapicollections.RestColServiceGetCollectionParams{
		Cid: createdCid,
	}
	restcolGetCollectionOk, err := client.Collections.RestColServiceGetCollection(getCollectionParams, noAuthInfo())
	assert.NoError(t, err)
	expectedSchema := []*restcolopenapimodel.APISchemaField{
		&restcolopenapimodel.APISchemaField{
			Datatype: restcolopenapimodel.APISchemaFieldDataTypeSCHEMAFIELDDATATYPESTRING.Pointer(),
			Example: &restcolopenapimodel.APISchemaFieldExampleValue{
				StringValue: "bar",
			},
			Name: "foo",
		},
		&restcolopenapimodel.APISchemaField{
			Datatype: restcolopenapimodel.APISchemaFieldDataTypeSCHEMAFIELDDATATYPESTRING.Pointer(),
			Example: &restcolopenapimodel.APISchemaFieldExampleValue{
				StringValue: "foo2bar",
			},
			Name: "foo2.foo2foo",
		},
		&restcolopenapimodel.APISchemaField{
			Datatype: restcolopenapimodel.APISchemaFieldDataTypeSCHEMAFIELDDATATYPENUMBER.Pointer(),
			Example: &restcolopenapimodel.APISchemaFieldExampleValue{
				NumberValue: 123,
			},
			Name: "foo3.foo3foo",
		},
	}
	assert.EqualValues(t, restcolGetCollectionOk.Payload.Schemas, expectedSchema)

	// get /apidoc
	getSwaggerDocParams := &restcolopenapiswagger.RestColServiceGetSwaggerDocParams{}
	restcolGetSwaggerDocOk, err := client.Swagger.RestColServiceGetSwaggerDoc(getSwaggerDocParams, noAuthInfo(), integrationtestclient.WithRawSwaggerDocReader())
	assert.NoError(t, err)
	assert.EqualValues(t, restcolGetSwaggerDocOk.Payload.ContentType, "application/json")

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
