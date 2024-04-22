package storagedocuments

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	storagepostgres "github.com/sdinsure/agent/pkg/storage/postgres"

	appmodelcollections "github.com/footprintai/restcol/pkg/models/collections"
	appmodeldocuments "github.com/footprintai/restcol/pkg/models/documents"
	appmodelprojects "github.com/footprintai/restcol/pkg/models/projects"
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

func (c *DocumentCURD) BatchWrite(ctx context.Context, tableName string, records []*appmodeldocuments.ModelDocument) error {
	err := c.With(ctx, tableName).Clauses(clause.OnConflict{UpdateAll: true}).Create(&records).Error
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

type QueryConditioner interface {
	apply(*conditions)
}

type conditions struct {
	startedAt *time.Time
	endedAt   *time.Time
	limit     *int
}

func newConditions(cnds ...QueryConditioner) *conditions {
	c := &conditions{}
	for _, cnd := range cnds {
		cnd.apply(c)
	}
	return c
}

var (
	_ QueryConditioner = withStartedAt{}
)

type withStartedAt struct {
	startedAt time.Time
}

func (s withStartedAt) apply(c *conditions) {
	c.startedAt = &s.startedAt
}

func WithStartedAt(startedAt time.Time) withStartedAt {
	return withStartedAt{startedAt: startedAt}
}

var (
	_ QueryConditioner = withEndedAt{}
)

type withEndedAt struct {
	endedAt time.Time
}

func (s withEndedAt) apply(c *conditions) {
	c.endedAt = &s.endedAt
}

func WithEndedAt(endedAt time.Time) withEndedAt {
	return withEndedAt{endedAt: endedAt}
}

var (
	_ QueryConditioner = withLimit{}
)

type withLimit struct {
	count int
}

func (s withLimit) apply(c *conditions) {
	c.limit = &s.count
}

func WithLimitCount(count int) withLimit {
	return withLimit{count: count}
}

func (c *DocumentCURD) Query(ctx context.Context,
	tableName string,
	pid appmodelprojects.ProjectID,
	cid appmodelcollections.CollectionID,
	cnds ...QueryConditioner,
) ([]*appmodeldocuments.ModelDocument, error) {
	conditions := newConditions(cnds...)

	records := []*appmodeldocuments.ModelDocument{}
	projections := c.With(ctx, tableName).
		Where("model_project_id = ?", pid.String()).
		Where("model_collection_id = ?", cid.String()).
		Order("created_at asc")
	if conditions.startedAt != nil {
		projections.Where("created_at >= ?", *conditions.startedAt)
	}
	if conditions.endedAt != nil {
		projections.Where("created_at < ?", *conditions.endedAt)
	}
	if conditions.limit != nil {
		projections.Limit(*conditions.limit)
	} else {
		projections.Limit(30) // default limit = 30
	}
	err := projections.Find(&records).Error
	if err != nil {
		return nil, storage.WrapStorageError(err)
	}
	return records, nil
}
