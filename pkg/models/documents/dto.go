package modeldocuments

import (
	"time"

	apppb "github.com/footprintai/restcol/api/pb/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewPbDocumentMetadata(md *ModelDocument) *apppb.DataMetadata {
	metadata := &apppb.DataMetadata{
		ProjectId:    md.ModelProjectID.String(),
		CollectionId: md.ModelCollectionID.String(),
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
