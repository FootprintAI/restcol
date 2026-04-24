package app

import (
	"context"
	"testing"

	sderrors "github.com/sdinsure/agent/pkg/errors"
	"github.com/sdinsure/agent/pkg/logger"
	sdinsureruntime "github.com/sdinsure/agent/pkg/runtime"
	storagetestutils "github.com/sdinsure/agent/pkg/storage/testutils"
	"github.com/stretchr/testify/assert"

	apppb "github.com/footprintai/restcol/api/pb"
	collectionsmodel "github.com/footprintai/restcol/pkg/models/collections"
	documentsmodel "github.com/footprintai/restcol/pkg/models/documents"
	projectsmodel "github.com/footprintai/restcol/pkg/models/projects"
	schemafinder "github.com/footprintai/restcol/pkg/schema"
	collectionsstorage "github.com/footprintai/restcol/pkg/storage/collections"
	documentsstorage "github.com/footprintai/restcol/pkg/storage/documents"
	projectsstorage "github.com/footprintai/restcol/pkg/storage/projects"
)

// fixedProjectResolver implements sdinsureruntime.ProjectResolver for tests:
// every call returns the configured project, so handlers see a stable tenant.
type fixedProjectResolver struct {
	pid projectsmodel.ProjectID
}

func (f *fixedProjectResolver) WithProjectInfo(ctx context.Context, _ string) context.Context {
	return ctx
}

func (f *fixedProjectResolver) ProjectInfo(_ context.Context) (sdinsureruntime.ProjectInfor, bool) {
	return &fixedProjectInfor{pid: f.pid}, true
}

type fixedProjectInfor struct {
	pid projectsmodel.ProjectID
}

func (f *fixedProjectInfor) GetProjectID() (string, error) { return f.pid.String(), nil }
func (f *fixedProjectInfor) GetProject(any) error          { return nil }

// newTestService spins up a full handler against a real postgres. Tests using
// it must be guarded with testing.Short() — CI runs -short.
func newTestService(t *testing.T) (*RestColServiceServerService, *projectsmodel.ModelProject) {
	t.Helper()
	log := logger.NewLogger()
	db, err := storagetestutils.NewTestPostgresCli(log)
	assert.NoError(t, err)

	projectCURD := projectsstorage.NewProjectCURD(db)
	assert.NoError(t, projectCURD.AutoMigrate())
	collectionCURD := collectionsstorage.NewCollectionCURD(db)
	assert.NoError(t, collectionCURD.AutoMigrate())
	documentCURD := documentsstorage.NewDocumentCURD(db)
	assert.NoError(t, documentCURD.AutoMigrate())

	proj := &projectsmodel.ModelProject{
		ID:   projectsmodel.NewProjectID(9001),
		Type: projectsmodel.RegularProjectType,
	}
	assert.NoError(t, projectCURD.Write(context.Background(), "", proj))

	svc := NewRestColServiceServerService(log, collectionCURD, documentCURD, schemafinder.NewSchemaBuilder(log))
	svc.SetDefaultProjectResolver(&fixedProjectResolver{pid: proj.ID})
	return svc, proj
}

func assertErrorCode(t *testing.T, err error, want sderrors.Code) {
	t.Helper()
	if !assert.Error(t, err) {
		return
	}
	ok, myerr := sderrors.As(err)
	if !assert.True(t, ok, "expected sderrors.Error, got %T: %v", err, err) {
		return
	}
	assert.Equal(t, want, myerr.Code(), "error code mismatch; message: %s", myerr.Error())
}

func TestDeleteCollection_ValidationAndExistence(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)
	ctx := context.Background()

	_, err := svc.DeleteCollection(ctx, &apppb.DeleteCollectionRequest{CollectionId: ""})
	assertErrorCode(t, err, sderrors.CodeBadParameters)

	_, err = svc.DeleteCollection(ctx, &apppb.DeleteCollectionRequest{
		CollectionId: collectionsmodel.NewCollectionID().String(),
	})
	assertErrorCode(t, err, sderrors.CodeNotFound)
}

func TestDeleteCollection_EmptyCollectionSucceeds(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)
	ctx := context.Background()

	created, err := svc.CreateCollection(ctx, &apppb.CreateCollectionRequest{
		Description: strPtr("empty-coll"),
	})
	assert.NoError(t, err)
	cid := created.XMetadata.CollectionId

	_, err = svc.DeleteCollection(ctx, &apppb.DeleteCollectionRequest{CollectionId: cid})
	assert.NoError(t, err)
}

func TestDeleteCollection_NonEmptyRejectedWithoutForce(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)
	ctx := context.Background()

	created, err := svc.CreateCollection(ctx, &apppb.CreateCollectionRequest{Description: strPtr("non-empty")})
	assert.NoError(t, err)
	cid := created.XMetadata.CollectionId

	_, err = svc.CreateDocument(ctx, &apppb.CreateDocumentRequest{
		CollectionId: cid,
		Data:         []byte(`{"k":"v"}`),
	})
	assert.NoError(t, err)

	_, err = svc.DeleteCollection(ctx, &apppb.DeleteCollectionRequest{CollectionId: cid, Force: false})
	assertErrorCode(t, err, sderrors.CodeStatusConflicted)
}

