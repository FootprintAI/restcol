package storagecollections

import (
	"context"

	"github.com/sdinsure/agent/pkg/errors"
	storageerrors "github.com/sdinsure/agent/pkg/storage/errors"
	storagepostgres "github.com/sdinsure/agent/pkg/storage/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	appmodelcollections "github.com/footprintai/restcol/pkg/models/collections"
)

type CollectionCURD struct {
	*storagepostgres.PostgresDb
}

func NewCollectionCURD(db *storagepostgres.PostgresDb) *CollectionCURD {
	return &CollectionCURD{
		PostgresDb: db,
	}
}

func (c *CollectionCURD) AutoMigrate() *errors.Error {
	tables := []interface{}{
		&appmodelcollections.ModelCollection{},
		&appmodelcollections.ModelSchema{},
		&appmodelcollections.ModelFieldSchema{},
	}
	return storageerrors.WrapStorageError(c.PostgresDb.AutoMigrate(tables))
}

func (c *CollectionCURD) Write(ctx context.Context, tableName string, record *appmodelcollections.ModelCollection) *errors.Error {
	err := c.With(ctx, tableName).Clauses(clause.OnConflict{UpdateAll: true}).Create(record).Error
	return storageerrors.WrapStorageError(err)
}

func (c *CollectionCURD) Update(ctx context.Context, tableName string, record *appmodelcollections.ModelCollection) *errors.Error {
	err := c.With(ctx, tableName).Session(&gorm.Session{FullSaveAssociations: true}).Updates(record).Error
	return storageerrors.WrapStorageError(err)
}

func (c *CollectionCURD) GetLatestSchema(ctx context.Context, tableName string, cid appmodelcollections.CollectionID) (*appmodelcollections.ModelCollection, *errors.Error) {
	s := &appmodelcollections.ModelSchema{}
	err := c.With(ctx, tableName).Where("model_collection_id = ?", cid.String()).Order("id desc").First(s).Error
	if err != nil {
		return nil, storageerrors.WrapStorageError(err)
	}
	return c.Get(ctx, tableName, cid, s.ID)
}

func (c *CollectionCURD) Get(ctx context.Context, tableName string, cid appmodelcollections.CollectionID, sid appmodelcollections.SchemaID) (*appmodelcollections.ModelCollection, *errors.Error) {
	record := &appmodelcollections.ModelCollection{}
	err := c.With(ctx, tableName).
		Preload("Schemas", func(db *gorm.DB) *gorm.DB {
			return db.Where(&appmodelcollections.ModelSchema{ID: sid})
		}).
		Preload("Schemas.Fields").Where("id = ?", cid.String()).Find(record).Error
	if err != nil {
		return nil, storageerrors.WrapStorageError(err)
	}
	return record, nil
}
