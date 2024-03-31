package modeldocuments

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	modelcollections "github.com/footprintai/restcol/pkg/models/collections"
	modelprojects "github.com/footprintai/restcol/pkg/models/projects"
)

type ModelDocument struct {
	ID DocumentID `gorm:"column:id;primarykey;type:string;"`

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	Data datatypes.JSON `gorm:"column:data"`

	ModelCollectionID modelcollections.CollectionID // foreign key to model collection
	ModelCollection   modelcollections.ModelCollection

	ModelProjectID modelprojects.ProjectID // foreigh key to model project
	ModelProject   modelprojects.ModelProject
}

type DocumentID uuid.UUID

func (d DocumentID) String() string {
	return uuid.UUID(d).String()
}

func Parse(s string) (DocumentID, error) {
	innerUuid, err := uuid.Parse(s)
	if err != nil {
		return DocumentID{}, err
	}
	return DocumentID(innerUuid), nil

}

func NewDocumentID() DocumentID {
	uid, _ := uuid.NewV7()
	return DocumentID(uid)
}
