package integrationtest

import (
	"fmt"
	"net/http"
	"testing"

	sderrors "github.com/sdinsure/agent/pkg/errors"
	"github.com/stretchr/testify/assert"

	restcolopenapidocument "github.com/footprintai/restcol/api/go-openapiv2/client/document"
	restcolopenapimodel "github.com/footprintai/restcol/api/go-openapiv2/models"
)

// assertNotFound asserts that a call was refused because the document is not
// there, in the given (project, collection) scope — and that the refusal
// reaches the caller as an HTTP 404.
//
// The status half is the regression guard for #59. Until the sdinsure/agent
// bump, *every* application error left this service as a 500: the pinned
// version had no GRPCStatus() on *sderrors.Error, so NotFound degraded to
// Unknown on the way out through the gateway, and a client could not tell a
// missing document from a broken server. An earlier version of this helper
// deliberately asserted only the application code, so it would keep passing
// once the status was fixed. It is fixed; assert both.
func assertNotFound(t *testing.T, err error) {
	t.Helper()
	if !assert.Error(t, err) {
		return
	}
	want := fmt.Sprintf("code(%d),", sderrors.CodeNotFound.Int())
	assert.Contains(t, err.Error(), want,
		"expected a NotFound application error, got: %v", err)

	coder, ok := err.(interface{ Code() int })
	if !assert.True(t, ok, "expected a generated error carrying a status, got %T", err) {
		return
	}
	assert.Equal(t, http.StatusNotFound, coder.Code(),
		"NotFound must reach the client as 404; a 500 here means the error code "+
			"is not being mapped to a gRPC status (#59)")
}

func TestDeleteDocument(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}

	suite := SetupTest(t)
	defer suite.Close()

	SetupCollection(t, suite)

	client := suite.NewClient()

	createDocumentParam := &restcolopenapidocument.RestColServiceCreateDocument2Params{
		Body: &restcolopenapimodel.RestColServiceCreateDocumentBody{
			Data: []byte(`{"test": "document for deletion"}`),
		},
		CollectionID: cid,
		ProjectID:    projectId,
	}
	createResp, err := client.Document.RestColServiceCreateDocument2(createDocumentParam, noAuthInfo())
	assert.NoError(t, err)
	assert.NotNil(t, createResp.Payload.Metadata)
	documentId := createResp.Payload.Metadata.DocumentID

	getDocumentParam := &restcolopenapidocument.RestColServiceGetDocumentParams{
		CollectionID: cid,
		DocumentID:   documentId,
		ProjectID:    projectId,
	}
	getResp, err := client.Document.RestColServiceGetDocument(getDocumentParam, noAuthInfo())
	assert.NoError(t, err)
	assert.Equal(t, documentId, getResp.Payload.Metadata.DocumentID)

	deleteDocumentParam := &restcolopenapidocument.RestColServiceDeleteDocumentParams{
		CollectionID: cid,
		DocumentID:   documentId,
		ProjectID:    projectId,
	}
	deleteResp, err := client.Document.RestColServiceDeleteDocument(deleteDocumentParam, noAuthInfo())
	assert.NoError(t, err)
	assert.NotNil(t, deleteResp.Payload)

	// The whole point of the feature (grandturks#936): the document is gone, and
	// the API says NotFound rather than answering 200 with blank metadata.
	_, err = client.Document.RestColServiceGetDocument(getDocumentParam, noAuthInfo())
	assertNotFound(t, err)
}

func TestDeleteDocument_NonExistent(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}

	suite := SetupTest(t)
	defer suite.Close()

	SetupCollection(t, suite)

	client := suite.NewClient()

	// Deleting something that was never there is a 404, not a silent success.
	// A caller that gets 200 back has no way to distinguish "removed it" from
	// "removed nothing", which is the failure mode this endpoint shipped with.
	deleteDocumentParam := &restcolopenapidocument.RestColServiceDeleteDocumentParams{
		CollectionID: cid,
		DocumentID:   "non-existent-document-id",
		ProjectID:    projectId,
	}
	_, err := client.Document.RestColServiceDeleteDocument(deleteDocumentParam, noAuthInfo())
	assertNotFound(t, err)
}

func TestDeleteDocument_InvalidParameters(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}

	suite := SetupTest(t)
	defer suite.Close()

	SetupCollection(t, suite)

	client := suite.NewClient()

	// Empty collection id
	deleteDocumentParam := &restcolopenapidocument.RestColServiceDeleteDocumentParams{
		CollectionID: "",
		DocumentID:   "some-document-id",
		ProjectID:    projectId,
	}
	_, err := client.Document.RestColServiceDeleteDocument(deleteDocumentParam, noAuthInfo())
	assert.Error(t, err)

	// Empty document id
	deleteDocumentParam = &restcolopenapidocument.RestColServiceDeleteDocumentParams{
		CollectionID: cid,
		DocumentID:   "",
		ProjectID:    projectId,
	}
	_, err = client.Document.RestColServiceDeleteDocument(deleteDocumentParam, noAuthInfo())
	assert.Error(t, err)
}

// TestDeleteDocument_UnknownProjectIsRejected pins where tenancy is actually
// decided, which is not where the handler code suggests.
//
// The handler never reads req.ProjectId; it calls getProjectIdFromCtx. What
// fills that context is the project identity middleware, and *that* reads the
// {projectId} path parameter and resolves it against the projects table. So the
// path segment is load-bearing after all: an id that resolves to no project is
// refused before the handler runs, with BadParameters rather than NotFound.
//
// Scoping between two *existing* projects cannot be exercised here — this suite
// boots a single bootstrap tenant — and is covered at the handler level in
// pkg/app/handlers_test.go, where the resolver can be varied.
func TestDeleteDocument_UnknownProjectIsRejected(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}

	suite := SetupTest(t)
	defer suite.Close()

	SetupCollection(t, suite)

	client := suite.NewClient()

	createResp, err := client.Document.RestColServiceCreateDocument2(
		&restcolopenapidocument.RestColServiceCreateDocument2Params{
			Body: &restcolopenapimodel.RestColServiceCreateDocumentBody{
				Data: []byte(`{"test": "path project id is ignored"}`),
			},
			CollectionID: cid,
			ProjectID:    projectId,
		}, noAuthInfo())
	assert.NoError(t, err)
	documentId := createResp.Payload.Metadata.DocumentID

	_, err = client.Document.RestColServiceDeleteDocument(
		&restcolopenapidocument.RestColServiceDeleteDocumentParams{
			CollectionID: cid,
			DocumentID:   documentId,
			ProjectID:    "wrong-project-id",
		}, noAuthInfo())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("code(%d),", sderrors.CodeBadParameters.Int()),
		"an unresolvable project must be rejected by the identity middleware, got: %v", err)
	if coder, ok := err.(interface{ Code() int }); assert.True(t, ok) {
		// The other half of the #59 mapping: BadParameters must arrive as 400,
		// not as the blanket 500 every application error used to become.
		assert.Equal(t, http.StatusBadRequest, coder.Code(),
			"BadParameters must reach the client as 400, got: %v", err)
	}

	// ... and the document is untouched under its real project.
	getResp, err := client.Document.RestColServiceGetDocument(
		&restcolopenapidocument.RestColServiceGetDocumentParams{
			CollectionID: cid,
			DocumentID:   documentId,
			ProjectID:    projectId,
		}, noAuthInfo())
	assert.NoError(t, err)
	assert.Equal(t, documentId, getResp.Payload.Metadata.DocumentID)
}
