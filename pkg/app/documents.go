package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	sderrors "github.com/sdinsure/agent/pkg/errors"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	apppb "github.com/footprintai/restcol/api/pb"
	collectionsmodel "github.com/footprintai/restcol/pkg/models/collections"
	documentsmodel "github.com/footprintai/restcol/pkg/models/documents"
	schemafinder "github.com/footprintai/restcol/pkg/schema"
	documentsstorage "github.com/footprintai/restcol/pkg/storage/documents"
)

// CreateDocument writes a document, auto-detecting its schema against the
// collection's latest known schema. When no collection is supplied, one is
// provisioned on the fly via getCollectionIDFromSchemas.
func (r *RestColServiceServerService) CreateDocument(ctx context.Context, req *apppb.CreateDocumentRequest) (*apppb.CreateDocumentResponse, error) {
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	var cid collectionsmodel.CollectionID
	if req.CollectionId != "" {
		cid = collectionsmodel.NewCollectionIDFromStr(req.CollectionId)
	} else {
		cid, err = r.getCollectionIDFromSchemas(ctx, projectId)
		if err != nil {
			return nil, err
		}
	}

	modelCollection, err := r.collectionCURD.GetLatestSchema(ctx, "", projectId, cid)
	if err != nil {
		return nil, err
	}

	_, inputDataSchema, valueHolder, err := r.schemaBuilder.Parse(req.Data)
	if err != nil {
		r.log.Error("failed to convert into modelschema, err:%+v\n", err)
		return nil, err
	}

	// Reuse the existing schema when it matches; otherwise attach the newly
	// parsed schema so schema evolution is tracked per-document.
	var docSchema *collectionsmodel.ModelSchema
	if len(modelCollection.Schemas) == 0 {
		docSchema = inputDataSchema
	} else if r.schemaBuilder.Equals(modelCollection.Schemas[0], inputDataSchema) {
		docSchema = modelCollection.Schemas[0]
	} else {
		docSchema = inputDataSchema
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
			[]*collectionsmodel.ModelSchema{docSchema},
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

// GetDocument returns a single document, optionally filtered to selected
// fields. A mismatch between the request's (project, collection) scope and the
// stored document surfaces as NotFound — callers never leak cross-tenant data.
func (r *RestColServiceServerService) GetDocument(ctx context.Context, req *apppb.GetDocumentRequest) (*apppb.GetDocumentResponse, error) {
	if req.CollectionId == "" {
		return nil, sderrors.NewBadParamsError(errors.New("collection_id is required"))
	}
	if req.DocumentId == "" {
		return nil, sderrors.NewBadParamsError(errors.New("document_id is required"))
	}
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	cid := collectionsmodel.NewCollectionIDFromStr(req.CollectionId)
	did, err := documentsmodel.Parse(req.DocumentId)
	if err != nil {
		return nil, sderrors.NewBadParamsError(fmt.Errorf("invalid document_id: %w", err))
	}
	docModel, err := r.documentCURD.Get(ctx, "", projectId, cid, did)
	if err != nil {
		return nil, err
	}
	// Storage .Get returns an empty record (no error) for soft-deleted or
	// scope-mismatched rows. Convert that into NotFound so the handler cannot
	// leak a blank response under the wrong project/collection.
	if docModel.ID.String() == "" {
		return nil, sderrors.NewNotFoundError(fmt.Errorf("document %s not found in collection %s", did.String(), cid.String()))
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

// DeleteDocument removes the named document; returns NotFound if it does not
// exist in the given (project, collection) scope.
func (r *RestColServiceServerService) DeleteDocument(ctx context.Context, req *apppb.DeleteDocumentRequest) (*apppb.DeleteDocumentResponse, error) {
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if req.CollectionId == "" {
		return nil, sderrors.NewBadParamsError(errors.New("collection_id is required"))
	}
	if req.DocumentId == "" {
		return nil, sderrors.NewBadParamsError(errors.New("document_id is required"))
	}
	cid := collectionsmodel.NewCollectionIDFromStr(req.CollectionId)
	did, err := documentsmodel.Parse(req.DocumentId)
	if err != nil {
		return nil, sderrors.NewBadParamsError(fmt.Errorf("invalid document_id: %w", err))
	}

	if _, err := r.documentCURD.Get(ctx, "", projectId, cid, did); err != nil {
		return nil, err
	}
	if err := r.documentCURD.Delete(ctx, "", projectId, cid, did); err != nil {
		return nil, err
	}
	return &apppb.DeleteDocumentResponse{}, nil
}

// QueryDocument returns a page of documents matching the request's time range
// and limit.
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
		filteredDoc, err := r.filterDocWithSelectedFields(doc, req.FieldSelectors)
		if err != nil {
			return nil, err
		}
		resp.Docs = append(resp.Docs, &apppb.GetDocumentResponse{
			XMetadata: documentsmodel.NewPbDocumentMetadata(doc),
			Data:      filteredDoc,
		})
	}
	return resp, nil
}

// QueryDocumentsStream streams matching documents. In follow-up mode it keeps
// polling until the client cancels the context.
func (r *RestColServiceServerService) QueryDocumentsStream(req *apppb.QueryDocumentStreamRequest, stream apppb.RestColService_QueryDocumentsStreamServer) error {
	ctx := stream.Context()
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return err
	}
	cid := collectionsmodel.NewCollectionIDFromStr(req.CollectionId)
	startedAt := req.SinceTs
	endedAt := req.EndedAt

	needsFollowUp := req.FollowUpMode != nil && *req.FollowUpMode

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
		if len(queryDocs) > 0 {
			startedAt = timestamppb.New(queryDocs[len(queryDocs)-1].CreatedAt)
		}
		select {
		case <-ctx.Done():
			r.log.Info("query is done as ctx is done\n")
			return nil
		case <-time.After(1 * time.Second):
		}
	}
}

// filterDocWithSelectedFields returns either the full document value or a
// pruned view containing only the field paths named in selectedFields.
func (r *RestColServiceServerService) filterDocWithSelectedFields(doc *documentsmodel.ModelDocument, selectedFields []string) (*structpb.Value, error) {
	r.log.Info("query doc with fields:%+v\n", selectedFields)

	if doc.Data == nil {
		return nil, nil
	}
	if len(selectedFields) == 0 {
		return structpb.NewValue(doc.Data.MapValue)
	}

	modelSchema, err := r.schemaBuilder.Flatten(doc.Data.MapValue)
	if err != nil {
		return nil, err
	}

	lookup := make(map[string]struct{}, len(selectedFields))
	for _, f := range selectedFields {
		lookup[strings.ToLower(f)] = struct{}{}
	}

	var fieldsInSelected []*collectionsmodel.ModelFieldSchema
	for _, dataField := range modelSchema.Fields {
		if _, ok := lookup[dataField.FieldName.String()]; ok {
			fieldsInSelected = append(fieldsInSelected, dataField)
		}
	}
	structWithSelectedFields, err := schemafinder.Build(fieldsInSelected)
	if err != nil {
		return nil, err
	}
	return structpb.NewValue(structWithSelectedFields)
}

// makeQueryConditioner translates optional proto filters into the storage
// layer's variadic query options.
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
