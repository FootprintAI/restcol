package app

import (
	"context"
	"errors"
	"strings"
	"time"

	sderrors "github.com/sdinsure/agent/pkg/errors"
	"github.com/sdinsure/agent/pkg/logger"
	sdinsureruntime "github.com/sdinsure/agent/pkg/runtime"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	apppb "github.com/footprintai/restcol/api/pb"
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
	schemaBuilder *schemafinder.SchemaBuilder,
) *RestColServiceServerService {
	return &RestColServiceServerService{
		log:            log,
		collectionCURD: collectionCURD,
		documentCURD:   documentCURD,
		schemaBuilder:  schemaBuilder,
	}
}

type RestColServiceServerService struct {
	apppb.UnimplementedRestColServiceServer

	log            logger.Logger
	collectionCURD *collectionsstorage.CollectionCURD
	documentCURD   *documentsstorage.DocumentCURD

	schemaBuilder *schemafinder.SchemaBuilder

	//optional
	defaultProjectResolver sdinsureruntime.ProjectResolver
}

func (r *RestColServiceServerService) SetDefaultProjectResolver(projectResolver sdinsureruntime.ProjectResolver) {
	r.defaultProjectResolver = projectResolver
}

func (r *RestColServiceServerService) GetSwaggerDoc(ctx context.Context, req *apppb.GetSwaggerDocRequest) (*httpbody.HttpBody, error) {
	var err error
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var collectionList []*collectionsmodel.ModelCollection
	if len(req.CollectionId) > 0 {
		var selectedCollection *collectionsmodel.ModelCollection
		cid := collectionsmodel.NewCollectionIDFromStr(req.CollectionId)
		selectedCollection, err = r.collectionCURD.GetLatestSchema(ctx, "", projectId, cid)
		if err != nil {
			return nil, err
		}
		collectionList = append(collectionList, selectedCollection)
	} else {
		// query collections from the project
		collectionList, err = r.collectionCURD.ListByProjectID(ctx, "", projectId)
		if err != nil {
			return nil, err
		}
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
	var cid collectionsmodel.CollectionID = collectionsmodel.NewCollectionID()
	if req.CollectionId != nil {
		cid = collectionsmodel.NewCollectionIDFromStr(*req.CollectionId)
	}
	collectionType := apppb.CollectionType_COLLECTION_TYPE_REGULAR_FILES
	if req.CollectionType != nil {
		collectionType = *req.CollectionType
	}
	var summary string
	if req.Description != nil {
		summary = *req.Description
	}
	var modelSchemaSlice []*collectionsmodel.ModelSchema
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
	}
	if len(mc.Schemas) > 0 {
		resp.Schemas = collectionsmodel.NewPbSchemaFields(mc.Schemas[0])
	}

	return resp, nil
}

func (r *RestColServiceServerService) getProjectIdFromCtx(ctx context.Context) (pid projectsmodel.ProjectID, reterr error) {
	if r.defaultProjectResolver == nil {
		pid = projectsmodel.ProjectID("invalid")
		reterr = errors.New("no project resolver")
		r.log.Error("getProjectIdFromCtx: no valid project resolver\n")
		return
	}
	projectInfor, found := r.defaultProjectResolver.ProjectInfo(ctx)
	if !found {
		r.log.Info("no valid project id found, use default: %+v\n", pid)
		return
	}
	rawPid, err := projectInfor.GetProjectID()
	if err != nil {
		r.log.Error("getProjectIdFromCtx: get projectid failed, err:%+v\n", err)
		return
	}
	return projectsmodel.ProjectID(rawPid), nil
}

