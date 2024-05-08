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
		XDeletedAt: nil,
	}
	if deletedAt, _ := md.DeletedAt.Value(); deletedAt != nil {
		metadata.XDeletedAt = timestamppb.New(deletedAt.(time.Time))
	}
	return metadata

}
