package getter

import (
	"context"
	"errors"

	sdinsureruntime "github.com/sdinsure/agent/pkg/runtime"

	projectsmodels "github.com/footprintai/restcol/pkg/models/projects"
	projectsstorage "github.com/footprintai/restcol/pkg/storage/projects"
)

func NewRuntimeProjectGetter(projectCURD *projectsstorage.ProjectCURD) *RuntimeProjectGetter {
	return &RuntimeProjectGetter{
		projectCURD: projectCURD,
	}
}

type RuntimeProjectGetter struct {
	projectCURD *projectsstorage.ProjectCURD
}

var (
	_ sdinsureruntime.ProjectGetter = &RuntimeProjectGetter{}
)

func (p *RuntimeProjectGetter) GetProject(ctx context.Context, projectId string) (sdinsureruntime.ProjectInfor, error) {
	modelProject, err := p.projectCURD.Get(ctx, "", projectsmodels.ProjectID(projectId))
	if err != nil {
		return sdinsureruntime.NewInvalidProjectInfor(), err
	}
	return projectInfor{modleProject: modelProject}, nil
}

var (
	_ sdinsureruntime.ProjectInfor = &projectInfor{}
)

type projectInfor struct {
	modleProject *projectsmodels.ModelProject
}

func (p projectInfor) GetProjectID() (string, error) {
	return p.modleProject.ID.String(), nil
}

func (p projectInfor) GetProject(v any) error {
	return errors.New("not impl")
}
