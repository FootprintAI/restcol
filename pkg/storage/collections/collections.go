package storagecollections

import (
	"context"

	sdinsureerrors "github.com/sdinsure/agent/pkg/errors"
	storagepostgres "github.com/sdinsure/agent/pkg/storage/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	appmodelcollections "github.com/footprintai/restcol/pkg/models/collections"
	appmodelprojects "github.com/footprintai/restcol/pkg/models/projects"
	"github.com/footprintai/restcol/pkg/storage"
)

type CollectionCURD struct {
	*storagepostgres.PostgresDb
}

func NewCollectionCURD(db *storagepostgres.PostgresDb) *CollectionCURD {
	return &CollectionCURD{
		PostgresDb: db,
	}
}

func (c *CollectionCURD) AutoMigrate() error {
	tables := []interface{}{
		&appmodelcollections.ModelCollection{},
		&appmodelcollections.ModelSchema{},
		&appmodelcollections.ModelFieldSchema{},
	}
	return storage.WrapStorageError(c.PostgresDb.AutoMigrate(tables))
}

func (c *CollectionCURD) Write(ctx context.Context, tableName string, record *appmodelcollections.ModelCollection) error {
	err := c.With(ctx, tableName).Clauses(clause.OnConflict{UpdateAll: true}).Create(record).Error
	return storage.WrapStorageError(err)
}

func (c *CollectionCURD) Update(ctx context.Context, tableName string, record *appmodelcollections.ModelCollection) error {
	err := c.With(ctx, tableName).Session(&gorm.Session{FullSaveAssociations: true}).Updates(record).Error
	return storage.WrapStorageError(err)
}

func (c *CollectionCURD) GetLatestSchema(ctx context.Context, tableName string, pid appmodelprojects.ProjectID, cid appmodelcollections.CollectionID) (*appmodelcollections.ModelCollection, error) {
	s := &appmodelcollections.ModelSchema{}
	err := c.With(ctx, tableName).Where("model_collection_id = ?", cid.String()).Order("id desc").First(s).Error
	if err != nil {
		wrappedErr := storage.WrapStorageError(err)
		ismyerr, myerr := sdinsureerrors.As(wrappedErr)
		if ismyerr && myerr.Code() == sdinsureerrors.CodeNotFound {
			return c.Get(ctx, tableName, pid, cid, appmodelcollections.NullSchemaID)
		}
		return nil, wrappedErr
	}
	return c.Get(ctx, tableName, pid, cid, s.ID)
}

func (c *CollectionCURD) ListByProjectID(ctx context.Context, tableName string, pid appmodelprojects.ProjectID) ([]*appmodelcollections.ModelCollection, error) {
	var cs []*appmodelcollections.ModelCollection
	err := c.With(ctx, tableName).
		Preload("Schemas", func(db *gorm.DB) *gorm.DB {
			return db.Order("id desc").Limit(1)
		}).
		Preload("Schemas.Fields").
		Where("model_project_id = ?", pid.String()).Order("id desc").Find(&cs).Error
	return cs, storage.WrapStorageError(err)
}

func (c *CollectionCURD) Get(ctx context.Context, tableName string, pid appmodelprojects.ProjectID, cid appmodelcollections.CollectionID, sid appmodelcollections.SchemaID) (*appmodelcollections.ModelCollection, error) {
	record := &appmodelcollections.ModelCollection{}
	db := c.With(ctx, tableName)
	if sid != appmodelcollections.NullSchemaID {
		db = db.Preload("Schemas", func(db *gorm.DB) *gorm.DB {
			return db.Where(&appmodelcollections.ModelSchema{ID: sid})
		}).Preload("Schemas.Fields")
	}
	if err := db.
		Where("id = ?", cid.String()).
		Where("model_project_id = ?", pid.String()).
		First(record).Error; err != nil {
		return nil, storage.WrapStorageError(err)
	}
	return record, nil
}
