package storageprojects

import (
	"context"

	appmodelprojects "github.com/footprintai/restcol/pkg/models/projects"
	storagepostgres "github.com/sdinsure/agent/pkg/storage/postgres"
)

func TestProjectSuite(postgrescli *storagepostgres.PostgresDb) (regularProject *appmodelprojects.ModelProject, proxyProject *appmodelprojects.ModelProject, retErr error) {

	ctx := context.Background()
	pcrud := NewProjectCURD(postgrescli)
	if retErr = pcrud.AutoMigrate(); retErr != nil {
		return
	}

	regularProject = &appmodelprojects.ModelProject{
		ID:   appmodelprojects.NewProjectID(1),
		Type: appmodelprojects.RegularProjectType,
	}
	if retErr = pcrud.Write(ctx, "", regularProject); retErr != nil {
		return
	}

	proxyProject = &appmodelprojects.ModelProject{
		ID:   appmodelprojects.NewProjectID(2),
		Type: appmodelprojects.ProxyProjectType,
	}
	if retErr = pcrud.Write(ctx, "", proxyProject); retErr != nil {
		return
	}
	return
}