func TestDeleteCollection_ForceCascades(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)
	ctx := context.Background()

	created, err := svc.CreateCollection(ctx, &apppb.CreateCollectionRequest{Description: strPtr("cascade")})
	assert.NoError(t, err)
	cid := created.XMetadata.CollectionId

	createdDoc, err := svc.CreateDocument(ctx, &apppb.CreateDocumentRequest{
		CollectionId: cid,
		Data:         []byte(`{"k":"v"}`),
	})
	assert.NoError(t, err)
	did := createdDoc.XMetadata.DocumentId

	_, err = svc.DeleteCollection(ctx, &apppb.DeleteCollectionRequest{CollectionId: cid, Force: true})
	assert.NoError(t, err)

	// Document should no longer be retrievable after cascade.
	_, err = svc.GetDocument(ctx, &apppb.GetDocumentRequest{CollectionId: cid, DocumentId: did})
	assertErrorCode(t, err, sderrors.CodeNotFound)
}

func TestGetDocument_ValidationAndScope(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)
	ctx := context.Background()

	cases := []struct {
		name string
		req  *apppb.GetDocumentRequest
		code sderrors.Code
	}{
		{
			name: "missing collection id",
			req:  &apppb.GetDocumentRequest{DocumentId: documentsmodel.NewDocumentID().String()},
			code: sderrors.CodeBadParameters,
		},
		{
			name: "missing document id",
			req:  &apppb.GetDocumentRequest{CollectionId: collectionsmodel.NewCollectionID().String()},
			code: sderrors.CodeBadParameters,
		},
		{
			name: "malformed document id",
			req: &apppb.GetDocumentRequest{
				CollectionId: collectionsmodel.NewCollectionID().String(),
				DocumentId:   "not-a-valid-doc-id",
			},
			code: sderrors.CodeBadParameters,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.GetDocument(ctx, tc.req)
			assertErrorCode(t, err, tc.code)
		})
	}
}

func TestGetDocument_WrongCollectionScopeReturnsNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)
	ctx := context.Background()

	// Create two collections and a document in the first.
	c1, err := svc.CreateCollection(ctx, &apppb.CreateCollectionRequest{Description: strPtr("c1")})
	assert.NoError(t, err)
	c2, err := svc.CreateCollection(ctx, &apppb.CreateCollectionRequest{Description: strPtr("c2")})
	assert.NoError(t, err)

	doc, err := svc.CreateDocument(ctx, &apppb.CreateDocumentRequest{
		CollectionId: c1.XMetadata.CollectionId,
		Data:         []byte(`{"k":"v"}`),
	})
	assert.NoError(t, err)

	// Reading the doc under the wrong collection must NOT leak an empty
	// response — it must surface as NotFound.
	_, err = svc.GetDocument(ctx, &apppb.GetDocumentRequest{
		CollectionId: c2.XMetadata.CollectionId,
		DocumentId:   doc.XMetadata.DocumentId,
	})
	assertErrorCode(t, err, sderrors.CodeNotFound)

	// And the happy path still works.
	got, err := svc.GetDocument(ctx, &apppb.GetDocumentRequest{
		CollectionId: c1.XMetadata.CollectionId,
		DocumentId:   doc.XMetadata.DocumentId,
	})
	assert.NoError(t, err)
	assert.Equal(t, doc.XMetadata.DocumentId, got.XMetadata.DocumentId)
}

func TestDeleteDocument_Validation(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)
	ctx := context.Background()

	_, err := svc.DeleteDocument(ctx, &apppb.DeleteDocumentRequest{
		CollectionId: "",
		DocumentId:   documentsmodel.NewDocumentID().String(),
	})
	assertErrorCode(t, err, sderrors.CodeBadParameters)

	_, err = svc.DeleteDocument(ctx, &apppb.DeleteDocumentRequest{
		CollectionId: collectionsmodel.NewCollectionID().String(),
		DocumentId:   "",
	})
	assertErrorCode(t, err, sderrors.CodeBadParameters)

	_, err = svc.DeleteDocument(ctx, &apppb.DeleteDocumentRequest{
		CollectionId: collectionsmodel.NewCollectionID().String(),
		DocumentId:   "not-a-valid-doc-id",
	})
	assertErrorCode(t, err, sderrors.CodeBadParameters)
}

func TestGetCollection_MissingIDRejected(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)
	ctx := context.Background()

	_, err := svc.GetCollection(ctx, &apppb.GetCollectionRequest{CollectionId: ""})
	assertErrorCode(t, err, sderrors.CodeBadParameters)
}

func strPtr(s string) *string { return &s }
