package integrationtest

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	restcolopenapiswagger "github.com/footprintai/restcol/api/go-openapiv2/client/swagger"
	integrationtestclient "github.com/footprintai/restcol/integrationtest/client"
)

func TestAPIDocTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skip now")
		return
	}

	suite := SetupTest(t)
	defer suite.Close()

	SetupCollection(t, suite)

	client := suite.NewClient()

	// get collection /apidoc/projects/${pid}/collections/${cid}
	getSwaggerDocParams := &restcolopenapiswagger.RestColServiceGetSwaggerDoc2Params{
		ProjectID:    projectId,
		CollectionID: cid,
	}
	restcolGetSwaggerDocOk, err := client.Swagger.RestColServiceGetSwaggerDoc2(getSwaggerDocParams, noAuthInfo(), integrationtestclient.WithRawSwaggerDocReader())
	assert.NoError(t, err)
	assert.EqualValues(t, restcolGetSwaggerDocOk.Payload.ContentType, "application/json")
	m := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal([]byte(restcolGetSwaggerDocOk.Payload.Data), &m))
	pathsShouldExists := []string{
		fmt.Sprintf("/v1/projects/%s/collections/%s/apidoc", projectId, cid),
		fmt.Sprintf("/v1/projects/%s/collections/%s:newdoc", projectId, cid),
		fmt.Sprintf("/v1/projects/%s/collections/%s/docs", projectId, cid),
		fmt.Sprintf("/v1/projects/%s/collections/%s/docs:stream", projectId, cid),
	}

	swaggerPaths := m["paths"].(map[string]interface{})
	//for p := range swaggerPaths {
	//	fmt.Printf("place:%s\n", p)
	//}

	for _, pathShouldExists := range pathsShouldExists {
		_, exist := swaggerPaths[pathShouldExists]
		assert.True(t, exist)
	}
}
