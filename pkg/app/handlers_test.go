package app

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	storagetestutils "github.com/footprintai/restcol/pkg/storage/testutil"
	sderrors "github.com/sdinsure/agent/pkg/errors"
	"github.com/sdinsure/agent/pkg/logger"
	sdinsureruntime "github.com/sdinsure/agent/pkg/runtime"
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

// Visibility matches the production getter, which reports "" because restcol
// has no per-project visibility. See pkg/runtime/getter.projectInfor.
func (f *fixedProjectInfor) Visibility() string { return "" }

// newTestService spins up a full handler against a real postgres. Tests using
// it must be guarded with testing.Short() — CI runs -short.
// Handed out by newTestService so each test gets its own tenant. See the
// comment there for why sharing one was a problem.
//
// Seeded from the clock rather than counting from a fixed base. A plain
// counter isolates tests WITHIN a run and not ACROSS runs: the database is
// never reset, so the second `go test` reuses the first's project ids and
// reads back rows the earlier run left behind. That is not hypothetical - it
// made a passing test fail on its next invocation with no code change between
// them, which reads as flakiness rather than as leftover state.
var nextTestProject = func() *atomic.Int64 {
	v := &atomic.Int64{}
	// Seconds, not nanoseconds: ProjectID is rendered with %d and the ids stay
	// readable in a failure message. Two runs in the same second would collide,
	// so leave room between runs.
	v.Store(time.Now().Unix() % 1_000_000 * 1000)
	return v
}()

func newTestService(t *testing.T) (*RestColServiceServerService, *projectsmodel.ModelProject) {
	t.Helper()
	log := logger.NewLogger(false)
	db, err := storagetestutils.NewTestPostgresCli(log)
	assert.NoError(t, err)

	projectCURD := projectsstorage.NewProjectCURD(db)
	assert.NoError(t, projectCURD.AutoMigrate())
	collectionCURD := collectionsstorage.NewCollectionCURD(db)
	assert.NoError(t, collectionCURD.AutoMigrate())
	documentCURD := documentsstorage.NewDocumentCURD(db)
	assert.NoError(t, documentCURD.AutoMigrate())

	// A DISTINCT project per test, not a shared 9001.
	//
	// Every test here talks to the same postgres, and the database is not reset
	// between them. With one shared project, whatever a test creates is still
	// there for the next one - fine for assertions that name a specific id, and
	// wrong for any assertion about how many collections a project has. The
	// list tests below are exactly that kind, and they failed on rows their own
	// test never created.
	//
	proj := &projectsmodel.ModelProject{
		ID:   projectsmodel.NewProjectID(int(nextTestProject.Add(1))),
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
			// Document ids are opaque strings, not UUIDs: documentsmodel.Parse
			// accepts anything and CreateDocument will happily store a
			// caller-supplied "not-a-valid-doc-id". Rejecting the same string on
			// read would make a document that can be created but never fetched,
			// so an unrecognised id is NotFound rather than BadParameters.
			name: "unrecognised document id",
			req: &apppb.GetDocumentRequest{
				CollectionId: collectionsmodel.NewCollectionID().String(),
				DocumentId:   "not-a-valid-doc-id",
			},
			code: sderrors.CodeNotFound,
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

	// See the note in TestGetDocument_ValidationAndScope: an id that does not
	// resolve is NotFound, not BadParameters, because ids are opaque.
	_, err = svc.DeleteDocument(ctx, &apppb.DeleteDocumentRequest{
		CollectionId: collectionsmodel.NewCollectionID().String(),
		DocumentId:   "not-a-valid-doc-id",
	})
	assertErrorCode(t, err, sderrors.CodeNotFound)
}

// TestDeleteDocument_RemovesDocument is the regression bar for grandturks#936:
// the document must actually be gone, and saying so twice must not report
// success the second time.
func TestDeleteDocument_RemovesDocument(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)
	ctx := context.Background()

	coll, err := svc.CreateCollection(ctx, &apppb.CreateCollectionRequest{Description: strPtr("delete-me")})
	assert.NoError(t, err)
	cid := coll.XMetadata.CollectionId

	created, err := svc.CreateDocument(ctx, &apppb.CreateDocumentRequest{
		CollectionId: cid,
		Data:         []byte(`{"field":"value"}`),
	})
	assert.NoError(t, err)
	did := created.XMetadata.DocumentId

	_, err = svc.DeleteDocument(ctx, &apppb.DeleteDocumentRequest{CollectionId: cid, DocumentId: did})
	assert.NoError(t, err)

	// The delete is a gorm soft delete, so this also proves the read path's
	// implicit "deleted_at IS NULL" filter is doing the work we rely on.
	_, err = svc.GetDocument(ctx, &apppb.GetDocumentRequest{CollectionId: cid, DocumentId: did})
	assertErrorCode(t, err, sderrors.CodeNotFound)

	_, err = svc.DeleteDocument(ctx, &apppb.DeleteDocumentRequest{CollectionId: cid, DocumentId: did})
	assertErrorCode(t, err, sderrors.CodeNotFound)
}

