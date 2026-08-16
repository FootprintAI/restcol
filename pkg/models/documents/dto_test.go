package modeldocuments

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	modelcollections "github.com/footprintai/restcol/pkg/models/collections"
	modelprojects "github.com/footprintai/restcol/pkg/models/projects"
)

// The metadata a caller sees has to say when a document was last written.
//
// Writing an existing documentId is an UPSERT - there is no Update rpc - so
// without _updatedAt an overwritten document is indistinguishable from one
// written once, in content and in metadata alike. A caller cannot tell that its
// data was replaced, or when.
//
// The column was never missing: ModelDocument has always carried updated_at and
// the upsert in DocumentCURD.Write lists it in DoUpdates. Only the API hid it,
// which is why this is a mapping fix rather than a migration.

func fixture(created, updated time.Time) *ModelDocument {
	return &ModelDocument{
		ID:                NewDocumentID(),
		CreatedAt:         created,
		UpdatedAt:         updated,
		Data:              NewModelDocumentData(map[string]interface{}{"k": "v"}),
		ModelCollectionID: modelcollections.CollectionID("col-1"),
		ModelProjectID:    modelprojects.ProjectID("proj-1"),
	}
}

func TestUpdatedAtIsSurfaced(t *testing.T) {
	created := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	updated := time.Date(2026, 8, 9, 17, 30, 0, 0, time.UTC)

	md := NewPbDocumentMetadata(fixture(created, updated))

	require.NotNil(t, md.XUpdatedAt, "_updatedAt must be populated; a nil here is the gap this fixes")
	assert.Equal(t, updated.UTC(), md.XUpdatedAt.AsTime().UTC())
	assert.Equal(t, created.UTC(), md.XCreatedAt.AsTime().UTC(),
		"_createdAt must not be overwritten by the new field")
}

// The distinguishing case, and the whole point: two documents that differ ONLY
// in having been rewritten must be distinguishable through the API.
func TestARewrittenDocumentIsDistinguishableFromAFreshOne(t *testing.T) {
	created := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)

	fresh := NewPbDocumentMetadata(fixture(created, created))
	rewritten := NewPbDocumentMetadata(fixture(created, created.Add(48*time.Hour)))

	assert.Equal(t, fresh.XCreatedAt.AsTime(), rewritten.XCreatedAt.AsTime(),
		"createdAt is identical - it is the only thing these two shared before")
	assert.True(t,
		rewritten.XUpdatedAt.AsTime().After(fresh.XUpdatedAt.AsTime()),
		"the rewritten document must report a later _updatedAt; if these compare "+
			"equal the metadata still cannot tell an overwrite from a first write")
}

// On a document nobody has rewritten, updatedAt equals createdAt - gorm sets
// both on insert. Pinned because it justifies the field being non-optional:
// there is no such thing as a document without an updatedAt.
func TestNeverRewrittenReportsUpdatedEqualToCreated(t *testing.T) {
	at := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)

	md := NewPbDocumentMetadata(fixture(at, at))

	require.NotNil(t, md.XUpdatedAt)
	assert.Equal(t, md.XCreatedAt.AsTime(), md.XUpdatedAt.AsTime())
}

// A soft-deleted document still reports when it was last written. The three
// timestamps are independent and must not overwrite one another - an earlier
// version of this mapping set XDeletedAt to nil unconditionally and then
// populated it, which is exactly the shape that loses a field silently.
func TestDeletedAtAndUpdatedAtCoexist(t *testing.T) {
	created := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	updated := created.Add(time.Hour)
	deleted := created.Add(2 * time.Hour)

	md := fixture(created, updated)
	md.DeletedAt = gorm.DeletedAt{Time: deleted, Valid: true}

	pb := NewPbDocumentMetadata(md)

	require.NotNil(t, pb.XDeletedAt, "deletedAt must still be populated")
	require.NotNil(t, pb.XUpdatedAt, "updatedAt must survive alongside deletedAt")
	assert.Equal(t, deleted.UTC(), pb.XDeletedAt.AsTime().UTC())
	assert.Equal(t, updated.UTC(), pb.XUpdatedAt.AsTime().UTC())
	assert.Equal(t, created.UTC(), pb.XCreatedAt.AsTime().UTC())
}

// Attribution must survive the DTO. The columns are useless if the API does not
// carry them, and the failure would be silent: the database would hold the
// right answer while every caller saw nothing.
func TestWriterAttributionIsSurfaced(t *testing.T) {
	const alice = "grandturk:apikey:42:alice-key"
	const bob = "grandturk:apikey:42:bob-key"

	md := fixture(time.Now(), time.Now())
	md.CreatedBy = alice
	md.UpdatedBy = bob

	pb := NewPbDocumentMetadata(md)

	assert.Equal(t, alice, pb.XCreatedBy)
	assert.Equal(t, bob, pb.XUpdatedBy)
}

// The distinguishing case, through the API: a document overwritten by someone
// else must be tellable from one written once by its creator.
func TestOverwrittenByAnotherPrincipalIsVisibleThroughTheApi(t *testing.T) {
	const alice = "grandturk:apikey:42:alice-key"
	const bob = "grandturk:apikey:42:bob-key"
	at := time.Now()

	own := fixture(at, at)
	own.CreatedBy, own.UpdatedBy = alice, alice

	replaced := fixture(at, at.Add(time.Hour))
	replaced.CreatedBy, replaced.UpdatedBy = alice, bob

	ownPb := NewPbDocumentMetadata(own)
	replacedPb := NewPbDocumentMetadata(replaced)

	assert.Equal(t, ownPb.XCreatedBy, replacedPb.XCreatedBy,
		"both were created by the same principal")
	assert.NotEqual(t, ownPb.XUpdatedBy, replacedPb.XUpdatedBy,
		"but only one was last written by somebody else - if these compared "+
			"equal the API could not tell an overwrite from a first write")
}

// Empty means unattributed, and must not become a placeholder on the way out.
func TestUnattributedStaysEmptyThroughTheApi(t *testing.T) {
	pb := NewPbDocumentMetadata(fixture(time.Now(), time.Now()))

	assert.Empty(t, pb.XCreatedBy,
		"an unattributed document must not acquire a writer in the DTO")
	assert.Empty(t, pb.XUpdatedBy)
}