// TODO getCollectionIDFromSchemas would lookup collection id with schema list given
// This should scan all collections and match by its schema and return the right collection id
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
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	cid = collectionsmodel.NewCollectionIDFromStr(req.CollectionId)
	mc, err := r.collectionCURD.GetLatestSchema(ctx, "", projectId, cid)
	if err != nil {
		ismyerr, myerr := sderrors.As(err)
		if ismyerr && myerr.Code() == sderrors.CodeNotFound {
			return &apppb.GetCollectionResponse{}, nil
		}
		return nil, err
	}
	resp := &apppb.GetCollectionResponse{
		XMetadata: collectionsmodel.NewPbCollectionMetadata(mc),
	}
	if mc == nil {
		return resp, nil
	}
	resp.Description = mc.Summary
	resp.CollectionType = mc.Type.Proto()
	if mc.Schemas != nil {
		resp.Schemas = collectionsmodel.NewPbSchemaFields(mc.Schemas[0])
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
		cid = collectionsmodel.NewCollectionIDFromStr(req.CollectionId)
	} else {
		cid, err = r.getCollectionIDFromSchemas()
		if err != nil {
			return nil, err
		}
	}

	modelCollection, err := r.collectionCURD.GetLatestSchema(ctx, "", projectId, cid)
	if err != nil {
		return nil, err
	}

	var docSchema *collectionsmodel.ModelSchema

	// auto detect schema
	_, inputDataSchema, valueHolder, err := r.schemaBuilder.Parse(req.Data)
	if err != nil {
		r.log.Error("failed to convert into modelschema, err:%+v\n", err)
		return nil, err
	}
	// check whether we need to create a new schema mapping
	// or use existing schema
	if len(modelCollection.Schemas) == 0 {
		// no previous schema, use the latest data
		docSchema = inputDataSchema
	} else {
		// has schema under the collection
		if r.schemaBuilder.Equals(modelCollection.Schemas[0], inputDataSchema) {
			docSchema = modelCollection.Schemas[0]
		} else {
			docSchema = inputDataSchema
		}
	}
	var docId documentsmodel.DocumentID
	if req.DocumentId == nil {
		docId = documentsmodel.NewDocumentID()
	} else {
		docId, err = documentsmodel.Parse(*req.DocumentId)
		if err != nil {
			return nil, err
		}
	}

	docModel := &documentsmodel.ModelDocument{
		ID:                docId,
		Data:              documentsmodel.NewModelDocumentData(valueHolder),
		ModelCollectionID: cid,
		ModelCollection: collectionsmodel.NewModelCollection(
			projectId,
			cid,
			apppb.CollectionType_COLLECTION_TYPE_REGULAR_FILES,
			"auto created collection",
			[]*collectionsmodel.ModelSchema{
				docSchema,
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
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	cid := collectionsmodel.NewCollectionIDFromStr(req.CollectionId)

	did, err := documentsmodel.Parse(req.DocumentId)
	if err != nil {
		return nil, err
	}
	docModel, err := r.documentCURD.Get(ctx, "", projectId, cid, did)
	if err != nil {
		return nil, err
	}
	filteredDoc, err := r.filterDocWithSelectedFields(docModel, req.FieldSelectors)
	if err != nil {
		return nil, err
	}

	return &apppb.GetDocumentResponse{
		XMetadata: documentsmodel.NewPbDocumentMetadata(docModel),
		Data:      filteredDoc,
	}, nil

}

func (r *RestColServiceServerService) filterDocWithSelectedFields(doc *documentsmodel.ModelDocument, selectedFields []string) (*structpb.Value, error) {
	r.log.Info("query doc with fields:%+v\n", selectedFields)

	if doc.Data == nil {
		return nil, nil
	}

	if len(selectedFields) == 0 {
		// no selectedFields, return all

		return structpb.NewValue(doc.Data.MapValue)
	}

	// get associated schema
	modelSchema, err := r.schemaBuilder.Flatten(doc.Data.MapValue)
	if err != nil {
		return nil, err
	}

	// make selectedFields into a lookup map
	lookupMap := make(map[string]struct{})
	for _, selectedField := range selectedFields {
		r.log.Info("query doc, add field:%s\n", strings.ToLower(selectedField))
		lookupMap[strings.ToLower(selectedField)] = struct{}{}
	}

	var fieldsInSelected []*collectionsmodel.ModelFieldSchema
	for _, dataField := range modelSchema.Fields {
		r.log.Info("query doc, select field:%s\n", dataField.FieldName.String())
		_, exist := lookupMap[dataField.FieldName.String()]
		if exist {
			r.log.Info("query doc, added field:%s\n", dataField.FieldName.String())
			fieldsInSelected = append(fieldsInSelected, dataField)
		}
	}
	// construct the whole struct with fieldsInSelected
	structWithSelectedFields, err := schemafinder.Build(fieldsInSelected)
	if err != nil {
		return nil, err
	}
	return structpb.NewValue(structWithSelectedFields)
}

func (r *RestColServiceServerService) DeleteDocument(ctx context.Context, req *apppb.DeleteDocumentRequest) (*apppb.DeleteDocumentResponse, error) {
	return nil, sderrors.NewNotImplError(errors.New("not implemented"))
}

func (r *RestColServiceServerService) QueryDocumentsStream(req *apppb.QueryDocumentStreamRequest, stream apppb.RestColService_QueryDocumentsStreamServer) error {
	ctx := stream.Context()
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return err
	}
	cid := collectionsmodel.NewCollectionIDFromStr(req.CollectionId)
	startedAt := req.SinceTs
	endedAt := req.EndedAt

	needsFollowUp := false
	if req.FollowUpMode != nil && *req.FollowUpMode {
		needsFollowUp = true
	}

	if !needsFollowUp {
		queryDocs, err := r.documentCURD.Query(
			ctx,
			"",
			projectId,
			cid,
			makeQueryConditioner(startedAt, endedAt, req.LimitCount)...,
		)
		if err != nil {
			return err
		}
		// TODO: apply field selector
		for _, doc := range queryDocs {
			filteredDoc, err := r.filterDocWithSelectedFields(doc, req.FieldSelectors)
			if err != nil {
				return err
			}

			if err := stream.Send(&apppb.GetDocumentResponse{
				XMetadata: documentsmodel.NewPbDocumentMetadata(doc),
				Data:      filteredDoc,
			}); err != nil {
				return err
			}
		}
		return nil
	}
	// enter followup mode, keep looking data
	sendingCount := 0
	for {
		r.log.Info("query with time range[%+v -> %+v], cid:%+v\n", startedAt.AsTime(), endedAt.AsTime(), cid)
		queryDocs, err := r.documentCURD.Query(
			ctx,
			"",
			projectId,
			cid,
			makeQueryConditioner(startedAt, endedAt, req.LimitCount)...,
		)
		if err != nil {
			return err
		}
		// TODO: apply field selector
		for _, doc := range queryDocs {
			filteredDoc, err := r.filterDocWithSelectedFields(doc, req.FieldSelectors)
			if err != nil {
				return err
			}

			sendingCount = sendingCount + 1
			if err := stream.Send(&apppb.GetDocumentResponse{
				XMetadata: documentsmodel.NewPbDocumentMetadata(doc),
				Data:      filteredDoc,
			}); err != nil {
				return err
			}
		}
		// update startedAt to be the last record of previous query
		if len(queryDocs) > 0 {
			startedAt = timestamppb.New(queryDocs[len(queryDocs)-1].CreatedAt)
		}

		// take a nap and would query again
		select {
		case <-ctx.Done():
			r.log.Info("query is done as ctx is done\n")
			return nil
		case <-time.After(1 * time.Second):
		}
	}
	return nil
}

func makeQueryConditioner(startedAt *timestamppb.Timestamp, endedAt *timestamppb.Timestamp, limitCount *int32) []documentsstorage.QueryConditioner {
	var cnds []documentsstorage.QueryConditioner
	if startedAt != nil {
		cnds = append(cnds, documentsstorage.WithStartedAt(startedAt.AsTime()))
	}
	if endedAt != nil {
		cnds = append(cnds, documentsstorage.WithEndedAt(endedAt.AsTime()))
	}
	if limitCount != nil {
		cnds = append(cnds, documentsstorage.WithLimitCount(*limitCount))
	}
	return cnds
}

func (r *RestColServiceServerService) QueryDocument(ctx context.Context, req *apppb.QueryDocumentRequest) (*apppb.QueryDocumentResponse, error) {
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	cid := collectionsmodel.NewCollectionIDFromStr(req.CollectionId)
	queryDocs, err := r.documentCURD.Query(
		ctx,
		"",
		projectId,
		cid,
		makeQueryConditioner(req.SinceTs, req.EndedAt, req.LimitCount)...,
	)
	if err != nil {
		return nil, err
	}
	resp := &apppb.QueryDocumentResponse{}
	for _, doc := range queryDocs {
		filteredDocBytes, err := r.filterDocWithSelectedFields(doc, req.FieldSelectors)
		if err != nil {
			return nil, err
		}
		resp.Docs = append(resp.Docs, &apppb.GetDocumentResponse{
			XMetadata: documentsmodel.NewPbDocumentMetadata(doc),
			Data:      filteredDocBytes,
		})
	}
	return resp, nil

}
