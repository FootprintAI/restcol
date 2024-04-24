package modelcollections

import (
	"database/sql/driver"
	"errors"
	"time"

	"gorm.io/gorm"

	apppb "github.com/footprintai/restcol/api/pb"
	modelprojects "github.com/footprintai/restcol/pkg/models/projects"
)

type ModelCollection struct {
	ID CollectionID `gorm:"column:id;primarykey;type:string;"`

	Type ModelCollectionType `gorm:"column:type;type:int;"`
	// Summary is a summary for the collections
	Summary   string         `gorm:"column:summary"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	Schemas        []ModelSchema           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // associated to many schemes
	ModelProjectID modelprojects.ProjectID // foreign key to model project
	ModelProject   modelprojects.ModelProject
}

func NewModelCollection(
	projectId modelprojects.ProjectID,
	id CollectionID,
	t apppb.CollectionType,
	summary string,
	schemas []ModelSchema,
) ModelCollection {
	return ModelCollection{
		ID:             id,
		Type:           ModelCollectionType(t),
		Summary:        summary,
		Schemas:        schemas,
		ModelProjectID: projectId,
	}
}

func (m ModelCollection) TableName() string {
	return "restcol-collections"
}

type ModelCollectionType apppb.CollectionType

func (m ModelCollectionType) Proto() apppb.CollectionType {
	return apppb.CollectionType(m)
}

func (m ModelCollectionType) Number() int {
	return int(apppb.CollectionType(m).Number())
}

func (m ModelCollectionType) Value() (driver.Value, error) {
	return m.Number(), nil
}

func (m *ModelCollectionType) Scan(in any) error {
	if int64Val, isInt64 := in.(int64); isInt64 {
		(*m) = ModelCollectionType(apppb.CollectionType(int64Val))
		return nil
	}
	return errors.New("shema: require int64")
}
