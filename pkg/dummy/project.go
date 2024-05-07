package dummy

import (
	"context"

	projectsmodel "github.com/footprintai/restcol/pkg/models/projects"
	projectsstorage "github.com/footprintai/restcol/pkg/storage/projects"
)

var (
	DummyModelProject = projectsmodel.ModelProject{
		ID:   projectsmodel.NewProjectID(1001),
		Type: projectsmodel.ProxyProjectType,
	}
)

type DummyProject struct {
	projectcurd *projectsstorage.ProjectCURD
}

func NewDummyProject(projectcurd *projectsstorage.ProjectCURD) *DummyProject {
	return &DummyProject{
		projectcurd: projectcurd,
	}
}

func (d *DummyProject) Init(ctx context.Context) error {
	return d.projectcurd.Write(ctx, "", &DummyModelProject)
}

func (d *DummyProject) GetProject(ctx context.Context, pid projectsmodel.ProjectID) (*projectsmodel.ModelProject, error) {
	docModel, err := d.projectcurd.Get(ctx, "", pid)
	return docModel, err
}
