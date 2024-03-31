package storagedocuments

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	storagepostgres "github.com/sdinsure/agent/pkg/storage/postgres"

	appmodeldocuments "github.com/footprintai/restcol/pkg/models/documents"
	"github.com/footprintai/restcol/pkg/storage"
)

type DocumentCURD struct {
	*storagepostgres.PostgresDb
}

func NewDocumentCURD(db *storagepostgres.PostgresDb) *DocumentCURD {
	return &DocumentCURD{
		PostgresDb: db,
	}
}

func (c *DocumentCURD) AutoMigrate() error {
	tables := []interface{}{
		&appmodeldocuments.ModelDocument{},
	}
	return storage.WrapStorageError(c.PostgresDb.AutoMigrate(tables))
}

func (c *DocumentCURD) Write(ctx context.Context, tableName string, record *appmodeldocuments.ModelDocument) error {
	err := c.With(ctx, tableName).Clauses(clause.OnConflict{UpdateAll: true}).Create(record).Error
	return storage.WrapStorageError(err)
}

func (c *DocumentCURD) Update(ctx context.Context, tableName string, record *appmodeldocuments.ModelDocument) error {
	err := c.With(ctx, tableName).Session(&gorm.Session{FullSaveAssociations: true}).Updates(record).Error
	return storage.WrapStorageError(err)
}

func (c *DocumentCURD) Get(ctx context.Context, tableName string, did appmodeldocuments.DocumentID) (*appmodeldocuments.ModelDocument, error) {
	record := &appmodeldocuments.ModelDocument{}
	err := c.With(ctx, tableName).Where("id = ?", did.String()).Find(record).Error
	if err != nil {
		return nil, storage.WrapStorageError(err)
	}
	return record, nil
}
