package integrationtest

import (
	"context"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	sdinsureopenapistream "github.com/sdinsure/agent/pkg/http/openapi/stream"
	"github.com/sdinsure/agent/pkg/logger"
	"github.com/stretchr/testify/assert"

	restcolopenapidocument "github.com/footprintai/restcol/api/go-openapiv2/client/document"
	restcolopenapimodel "github.com/footprintai/restcol/api/go-openapiv2/models"
	apihelper "github.com/footprintai/restcol/api/helper"
	"github.com/footprintai/restcol/pkg/nullable"
)

func TestIntegrationStreamTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skip now")
		return
	}

	suite := SetupTest(t)
	defer suite.Close()

	log := logger.NewLogger()
	client := suite.NewClient()

	// post /api/newdoc
	createDocumentParam := &restcolopenapidocument.RestColServiceCreateDocumentParams{
		Body: &restcolopenapimodel.APICreateDocumentRequest{
			Data: []byte(jsonData),
		},
	}
	restcolCreateDocumentOk, err := client.Document.RestColServiceCreateDocument(createDocumentParam, noAuthInfo())
	assert.NoError(t, err)
	createdPid := restcolCreateDocumentOk.Payload.Metadata.ProjectID
	createdCid := restcolCreateDocumentOk.Payload.Metadata.CollectionID

	startedTs := strfmt.DateTime(time.Now())

	// create a goroutine to continueous write docs
	go func(count int32) {
		for i := int32(0); i < count; i++ {
			createDocumentParam := &restcolopenapidocument.RestColServiceCreateDocumentParams{
				Body: &restcolopenapimodel.APICreateDocumentRequest{
					Data:         []byte(jsonData),
					CollectionID: createdCid,
					ProjectID:    createdPid,
				},
			}
			_, err := client.Document.RestColServiceCreateDocument(createDocumentParam, noAuthInfo())
			assert.NoError(t, err)

			time.Sleep(10 * time.Millisecond)
		}
	}(100)
	cctx, cancel := context.WithCancel(context.Background())

	// create a stream to read the data,
	queryDocumentsStreamParam := &restcolopenapidocument.RestColServiceQueryDocumentsStreamParams{
		ProjectID:    nullable.StringP(createdPid),
		CollectionID: nullable.StringP(createdCid),
		SinceTs:      &startedTs,
		FollowUpMode: nullable.BoolP(true),
	}
	queryDocumentsStreamParam.WithContext(cctx)
	sinkCloser := apihelper.NewDocumentStreamSinkCloser()

	go func() {
		// NOTE: this would never ended blocking call until stream is ended
		_, err = client.Document.RestColServiceQueryDocumentsStream(queryDocumentsStreamParam, noAuthInfo(), sdinsureopenapistream.StreamReceiveClientOption(cctx, log, sinkCloser))
		assert.NoError(t, err)
	}()

	recvDocIDMap := make(map[string]struct{})
	for {
		queryDocumentsStreamOk, err := apihelper.WithError(sinkCloser.Recv())
		if err != nil {
			if err == apihelper.ErrEOF {
				cancel()
				break
			}
			assert.NoError(t, err)
		}
		if queryDocumentsStreamOk != nil && queryDocumentsStreamOk.Payload != nil && queryDocumentsStreamOk.Payload.Result != nil {
			recvDocIDMap[queryDocumentsStreamOk.Payload.Result.Metadata.DocumentID] = struct{}{}
		}

		if len(recvDocIDMap) == 100 {
			cancel()
			break
		}
	}

	assert.EqualValues(t, len(recvDocIDMap), 100)
}
