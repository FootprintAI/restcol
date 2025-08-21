package integrationtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	restcolopenapidocument "github.com/footprintai/restcol/api/go-openapiv2/client/document"
	restcolopenapimodel "github.com/footprintai/restcol/api/go-openapiv2/models"
)

func TestDeleteDocument(t *testing.T) {
	if testing.Short() {
		t.Skip("skip now")
		return
	}

	suite := SetupTest(t)
	defer suite.Close()

	SetupCollection(t, suite)

	client := suite.NewClient()

	// Create a document to delete
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

	// Verify document exists by getting it
	getDocumentParam := &restcolopenapidocument.RestColServiceGetDocumentParams{
		CollectionID: cid,
		DocumentID:   documentId,
		ProjectID:    projectId,
	}
	getResp, err := client.Document.RestColServiceGetDocument(getDocumentParam, noAuthInfo())
	assert.NoError(t, err)
	assert.Equal(t, documentId, getResp.Payload.Metadata.DocumentID)

	// Delete the document
	deleteDocumentParam := &restcolopenapidocument.RestColServiceDeleteDocumentParams{
		CollectionID: cid,
		DocumentID:   documentId,
		ProjectID:    projectId,
	}
	deleteResp, err := client.Document.RestColServiceDeleteDocument(deleteDocumentParam, noAuthInfo())
	assert.NoError(t, err)
	assert.NotNil(t, deleteResp.Payload) // Should return empty response

	// Verify document is deleted (should return empty response)
	getResp2, err := client.Document.RestColServiceGetDocument(getDocumentParam, noAuthInfo())
	assert.NoError(t, err) // GetDocument doesn't error for soft-deleted records
	assert.Empty(t, getResp2.Payload.Metadata.DocumentID) // Should return empty metadata
}

func TestDeleteDocument_NonExistent(t *testing.T) {
	if testing.Short() {
		t.Skip("skip now")
		return
	}

	suite := SetupTest(t)
	defer suite.Close()

	SetupCollection(t, suite)

	client := suite.NewClient()

	// Try to delete a non-existent document (should succeed due to idempotent behavior)
	deleteDocumentParam := &restcolopenapidocument.RestColServiceDeleteDocumentParams{
		CollectionID: cid,
		DocumentID:   "non-existent-document-id",
		ProjectID:    projectId,
	}
	deleteResp, err := client.Document.RestColServiceDeleteDocument(deleteDocumentParam, noAuthInfo())
	assert.NoError(t, err) // DELETE is idempotent - should succeed even for non-existent documents
	assert.NotNil(t, deleteResp.Payload) // Should return empty response
}

func TestDeleteDocument_InvalidParameters(t *testing.T) {
	if testing.Short() {
		t.Skip("skip now")
		return
	}

	suite := SetupTest(t)
	defer suite.Close()

	SetupCollection(t, suite)

	client := suite.NewClient()

	// Test with empty collection ID
	deleteDocumentParam := &restcolopenapidocument.RestColServiceDeleteDocumentParams{
		CollectionID: "", // Empty collection ID
		DocumentID:   "some-document-id",
		ProjectID:    projectId,
	}
	_, err := client.Document.RestColServiceDeleteDocument(deleteDocumentParam, noAuthInfo())
	assert.Error(t, err) // Should return validation error

	// Test with empty document ID
	deleteDocumentParam = &restcolopenapidocument.RestColServiceDeleteDocumentParams{
		CollectionID: cid,
		DocumentID:   "", // Empty document ID
		ProjectID:    projectId,
	}
	_, err = client.Document.RestColServiceDeleteDocument(deleteDocumentParam, noAuthInfo())
	assert.Error(t, err) // Should return validation error
}

func TestDeleteDocument_CrossProject(t *testing.T) {
	if testing.Short() {
		t.Skip("skip now")
		return
	}

	suite := SetupTest(t)
	defer suite.Close()

	SetupCollection(t, suite)

	client := suite.NewClient()

	// Create a document
	createDocumentParam := &restcolopenapidocument.RestColServiceCreateDocument2Params{
		Body: &restcolopenapimodel.RestColServiceCreateDocumentBody{
			Data: []byte(`{"test": "document for cross-project deletion test"}`),
		},
		CollectionID: cid,
		ProjectID:    projectId,
	}
	createResp, err := client.Document.RestColServiceCreateDocument2(createDocumentParam, noAuthInfo())
	assert.NoError(t, err)
	documentId := createResp.Payload.Metadata.DocumentID

	// Try to delete with wrong project ID (should succeed but not affect the actual document)
	deleteDocumentParam := &restcolopenapidocument.RestColServiceDeleteDocumentParams{
		CollectionID: cid,
		DocumentID:   documentId,
		ProjectID:    "wrong-project-id", // Wrong project ID
	}
	_, err = client.Document.RestColServiceDeleteDocument(deleteDocumentParam, noAuthInfo())
	assert.NoError(t, err) // DELETE succeeds (0 rows affected)

	// Verify document still exists with correct project ID
	getDocumentParam := &restcolopenapidocument.RestColServiceGetDocumentParams{
		CollectionID: cid,
		DocumentID:   documentId,
		ProjectID:    projectId,
	}
	getResp, err := client.Document.RestColServiceGetDocument(getDocumentParam, noAuthInfo())
	assert.NoError(t, err)
	assert.Equal(t, documentId, getResp.Payload.Metadata.DocumentID)
}