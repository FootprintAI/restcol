package modeldocuments

import (
	"time"

	apppb "github.com/footprintai/restcol/api/pb/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewPbDocumentMetadata(md *ModelDocument) *apppb.DataMetadata {
	metadata := &apppb.DataMetadata{
		Pid: md.ModelProjectID.String(),
		Cid: md.ModelCollectionID.String(),
		Did: md.ID.String(),
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
