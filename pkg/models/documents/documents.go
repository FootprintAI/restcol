package modeldocuments

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	modelcollections "github.com/footprintai/restcol/pkg/models/collections"
	modelprojects "github.com/footprintai/restcol/pkg/models/projects"
)

type ModelDocument struct {
	ID DocumentID `gorm:"column:id;primaryKey;type:string;"`

	// The docScope composite index backs ModelDocument.Query,
	// CountByCollection, and DeleteByCollection — all of which filter by
	// (project, collection) and sort by created_at.
	CreatedAt time.Time      `gorm:"column:created_at;index:docScope,priority:3"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	// CreatedBy is the principal that first wrote the document; UpdatedBy is
	// the one that wrote it last.
	//
	// Two columns because writing an existing documentId is an UPSERT: with
	// only a creator, a document overwritten by a DIFFERENT principal would
	// still report the original author and the second write would leave no
	// trace of who made it.
	//
	// They must move independently, and the upsert in DocumentCURD.Write is
	// what makes that true: updated_by is in its DoUpdates list and created_by
	// is deliberately NOT, exactly as with updated_at and created_at. If
	// created_by ever joined that list, both fields would report the most
	// recent writer and the pair would answer nothing.
	//
	// Empty means unattributed - a deployment with no caller resolver wired
	// records no writer, which is a different thing from a writer named
	// "anonymous".
	CreatedBy string `gorm:"column:created_by"`
	UpdatedBy string `gorm:"column:updated_by"`

	Data *ModelDocumentData `gorm:"column:data;type:jsonb"`

	ModelCollectionID modelcollections.CollectionID `gorm:"column:model_collection_id;primaryKey;index:docScope,priority:2"` // foreign key to model collection
	ModelCollection   modelcollections.ModelCollection

	ModelProjectID modelprojects.ProjectID `gorm:"column:model_project_id;primaryKey;index:docScope,priority:1"` // foreign key to model project
	ModelProject   modelprojects.ModelProject
}

func (m ModelDocument) TableName() string {
	return "restcol-documents"
}

type DocumentID struct {
	S string
}

func (d DocumentID) String() string {
	return d.S
}

var (
	_ driver.Valuer = DocumentID{}
	_ sql.Scanner   = &DocumentID{}
)

func (d DocumentID) Value() (driver.Value, error) {
	return d.String(), nil
}

func (d *DocumentID) Scan(value interface{}) error {
	_, isStringType := value.(string)
	if !isStringType {
		return fmt.Errorf("db.model: invalid type, expect string")
	}
	did, err := Parse(value.(string))
	if err != nil {
		return err
	}
	(*d) = did
	return nil
}

func Parse(s string) (DocumentID, error) {
	//innerUuid, err := uuid.Parse(s)
	//if err != nil {
	//	fmt.Printf("parse doc uuid failed, err:+%v, s:%s\n", err, s)
	//	return DocumentID{}, err
	//}
	return DocumentID{S: s}, nil

}

func NewDocumentID() DocumentID {
	uid, _ := uuid.NewV7()
	return DocumentID{S: uid.String()}
}

type ModelDocumentData struct {
	MapValue map[string]interface{} `json:"json_value"`
}

func NewModelDocumentData(v map[string]interface{}) *ModelDocumentData {
	return &ModelDocumentData{
		MapValue: v,
	}
}

var (
	_ driver.Valuer = ModelDocumentData{}
	_ sql.Scanner   = &ModelDocumentData{}
)

func (d ModelDocumentData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *ModelDocumentData) Scan(value interface{}) error {
	_, isByte := value.([]byte)
	if !isByte {
		return fmt.Errorf("db.model: invalid type, expect []byte")
	}
	return json.Unmarshal(value.([]byte), d)
}
