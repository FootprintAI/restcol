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

// Visibility reports the project's visibility tag, and restcol has none.
//
// ModelProject carries an ID and a ProjectType (regular/external/proxy) —
// nothing about who may see it — and authorization here is
// appauthz.AllowEveryOne, so there is no visibility to report and nothing that
// would consume it. The interface defines the empty string as "unknown /
// unresolved", which is exactly the truthful answer; sdinsure's own
// invalidProjectInfor returns the same.
//
// Inventing a value ("PUBLIC") would be worse than saying nothing: a consumer
// that later reads this on an authorization path would be told every restcol
// project is world-visible on no evidence.
func (p projectInfor) Visibility() string {
	return ""
}
