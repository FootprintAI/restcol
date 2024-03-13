package app

import (
	"context"
	"errors"

	sderrors "github.com/sdinsure/agent/pkg/errors"
	"github.com/sdinsure/agent/pkg/logger"

	apppb "github.com/footprintai/restcol/api/pb/proto"
	collectionsmodel "github.com/footprintai/restcol/pkg/models/collections"
	projectsmodel "github.com/footprintai/restcol/pkg/models/projects"
	"github.com/footprintai/restcol/pkg/nullable"
	collectionsstorage "github.com/footprintai/restcol/pkg/storage/collections"
)

type RestColServiceServerService struct {
	apppb.UnimplementedRestColServiceServer

	log            logger.Logger
	collectionCURD *collectionsstorage.CollectionCURD
}

func NewRestColServiceServerService(log logger.Logger, collectionCURD *collectionsstorage.CollectionCURD) *RestColServiceServerService {
	return &RestColServiceServerService{
		log:            log,
		collectionCURD: collectionCURD,
	}
}

func (r *RestColServiceServerService) CreateCollection(ctx context.Context, req *apppb.CreateCollectionRequest) (*apppb.CreateCollectionResponse, error) {
	cid, err := getOrCreateCid(req.Cid)
	if err != nil {
		return nil, err
	}
	collectionType := apppb.CollectionType_COLLECTION_TYPE_REGULAR_FILES
	if req.CollectionType != nil {
		collectionType = *req.CollectionType
	}
	var summary string
	if req.Description != nil {
		summary = *req.Description
	}
	var modelSchemaSlice []collectionsmodel.ModelSchema
	if len(req.Schemas) > 0 {
		// request is with a specific schema, use it
		modelSchema, _ := collectionsmodel.NewModelSchema(req.Schemas)
		modelSchemaSlice = append(modelSchemaSlice, modelSchema)
	}

	mc := collectionsmodel.NewModelCollection(
		getPid(req.Pid),
		cid,
		collectionType,
		summary,
		modelSchemaSlice,
	)
	if err := r.collectionCURD.Write(ctx, "", mc); err != nil {
		return nil, err
	}

	resp := &apppb.CreateCollectionResponse{
		XMetadata:      collectionsmodel.NewPbCollectionMetadata(mc),
		Description:    mc.Summary,
		CollectionType: mc.Type.Proto(),
		Schemas:        collectionsmodel.NewPbSchemaFields(mc.Schemas[0]),
	}

	return resp, nil
}

func getOrCreateCid(cidStrP *string) (collectionsmodel.CollectionID, error) {
	var cid collectionsmodel.CollectionID
	var err error
	if cidStrP != nil {
		cid, err = collectionsmodel.NewCollectionIDFromStr(*cidStrP)
		if err != nil {
			return collectionsmodel.CollectionID{}, err
		}
		return cid, nil
	}
	return collectionsmodel.NewCollectionID(), nil
}

func getPid(pid string) projectsmodel.ProjectID {
	return projectsmodel.NewProjectIDStr(pid)
}

func (r *RestColServiceServerService) ListCollections(ctx context.Context, req *apppb.ListCollectionsRequest) (*apppb.ListCollectionsResponse, error) {
	return nil, sderrors.NewNotImplError(errors.New("not implemented"))
}

func (r *RestColServiceServerService) GetCollection(ctx context.Context, req *apppb.GetCollectionRequest) (*apppb.GetCollectionResponse, error) {
	// pid := getPid(req.Pid)
	cid, err := getOrCreateCid(nullable.StringP(req.Cid))
	if err != nil {
		return nil, err
	}
	mc, err := r.collectionCURD.GetLatestSchema(ctx, "", cid)
	if err != nil {
		return nil, err
	}
	resp := &apppb.GetCollectionResponse{
		XMetadata:      collectionsmodel.NewPbCollectionMetadata(mc),
		Description:    mc.Summary,
		CollectionType: mc.Type.Proto(),
		Schemas:        collectionsmodel.NewPbSchemaFields(mc.Schemas[0]),
	}
	return resp, nil
}
func (r *RestColServiceServerService) DeleteCollection(ctx context.Context, req *apppb.DeleteCollectionRequest) (*apppb.DeleteCollectionResponse, error) {
	return nil, sderrors.NewNotImplError(errors.New("not implemented"))
}
func (r *RestColServiceServerService) CreateDocument(ctx context.Context, req *apppb.CreateDocumentRequest) (*apppb.CreateDocumentResponse, error) {
	return nil, sderrors.NewNotImplError(errors.New("not implemented"))
}
func (r *RestColServiceServerService) GetDocument(ctx context.Context, req *apppb.GetDocumentRequest) (*apppb.GetDocumentResponse, error) {
	return nil, sderrors.NewNotImplError(errors.New("not implemented"))
}
func (r *RestColServiceServerService) DeleteDocument(ctx context.Context, req *apppb.DeleteDocumentRequest) (*apppb.DeleteDocumentResponse, error) {
	return nil, sderrors.NewNotImplError(errors.New("not implemented"))
}