// TestDeleteDocument_DisappearsFromQuery pins the collection listing, which is
// the surface most callers poll rather than fetching documents by id.
func TestDeleteDocument_DisappearsFromQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)
	ctx := context.Background()

	coll, err := svc.CreateCollection(ctx, &apppb.CreateCollectionRequest{Description: strPtr("query-after-delete")})
	assert.NoError(t, err)
	cid := coll.XMetadata.CollectionId

	keep, err := svc.CreateDocument(ctx, &apppb.CreateDocumentRequest{
		CollectionId: cid,
		Data:         []byte(`{"keep":true}`),
	})
	assert.NoError(t, err)
	drop, err := svc.CreateDocument(ctx, &apppb.CreateDocumentRequest{
		CollectionId: cid,
		Data:         []byte(`{"keep":false}`),
	})
	assert.NoError(t, err)

	_, err = svc.DeleteDocument(ctx, &apppb.DeleteDocumentRequest{
		CollectionId: cid,
		DocumentId:   drop.XMetadata.DocumentId,
	})
	assert.NoError(t, err)

	resp, err := svc.QueryDocument(ctx, &apppb.QueryDocumentRequest{CollectionId: cid})
	assert.NoError(t, err)

	var got []string
	for _, doc := range resp.Docs {
		got = append(got, doc.XMetadata.DocumentId)
	}
	assert.Equal(t, []string{keep.XMetadata.DocumentId}, got)
}

// TestDeleteDocument_WrongCollectionReturnsNotFound proves the delete is scoped
// rather than keyed on the document id alone. Before the scope check this
// returned success while deleting nothing.
func TestDeleteDocument_WrongCollectionReturnsNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)
	ctx := context.Background()

	coll, err := svc.CreateCollection(ctx, &apppb.CreateCollectionRequest{Description: strPtr("owner")})
	assert.NoError(t, err)
	cid := coll.XMetadata.CollectionId

	other, err := svc.CreateCollection(ctx, &apppb.CreateCollectionRequest{Description: strPtr("bystander")})
	assert.NoError(t, err)

	created, err := svc.CreateDocument(ctx, &apppb.CreateDocumentRequest{
		CollectionId: cid,
		Data:         []byte(`{"field":"value"}`),
	})
	assert.NoError(t, err)
	did := created.XMetadata.DocumentId

	_, err = svc.DeleteDocument(ctx, &apppb.DeleteDocumentRequest{
		CollectionId: other.XMetadata.CollectionId,
		DocumentId:   did,
	})
	assertErrorCode(t, err, sderrors.CodeNotFound)

	// ... and the document is untouched in its own collection.
	got, err := svc.GetDocument(ctx, &apppb.GetDocumentRequest{CollectionId: cid, DocumentId: did})
	assert.NoError(t, err)
	assert.Equal(t, did, got.XMetadata.DocumentId)
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

// ListCollections returned 501 while the OpenAPI spec advertised it
// (FootprintAI/restcol#144). It is the only listed call that needs no prior
// knowledge - you cannot GET /collections/{id} without an id - so it is the
// first thing anyone tries from the API doc.

