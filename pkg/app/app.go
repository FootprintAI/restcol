package app

import (
	"context"
	"errors"

	sderrors "github.com/sdinsure/agent/pkg/errors"
	"github.com/sdinsure/agent/pkg/logger"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"gorm.io/datatypes"

	apppb "github.com/footprintai/restcol/api/pb/proto"
	collectionsmodel "github.com/footprintai/restcol/pkg/models/collections"
	documentsmodel "github.com/footprintai/restcol/pkg/models/documents"
	projectsmodel "github.com/footprintai/restcol/pkg/models/projects"
	collectionsstorage "github.com/footprintai/restcol/pkg/storage/collections"
	documentsstorage "github.com/footprintai/restcol/pkg/storage/documents"
	collectionsswagger "github.com/footprintai/restcol/pkg/swagger/collections"
)

func NewRestColServiceServerService(
	log logger.Logger,
	collectionCURD *collectionsstorage.CollectionCURD,
	documentCURD *documentsstorage.DocumentCURD,
	projectGetter ProjectGetter,
) *RestColServiceServerService {
	return &RestColServiceServerService{
		log:            log,
		collectionCURD: collectionCURD,
		documentCURD:   documentCURD,
		projectGetter:  projectGetter,
	}
}

type ProjectGetter interface {
	GetProject(ctx context.Context, pid projectsmodel.ProjectID) (*projectsmodel.ModelProject, error)
}

//type SchemaGetter interface {
//	GetCollection(ctx context.Context) (*collectionsmodel.ModelCollection, error)
//}

type RestColServiceServerService struct {
	apppb.UnimplementedRestColServiceServer

	log            logger.Logger
	collectionCURD *collectionsstorage.CollectionCURD
	documentCURD   *documentsstorage.DocumentCURD
	projectGetter  ProjectGetter
	//collectionGetter CollectionGetter
}

func (r *RestColServiceServerService) GetSwaggerDoc(ctx context.Context, req *apppb.GetSwaggerDocRequest) (*httpbody.HttpBody, error) {
	modelProject, err := r.getModelProject(ctx, req.Pid)
	collectionList, err := r.collectionCURD.ListByProjectID(ctx, "", modelProject.ID)
	if err != nil {
		return nil, err
	}
	colSwaggerDoc := collectionsswagger.NewCollectionSwaggerDoc(collectionList...)
	colSwaggerDocInStr, err := colSwaggerDoc.RenderDoc()
	if err != nil {
		return nil, err
	}
	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        []byte(colSwaggerDocInStr),
	}, nil
}

func (r *RestColServiceServerService) CreateCollection(ctx context.Context, req *apppb.CreateCollectionRequest) (*apppb.CreateCollectionResponse, error) {
	modelProject, err := r.getModelProject(ctx, req.Pid)
	if err != nil {
		return nil, err
	}
	var cid collectionsmodel.CollectionID
	if req.Cid == nil {
		cid = collectionsmodel.NewCollectionID()
	} else {
		cid, err = collectionsmodel.Parse(*req.Cid)
		if err != nil {
			return nil, err
		}
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
		modelProject,
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

func (r *RestColServiceServerService) getModelProject(ctx context.Context, pidStr string) (*projectsmodel.ModelProject, error) {
	var pid projectsmodel.ProjectID
	if pidStr == "" {
		// use default pid
		pid = projectsmodel.NewProjectID(1001)
		r.log.Info("no valid project id found, use default: %+v\n", pid)
	} else {
		pid = projectsmodel.NewProjectIDStr(pidStr)
	}
	return r.projectGetter.GetProject(ctx, pid)
}

// TODO getCollectionIDFromSchemas would lookup collection id with schema list given// This should scan all collections and match by its schema and return the right collection id
// For now, we do nothing but return a new one
func (r *RestColServiceServerService) getCollectionIDFromSchemas() (collectionsmodel.CollectionID, error) {
	return collectionsmodel.NewCollectionID(), nil
}

func (r *RestColServiceServerService) ListCollections(ctx context.Context, req *apppb.ListCollectionsRequest) (*apppb.ListCollectionsResponse, error) {
	return nil, sderrors.NewNotImplError(errors.New("not implemented"))
}

func (r *RestColServiceServerService) GetCollection(ctx context.Context, req *apppb.GetCollectionRequest) (*apppb.GetCollectionResponse, error) {
	var cid collectionsmodel.CollectionID
	if len(req.Cid) == 0 {
		return nil, sderrors.NewBadParamsError(errors.New("missing required field"))
	}
	cid, err := collectionsmodel.Parse(req.Cid)
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
	modelProject, err := r.getModelProject(ctx, req.Pid)
	if err != nil {
		return nil, err
	}
	var cid collectionsmodel.CollectionID
	if req.Cid != "" {
		cid, err = collectionsmodel.Parse(req.Cid)
		if err != nil {
			return nil, err
		}
	} else {
		cid, err = r.getCollectionIDFromSchemas()
		if err != nil {
			return nil, err
		}
	}
	docModel := &documentsmodel.ModelDocument{
		ID:                documentsmodel.NewDocumentID(),
		Data:              datatypes.JSON(req.Data),
		ModelCollectionID: cid,
		ModelProjectID:    modelProject.ID,
	}
	if err := r.documentCURD.Write(ctx, "", docModel); err != nil {
		return nil, err
	}
	return &apppb.CreateDocumentResponse{
		XMetadata: documentsmodel.NewPbDocumentMetadata(docModel),
	}, nil
}
func (r *RestColServiceServerService) GetDocument(ctx context.Context, req *apppb.GetDocumentRequest) (*apppb.GetDocumentResponse, error) {
	// TODO: use pid and cid for permission checking
	// as for retrieving data, did is only required field
	did, err := documentsmodel.Parse(req.Did)
	if err != nil {
		return nil, err
	}
	docModel, err := r.documentCURD.Get(ctx, "", did)
	if err != nil {
		return nil, err
	}

	return &apppb.GetDocumentResponse{
		XMetadata: documentsmodel.NewPbDocumentMetadata(docModel),
		Data:      []byte(docModel.Data),
	}, nil

}
func (r *RestColServiceServerService) DeleteDocument(ctx context.Context, req *apppb.DeleteDocumentRequest) (*apppb.DeleteDocumentResponse, error) {
	return nil, sderrors.NewNotImplError(errors.New("not implemented"))
}
