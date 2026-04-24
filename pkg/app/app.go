// Package app implements the gRPC RestColService handlers. Handlers are split
// by resource into app.go (shared plumbing + swagger), collections.go, and
// documents.go.
package app

import (
	"context"
	"errors"

	"github.com/sdinsure/agent/pkg/logger"
	sdinsureruntime "github.com/sdinsure/agent/pkg/runtime"
	"google.golang.org/genproto/googleapis/api/httpbody"

	apppb "github.com/footprintai/restcol/api/pb"
	collectionsmodel "github.com/footprintai/restcol/pkg/models/collections"
	projectsmodel "github.com/footprintai/restcol/pkg/models/projects"
	schemafinder "github.com/footprintai/restcol/pkg/schema"
	collectionsstorage "github.com/footprintai/restcol/pkg/storage/collections"
	documentsstorage "github.com/footprintai/restcol/pkg/storage/documents"
	collectionsswagger "github.com/footprintai/restcol/pkg/swagger/collections"
)

// RestColServiceServerService implements apppb.RestColServiceServer.
type RestColServiceServerService struct {
	apppb.UnimplementedRestColServiceServer

	log            logger.Logger
	collectionCURD *collectionsstorage.CollectionCURD
	documentCURD   *documentsstorage.DocumentCURD

	schemaBuilder *schemafinder.SchemaBuilder

	defaultProjectResolver sdinsureruntime.ProjectResolver
}

// NewRestColServiceServerService wires a new service handler. Call
// SetDefaultProjectResolver before serving requests.
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

// SetDefaultProjectResolver installs the resolver that maps incoming requests
// to a project tenant. Without it, every handler returns an error.
func (r *RestColServiceServerService) SetDefaultProjectResolver(projectResolver sdinsureruntime.ProjectResolver) {
	r.defaultProjectResolver = projectResolver
}

// getProjectIdFromCtx extracts the tenant project from the request context
// populated by identity middleware.
func (r *RestColServiceServerService) getProjectIdFromCtx(ctx context.Context) (projectsmodel.ProjectID, error) {
	if r.defaultProjectResolver == nil {
		r.log.Error("getProjectIdFromCtx: no valid project resolver\n")
		return projectsmodel.ProjectID("invalid"), errors.New("no project resolver")
	}
	projectInfor, found := r.defaultProjectResolver.ProjectInfo(ctx)
	if !found {
		r.log.Info("no valid project id found\n")
		return "", nil
	}
	rawPid, err := projectInfor.GetProjectID()
	if err != nil {
		r.log.Error("getProjectIdFromCtx: get projectid failed, err:%+v\n", err)
		return "", err
	}
	return projectsmodel.ProjectID(rawPid), nil
}

// GetSwaggerDoc renders the OpenAPI document for one or all collections of the
// caller's project.
func (r *RestColServiceServerService) GetSwaggerDoc(ctx context.Context, req *apppb.GetSwaggerDocRequest) (*httpbody.HttpBody, error) {
	projectId, err := r.getProjectIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var collectionList []*collectionsmodel.ModelCollection
	if len(req.CollectionId) > 0 {
		cid := collectionsmodel.NewCollectionIDFromStr(req.CollectionId)
		selected, err := r.collectionCURD.GetLatestSchema(ctx, "", projectId, cid)
		if err != nil {
			return nil, err
		}
		collectionList = append(collectionList, selected)
	} else {
		collectionList, err = r.collectionCURD.ListByProjectID(ctx, "", projectId)
		if err != nil {
			return nil, err
		}
	}
	colSwaggerDoc := collectionsswagger.NewCollectionSwaggerDoc(collectionList...)
	docStr, err := colSwaggerDoc.RenderDoc()
	if err != nil {
		return nil, err
	}
	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        []byte(docStr),
	}, nil
}