func TestListCollections_EmptyProjectIsAnEmptyListNotAnError(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)

	// A new project has no collections. That is a normal state, not a failure:
	// returning NotFound would make every client special-case its own first run.
	resp, err := svc.ListCollections(context.Background(), &apppb.ListCollectionsRequest{})
	assert.NoError(t, err)
	if assert.NotNil(t, resp) {
		assert.Empty(t, resp.Collections)
	}
}

func TestListCollections_ReturnsCreatedCollectionsNewestFirst(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)
	ctx := context.Background()

	var ids []string
	for _, name := range []string{"first", "second", "third"} {
		created, err := svc.CreateCollection(ctx, &apppb.CreateCollectionRequest{
			Description: strPtr(name),
		})
		assert.NoError(t, err)
		ids = append(ids, created.XMetadata.CollectionId)
	}

	resp, err := svc.ListCollections(ctx, &apppb.ListCollectionsRequest{})
	assert.NoError(t, err)
	assert.Len(t, resp.Collections, 3)

	got := make([]string, 0, len(resp.Collections))
	for _, c := range resp.Collections {
		got = append(got, c.XMetadata.CollectionId)
	}
	// Newest first, which is what the storage query orders by and what the
	// proto comment promises. Asserting the order rather than set membership,
	// because "newest first" is part of the contract.
	assert.Equal(t, []string{ids[2], ids[1], ids[0]}, got)
}

func TestListCollections_ItemMatchesTheSingleCollectionGet(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)
	ctx := context.Background()

	created, err := svc.CreateCollection(ctx, &apppb.CreateCollectionRequest{
		Description: strPtr("same-shape"),
	})
	assert.NoError(t, err)
	cid := created.XMetadata.CollectionId

	list, err := svc.ListCollections(ctx, &apppb.ListCollectionsRequest{})
	assert.NoError(t, err)
	if !assert.Len(t, list.Collections, 1) {
		return
	}
	one, err := svc.GetCollection(ctx, &apppb.GetCollectionRequest{CollectionId: cid})
	assert.NoError(t, err)

	// One resource, one shape. Both endpoints render through newPbCollection,
	// and this is what stops them drifting apart later.
	item := list.Collections[0]
	assert.Equal(t, one.XMetadata.CollectionId, item.XMetadata.CollectionId)
	assert.Equal(t, one.Description, item.Description)
	assert.Equal(t, one.CollectionType, item.CollectionType)
	assert.Equal(t, len(one.Schemas), len(item.Schemas))
}

func TestListCollections_IgnoresProjectIdInTheRequest(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, proj := newTestService(t)
	ctx := context.Background()

	_, err := svc.CreateCollection(ctx, &apppb.CreateCollectionRequest{
		Description: strPtr("mine"),
	})
	assert.NoError(t, err)

	// THE TENANT CHECK. The path carries a projectId because the URL shape
	// needs one, but the project must come from the caller's credential. If
	// req.ProjectId were honoured, editing the path would list somebody else's
	// collections - and the naive implementation is exactly the one that does.
	resp, err := svc.ListCollections(ctx, &apppb.ListCollectionsRequest{
		ProjectId: projectsmodel.NewProjectID(4242).String(),
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Collections, 1,
		"a foreign projectId in the request must not change whose collections are returned")
	assert.NotEqual(t, "4242", proj.ID.String())
}

func TestListCollections_DeletedCollectionDisappears(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a local postgres; skipping under -short")
	}
	svc, _ := newTestService(t)
	ctx := context.Background()

	created, err := svc.CreateCollection(ctx, &apppb.CreateCollectionRequest{
		Description: strPtr("doomed"),
	})
	assert.NoError(t, err)

	_, err = svc.DeleteCollection(ctx, &apppb.DeleteCollectionRequest{
		CollectionId: created.XMetadata.CollectionId,
	})
	assert.NoError(t, err)

	// Soft delete. gorm excludes it, but only if the model carries DeletedAt -
	// worth pinning, because a list that keeps returning deleted rows looks
	// like a caching bug from the client side.
	resp, err := svc.ListCollections(ctx, &apppb.ListCollectionsRequest{})
	assert.NoError(t, err)
	assert.Empty(t, resp.Collections)
}
