package modeldocuments

import (
	"time"

	apppb "github.com/footprintai/restcol/api/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewPbDocumentMetadata(md *ModelDocument) *apppb.DataMetadata {
	var schemaId string
	if len(md.ModelCollection.Schemas) > 0 {
		schemaId = md.ModelCollection.Schemas[0].ID.String()
	}

	metadata := &apppb.DataMetadata{
		ProjectId:    md.ModelProjectID.String(),
		CollectionId: md.ModelCollectionID.String(),
		SchemaId:     schemaId,
		DocumentId:   md.ID.String(),
		//Dataformat: nil;
		// FIXME(hsiny): need to fix dataformat
		XCreatedAt: timestamppb.New(md.CreatedAt),
		// Surfaced so a caller can tell an overwritten document from one written
		// once. Writing an existing documentId is an upsert and there is no
		// Update rpc, so before this the two were indistinguishable in both
		// content and metadata. gorm has always maintained the column; only the
		// API was hiding it.
		XUpdatedAt: timestamppb.New(md.UpdatedAt),
		// Empty when no caller resolver is wired - which says "unattributed"
		// rather than naming a principal that does not exist.
		XCreatedBy: md.CreatedBy,
		XUpdatedBy: md.UpdatedBy,
		XDeletedAt: nil,
	}
	if deletedAt, _ := md.DeletedAt.Value(); deletedAt != nil {
		metadata.XDeletedAt = timestamppb.New(deletedAt.(time.Time))
	}
	return metadata

}
