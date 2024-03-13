package storageprojects

import (
	"context"

	"github.com/sdinsure/agent/pkg/errors"
	storageerrors "github.com/sdinsure/agent/pkg/storage/errors"
	storagepostgres "github.com/sdinsure/agent/pkg/storage/postgres"
	"gorm.io/gorm/clause"

	appmodelprojects "github.com/footprintai/restcol/pkg/models/projects"
)

type ProjectCURD struct {
	*storagepostgres.PostgresDb
}

func NewProjectCURD(db *storagepostgres.PostgresDb) *ProjectCURD {
	return &ProjectCURD{
		PostgresDb: db,
	}
}

func (c *ProjectCURD) AutoMigrate() *errors.Error {
	tables := []interface{}{
		&appmodelprojects.ModelProject{},
	}
	return storageerrors.WrapStorageError(c.PostgresDb.AutoMigrate(tables))
}

func (c *ProjectCURD) Write(ctx context.Context, tableName string, record *appmodelprojects.ModelProject) *errors.Error {
	err := c.With(ctx, tableName).Clauses(clause.OnConflict{DoNothing: true}).Create(record).Error
	return storageerrors.WrapStorageError(err)
}

func (c *ProjectCURD) Get(ctx context.Context, tableName string, pid appmodelprojects.ProjectID) (*appmodelprojects.ModelProject, *errors.Error) {
	record := &appmodelprojects.ModelProject{}
	err := c.With(ctx, tableName).Where(" id = ? ", pid.String()).Create(record).Error
	return record, storageerrors.WrapStorageError(err)
}
