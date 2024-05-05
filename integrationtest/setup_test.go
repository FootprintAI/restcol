package integrationtest

import (
	"fmt"
	"testing"

	restcolopenapicollections "github.com/footprintai/restcol/api/go-openapiv2/client/collections"
	restcolopenapimodels "github.com/footprintai/restcol/api/go-openapiv2/models"
	restcoldummy "github.com/footprintai/restcol/pkg/dummy"
	"github.com/stretchr/testify/assert"
)

var (
	projectId = restcoldummy.DummyModelProject.ID.String()
	cid       string
)

func SetupCollection(t *testing.T, s *suite) {
	if cid == "" {
		fmt.Printf("create collection\n")

		client := s.NewClient()
		// post create test collection
		createCollectionParam := &restcolopenapicollections.RestColServiceCreateCollectionParams{
			Body: &restcolopenapimodels.RestColServiceCreateCollectionBody{
				Description: "unittest",
			},
			ProjectID: projectId,
		}
		restColServiceCreateCollectionOK, err := client.Collections.RestColServiceCreateCollection(createCollectionParam, noAuthInfo())
		assert.NoError(t, err)
		cid = restColServiceCreateCollectionOK.Payload.Metadata.CollectionID
	}
	fmt.Printf("use collection id:%+v\n", cid)
}
