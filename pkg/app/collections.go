package app

import (
	"context"
	"errors"
	"fmt"

	sderrors "github.com/sdinsure/agent/pkg/errors"

	apppb "github.com/footprintai/restcol/api/pb"
	collectionsmodel "github.com/footprintai/restcol/pkg/models/collections"
	projectsmodel "github.com/footprintai/restcol/pkg/models/projects"
)

// CreateCollection creates a new collection under the caller's project.
func (r *RestColServiceServerService) CreateCollection(ctx context.Context, req *apppb.CreateCollectionRequest) (*apppb.CreateCollectionResponse, error) {
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	cid := collectionsmodel.NewCollectionID()
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

// ListCollections is not yet implemented.
func (r *RestColServiceServerService) ListCollections(ctx context.Context, req *apppb.ListCollectionsRequest) (*apppb.ListCollectionsResponse, error) {
	return nil, sderrors.NewNotImplError(errors.New("not implemented"))
}

// GetCollection returns metadata plus the latest schema for the requested
// collection. A not-found is mapped to an empty response (not an error) to
// match existing client expectations.
func (r *RestColServiceServerService) GetCollection(ctx context.Context, req *apppb.GetCollectionRequest) (*apppb.GetCollectionResponse, error) {
	if len(req.CollectionId) == 0 {
		return nil, sderrors.NewBadParamsError(errors.New("missing required field"))
	}
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	cid := collectionsmodel.NewCollectionIDFromStr(req.CollectionId)
	mc, err := r.collectionCURD.GetLatestSchema(ctx, "", projectId, cid)
	if err != nil {
		if ismyerr, myerr := sderrors.As(err); ismyerr && myerr.Code() == sderrors.CodeNotFound {
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

// DeleteCollection soft-deletes a collection. If the collection still contains
// documents the call fails with FailedPrecondition unless req.Force is true, in
// which case all documents in the collection are soft-deleted first.
func (r *RestColServiceServerService) DeleteCollection(ctx context.Context, req *apppb.DeleteCollectionRequest) (*apppb.DeleteCollectionResponse, error) {
	if len(req.CollectionId) == 0 {
		return nil, sderrors.NewBadParamsError(errors.New("collection_id is required"))
	}
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	cid := collectionsmodel.NewCollectionIDFromStr(req.CollectionId)

	if _, err := r.collectionCURD.Get(ctx, "", projectId, cid, collectionsmodel.NullSchemaID); err != nil {
		return nil, err
	}

	count, err := r.documentCURD.CountByCollection(ctx, "", projectId, cid)
	if err != nil {
		return nil, err
	}
	if count > 0 && !req.Force {
		return nil, sderrors.NewStatusConflicted(fmt.Errorf("collection %s contains %d documents; pass force=true to cascade-delete", cid.String(), count))
	}
	if count > 0 {
		if err := r.documentCURD.DeleteByCollection(ctx, "", projectId, cid); err != nil {
			return nil, err
		}
	}
	if err := r.collectionCURD.Delete(ctx, "", projectId, cid); err != nil {
		return nil, err
	}
	return &apppb.DeleteCollectionResponse{}, nil
}

// getCollectionIDFromSchemas auto-provisions a collection when a document is
// written without an explicit collection ID.
//
// TODO: match incoming data against existing collections by schema similarity
// instead of always creating a new collection.
func (r *RestColServiceServerService) getCollectionIDFromSchemas(ctx context.Context, projectId projectsmodel.ProjectID) (collectionsmodel.CollectionID, error) {
	cid := collectionsmodel.NewCollectionID()
	mc := collectionsmodel.NewModelCollection(
		projectId,
		cid,
		apppb.CollectionType_COLLECTION_TYPE_REGULAR_FILES,
		"auto created collection",
		[]*collectionsmodel.ModelSchema{},
	)
	if err := r.collectionCURD.Write(ctx, "", &mc); err != nil {
		return collectionsmodel.CollectionID(""), err
	}
	return cid, nil
}
