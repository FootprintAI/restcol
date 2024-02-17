package modelcollections

import (
	"time"

	"gorm.io/gorm"
)

type ModelCollection struct {
	ID CollectionID `gorm:"column:id;primarykey;type:string;"`

	// Summary is a summary for the collections
	Summary   string         `gorm:"column:summary"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	Schemas []ModelSchema `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // associated to many schemes
}

func (m ModelCollection) TableName() string {
	return "restcol-collections"
}
