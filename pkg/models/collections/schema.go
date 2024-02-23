package modelcollections

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cthulhu/jsonpath"
	"gorm.io/gorm"
)

type SchemaID int

func (s SchemaID) String() string {
	return fmt.Sprintf("%d", s)
}

type ModelSchema struct {
	ID        SchemaID       `gorm:"column:id;primarykey;type:int;uniqueIndex:cidsid;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	Fields            []ModelFieldSchema `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ModelCollectionID CollectionID       // foreigh key to ModelCollection -> ID
}

func (m ModelSchema) TableName() string {
	return "restcol-collections-schema"
}

type FieldID int

func (f FieldID) String() string {
	return fmt.Sprintf("%d", f)
}

type fieldValueType interface {
	valueType() string
}

type SwagValueType string

var (
	_ fieldValueType = new(SwagValueType)
)

func (f SwagValueType) valueType() string {
	return string(f)
}

var (
	StringSwagValueType  SwagValueType = "string"
	NumberSwagValueType  SwagValueType = "number"
	IntegerSwagValueType SwagValueType = "integer"
	BoolSwagValueType    SwagValueType = "bool"
)

type fieldValueValue interface {
}

type SwagValueValue struct {
	StringValue  *string  `json:"str,omitempty"`
	NumberValue  *float64 `json:"num,omitempty"`
	IntegerValue *int     `json:"int,omitempty"`
	BoolValue    *bool    `json:"bool,omitempty"`
}

func Must(s SwagValueValue, e error) SwagValueValue {
	if e != nil {
		panic(e)
	}
	return s
}

func NewSwagValue(v any) (SwagValueValue, error) {
	switch v.(type) {
	case int:
		i := v.(int)
		return SwagValueValue{
			IntegerValue: &i,
		}, nil
	case float64:
		f := v.(float64)
		return SwagValueValue{
			NumberValue: &f,
		}, nil
	case bool:
		b := v.(bool)
		return SwagValueValue{
			BoolValue: &b,
		}, nil
	case string:
		s := v.(string)
		return SwagValueValue{
			StringValue: &s,
		}, nil
	default:
		return SwagValueValue{}, errors.New("Schema: no available type for swag value")
	}
}

var (
	_ sql.Scanner   = &SwagValueValue{}
	_ driver.Valuer = SwagValueValue{}
)

func (s *SwagValueValue) Scan(in any) error {
	if _, isBytes := in.([]byte); isBytes {
		return json.Unmarshal(in.([]byte), s)
	}
	return errors.New("shema: require []bytes")

}

func (s SwagValueValue) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type ModelFieldSchema struct {
	ID        FieldID        `gorm:"column:id;primarykey;type:int;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	FieldName      string         `gorm:"column:field_name"` // dot concated path, a.b.c represents a -> b -> c path
	FieldValueType SwagValueType  `gorm:"column:value_type"`
	FieldExample   SwagValueValue `gorm:"column:value_example;type:jsonb"`

	ModelSchemaID SchemaID // foreign key to ModelSchema -> ID
}

func (m ModelFieldSchema) TableName() string {
	return "restcol-collections-schematable"
}

type ModelFieldsSchema []ModelFieldSchema

func (m ModelFieldsSchema) ToJSON(dotPrefixs ...string) ([]byte, error) {

	withPrefix := func(fieldName string) string {
		if len(dotPrefixs) == 0 {
			return fieldName
		}
		prefixAndFieldName := append(dotPrefixs, fieldName)
		return strings.Join(prefixAndFieldName, ".")
	}

	// convert fields into a map with a common dotPrefix
	fieldsMap := make(map[string]string)
	for _, field := range m {
		fieldName := withPrefix(field.FieldName)
		if field.FieldValueType == NumberSwagValueType && field.FieldExample.NumberValue != nil {
			fieldsMap[fmt.Sprintf("%s.num()", fieldName)] = fmt.Sprintf("%f", *field.FieldExample.NumberValue)
		} else if field.FieldValueType == IntegerSwagValueType && field.FieldExample.IntegerValue != nil {
			fieldsMap[fmt.Sprintf("%s.num()", fieldName)] = fmt.Sprintf("%d", *field.FieldExample.IntegerValue)
		} else if field.FieldValueType == BoolSwagValueType && field.FieldExample.BoolValue != nil {
			fieldsMap[fmt.Sprintf("%s.bool()", fieldName)] = fmt.Sprintf("%b", *field.FieldExample.BoolValue)
		} else if field.FieldValueType == StringSwagValueType && field.FieldExample.StringValue != nil {
			fieldsMap[fieldName] = *field.FieldExample.StringValue
		} else {
			return nil, errors.New("fieldschema: invalid condition")
		}
	}
	return jsonpath.Marshal(fieldsMap)
}
