package storageprojects

import (
	"context"

	storagepostgres "github.com/sdinsure/agent/pkg/storage/postgres"
	"gorm.io/gorm/clause"

	appmodelprojects "github.com/footprintai/restcol/pkg/models/projects"
	"github.com/footprintai/restcol/pkg/storage"
)

type ProjectCURD struct {
	*storagepostgres.PostgresDb
}

func NewProjectCURD(db *storagepostgres.PostgresDb) *ProjectCURD {
	return &ProjectCURD{
		PostgresDb: db,
	}
}

func (c *ProjectCURD) AutoMigrate() error {
	tables := []interface{}{
		&appmodelprojects.ModelProject{},
	}
	return storage.WrapStorageError(c.PostgresDb.AutoMigrate(tables))
}

func (c *ProjectCURD) Write(ctx context.Context, tableName string, record *appmodelprojects.ModelProject) error {
	err := c.With(ctx, tableName).Clauses(clause.OnConflict{DoNothing: true}).Create(record).Error
	return storage.WrapStorageError(err)
}

func (c *ProjectCURD) Get(ctx context.Context, tableName string, pid appmodelprojects.ProjectID) (*appmodelprojects.ModelProject, error) {
	record := &appmodelprojects.ModelProject{}
	err := c.With(ctx, tableName).Where(" id = ? ", pid.String()).Create(record).Error
	return record, storage.WrapStorageError(err)
}
