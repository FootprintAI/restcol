// Package bootstrap seeds the database with the baseline records required for
// the server to start. It currently provisions the default project that backs
// anonymous authentication so requests without an explicit project have a
// tenant to attach to.
package bootstrap

import (
	"context"

	projectsmodel "github.com/footprintai/restcol/pkg/models/projects"
	projectsstorage "github.com/footprintai/restcol/pkg/storage/projects"
)

// DefaultModelProject is the project used as the fallback tenant for anonymous
// requests. Its ID is stable so swagger/OpenAPI URLs like
// /v1/projects/1001/apidoc remain valid across restarts.
var DefaultModelProject = projectsmodel.ModelProject{
	ID:   projectsmodel.NewProjectID(1001),
	Type: projectsmodel.ProxyProjectType,
}

// DefaultProject upserts and retrieves DefaultModelProject.
type DefaultProject struct {
	projectcurd *projectsstorage.ProjectCURD
}

// NewDefaultProject wires a DefaultProject against the given storage handle.
func NewDefaultProject(projectcurd *projectsstorage.ProjectCURD) *DefaultProject {
	return &DefaultProject{
		projectcurd: projectcurd,
	}
}

// Init persists DefaultModelProject so downstream code can resolve it by ID.
// Safe to call on every startup.
func (d *DefaultProject) Init(ctx context.Context) error {
	return d.projectcurd.Write(ctx, "", &DefaultModelProject)
}

// GetProject returns the project record identified by pid.
func (d *DefaultProject) GetProject(ctx context.Context, pid projectsmodel.ProjectID) (*projectsmodel.ModelProject, error) {
	return d.projectcurd.Get(ctx, "", pid)
}
