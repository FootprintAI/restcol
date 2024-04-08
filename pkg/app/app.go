package app

import (
	"context"
	"errors"

	sderrors "github.com/sdinsure/agent/pkg/errors"
	"github.com/sdinsure/agent/pkg/logger"
	sdinsureruntime "github.com/sdinsure/agent/pkg/runtime"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"gorm.io/datatypes"

	apppb "github.com/footprintai/restcol/api/pb/proto"
	collectionsmodel "github.com/footprintai/restcol/pkg/models/collections"
	documentsmodel "github.com/footprintai/restcol/pkg/models/documents"
	projectsmodel "github.com/footprintai/restcol/pkg/models/projects"
	schemafinder "github.com/footprintai/restcol/pkg/schema"
	collectionsstorage "github.com/footprintai/restcol/pkg/storage/collections"
	documentsstorage "github.com/footprintai/restcol/pkg/storage/documents"
	collectionsswagger "github.com/footprintai/restcol/pkg/swagger/collections"
)

func NewRestColServiceServerService(
	log logger.Logger,
	collectionCURD *collectionsstorage.CollectionCURD,
	documentCURD *documentsstorage.DocumentCURD,
) *RestColServiceServerService {
	return &RestColServiceServerService{
		log:            log,
		collectionCURD: collectionCURD,
		documentCURD:   documentCURD,
	}
}

type RestColServiceServerService struct {
	apppb.UnimplementedRestColServiceServer

	log            logger.Logger
	collectionCURD *collectionsstorage.CollectionCURD
	documentCURD   *documentsstorage.DocumentCURD

	//optional
	defaultProjectResolver sdinsureruntime.ProjectResolver
}

func (r *RestColServiceServerService) SetDefaultProjectResolver(projectResolver sdinsureruntime.ProjectResolver) {
	r.defaultProjectResolver = projectResolver
}

func (r *RestColServiceServerService) GetSwaggerDoc(ctx context.Context, req *apppb.GetSwaggerDocRequest) (*httpbody.HttpBody, error) {
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	collectionList, err := r.collectionCURD.ListByProjectID(ctx, "", projectId)
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
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	var cid collectionsmodel.CollectionID
	if req.CollectionId == nil {
		cid = collectionsmodel.NewCollectionID()
	} else {
		cid, err = collectionsmodel.Parse(*req.CollectionId)
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
		projectId,
		cid,
		collectionType,
		summary,
		modelSchemaSlice,
	)
	if err := r.collectionCURD.Write(ctx, "", &mc); err != nil {
		return nil, err
	}

	resp := &apppb.CreateCollectionResponse{
		XMetadata:      collectionsmodel.NewPbCollectionMetadata(&mc),
		Description:    mc.Summary,
		CollectionType: mc.Type.Proto(),
		Schemas:        collectionsmodel.NewPbSchemaFields(mc.Schemas[0]),
	}

	return resp, nil
}

func (r *RestColServiceServerService) getProjectIdFromCtx(ctx context.Context) (pid projectsmodel.ProjectID, reterr error) {
	if r.defaultProjectResolver == nil {
		pid = projectsmodel.ProjectID("invalid")
		reterr = errors.New("no project resolver")
		return
	}
	projectInfor, found := r.defaultProjectResolver.ProjectInfo(ctx)
	if !found {
		r.log.Info("no valid project id found, use default: %+v\n", pid)
		pid = projectsmodel.NewProjectID(1001)
		return
	}
	rawPid, err := projectInfor.GetProjectID()
	if err != nil {
		pid = projectsmodel.NewProjectID(1001)
		return
	}
	return projectsmodel.ProjectID(rawPid), nil
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
	if len(req.CollectionId) == 0 {
		return nil, sderrors.NewBadParamsError(errors.New("missing required field"))
	}
	cid, err := collectionsmodel.Parse(req.CollectionId)
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
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	var cid collectionsmodel.CollectionID
	if req.CollectionId != "" {
		cid, err = collectionsmodel.Parse(req.CollectionId)
		if err != nil {
			return nil, err
		}
	} else {
		cid, err = r.getCollectionIDFromSchemas()
		if err != nil {
			return nil, err
		}
	}
	// auto detect schema
	schemaBuilder := schemafinder.NewSchemaBuilder()
	_, modelSchema, err := schemaBuilder.Parse(req.Data)
	if err != nil {
		r.log.Error("failed to convert into modelschema, err:%+v\n", err)
		return nil, err
	}
	docModel := &documentsmodel.ModelDocument{
		ID:                documentsmodel.NewDocumentID(),
		Data:              datatypes.JSON(req.Data),
		ModelCollectionID: cid,
		ModelCollection: collectionsmodel.NewModelCollection(
			projectId,
			cid,
			apppb.CollectionType_COLLECTION_TYPE_REGULAR_FILES,
			"auto created collection",
			[]collectionsmodel.ModelSchema{
				*modelSchema,
			},
		),
		ModelProjectID: projectId,
	}
	if err := r.documentCURD.Write(ctx, "", docModel); err != nil {
		r.log.Error("failed to write docmodel, err:%+v\n", err)
		return nil, err
	}
	return &apppb.CreateDocumentResponse{
		XMetadata: documentsmodel.NewPbDocumentMetadata(docModel),
	}, nil
}
func (r *RestColServiceServerService) GetDocument(ctx context.Context, req *apppb.GetDocumentRequest) (*apppb.GetDocumentResponse, error) {
	// TODO: use pid and cid for permission checking
	// as for retrieving data, did is only required field
	did, err := documentsmodel.Parse(req.DocumentId)
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
