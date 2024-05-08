package modelcollections

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"

	"github.com/cthulhu/jsonpath"
	"google.golang.org/protobuf/types/known/structpb"
	_ "google.golang.org/protobuf/types/known/structpb"
	"gorm.io/gorm"

	apppb "github.com/footprintai/restcol/api/pb"
	dotnotation "github.com/footprintai/restcol/pkg/notation/dot"
)

var (
	NullSchemaID SchemaID = SchemaID(-1)
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

	Fields            []*ModelFieldSchema `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ModelCollectionID CollectionID        // foreigh key to ModelCollection -> ID
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
	NoneSwagValueType    SwagValueType = "none"
	NullSwagValueType    SwagValueType = "null"
	StringSwagValueType  SwagValueType = "string"
	NumberSwagValueType  SwagValueType = "number"
	IntegerSwagValueType SwagValueType = "integer"
	BoolSwagValueType    SwagValueType = "bool"
	ObjectSwagValueType  SwagValueType = "object"
	ArraySwagValueType   SwagValueType = "array"
)

func (f SwagValueType) Proto() apppb.SchemaFieldDataType {
	switch f {
	case StringSwagValueType:
		return apppb.SchemaFieldDataType_SCHEMA_FIELD_DATA_TYPE_STRING
	case NumberSwagValueType:
		return apppb.SchemaFieldDataType_SCHEMA_FIELD_DATA_TYPE_NUMBER
	case IntegerSwagValueType:
		return apppb.SchemaFieldDataType_SCHEMA_FIELD_DATA_TYPE_INTEGER
	case BoolSwagValueType:
		return apppb.SchemaFieldDataType_SCHEMA_FIELD_DATA_TYPE_BOOL
	case ObjectSwagValueType:
		return apppb.SchemaFieldDataType_SCHEMA_FIELD_DATA_TYPE_OBJECT
	case ArraySwagValueType:
		return apppb.SchemaFieldDataType_SCHEMA_FIELD_DATA_TYPE_ARRAY
	default:
		return apppb.SchemaFieldDataType_SCHEMA_FIELD_DATA_TYPE_NONE
	}
}

func NewSwaggerValueType(pbDataType apppb.SchemaFieldDataType) SwagValueType {
	if pbDataType == apppb.SchemaFieldDataType_SCHEMA_FIELD_DATA_TYPE_STRING {
		return StringSwagValueType
	}
	if pbDataType == apppb.SchemaFieldDataType_SCHEMA_FIELD_DATA_TYPE_NUMBER {
		return NumberSwagValueType
	}
	if pbDataType == apppb.SchemaFieldDataType_SCHEMA_FIELD_DATA_TYPE_INTEGER {
		return IntegerSwagValueType
	}
	if pbDataType == apppb.SchemaFieldDataType_SCHEMA_FIELD_DATA_TYPE_BOOL {
		return BoolSwagValueType
	}
	if pbDataType == apppb.SchemaFieldDataType_SCHEMA_FIELD_DATA_TYPE_OBJECT {
		return ObjectSwagValueType
	}
	if pbDataType == apppb.SchemaFieldDataType_SCHEMA_FIELD_DATA_TYPE_ARRAY {
		return ArraySwagValueType
	}
	return NoneSwagValueType
}

type SwagValueValue structpb.Value

func (s SwagValueValue) Interface() interface{} {
	return s.Proto().AsInterface()
}

func Must(s *SwagValueValue, e error) *SwagValueValue {
	if e != nil {
		panic(e)
	}
	return s
}

func NewSwagValue(v any) (*SwagValueValue, error) {
	pbValue, err := structpb.NewValue(wrapSlice(v))
	if err != nil {
		return &SwagValueValue{}, err
	}
	swagVal := SwagValueValue(*pbValue)
	return &swagVal, nil
}

func wrapSlice(v any) any {
	s := reflect.ValueOf(v)
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		ret := make([]interface{}, s.Len())
		for i := 0; i < s.Len(); i++ {
			ret[i] = s.Index(i).Interface()
		}
		return ret
	}
	return v
}

func (s SwagValueValue) Proto() *structpb.Value {
	pbValue := structpb.Value(s)
	return &pbValue
}

func (s *SwagValueValue) Type() SwagValueType {
	pbValue := structpb.Value(*s)
	switch pbValue.Kind.(type) {
	case *structpb.Value_NullValue:
		return NullSwagValueType
	case *structpb.Value_BoolValue:
		return BoolSwagValueType
	case *structpb.Value_NumberValue:
		return NumberSwagValueType
	case *structpb.Value_StringValue:
		return StringSwagValueType
	case *structpb.Value_StructValue:
		return ObjectSwagValueType
	case *structpb.Value_ListValue:
		return ArraySwagValueType
	default:
		return NoneSwagValueType
	}
}

var (
	_ sql.Scanner   = &SwagValueValue{}
	_ driver.Valuer = SwagValueValue{}
)

func (s *SwagValueValue) Scan(in any) error {
	pbValue := &structpb.Value{}
	if err := pbValue.UnmarshalJSON(in.([]byte)); err != nil {
		return err
	}
	(*s) = SwagValueValue(*pbValue)
	return nil
}

func (s SwagValueValue) Value() (driver.Value, error) {
	return s.Proto().MarshalJSON()
}

type ModelFieldSchema struct {
	ID        FieldID        `gorm:"column:id;primarykey;type:int;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	FieldName      *dotnotation.DotNotation `gorm:"column:field_name";type:string` // dot concated path, a.b.c represents a -> b -> c path
	FieldValueType SwagValueType            `gorm:"column:value_type"`
	FieldExample   *SwagValueValue          `gorm:"column:value_example;type:jsonb"`

	ModelSchemaID SchemaID // foreign key to ModelSchema -> ID
}

func (m ModelFieldSchema) TableName() string {
	return "restcol-collections-schematable"
}

type ModelFieldsSchema []*ModelFieldSchema

func (m ModelFieldsSchema) ToJSON(dotPrefixs ...string) ([]byte, error) {

	withPrefix := func(fieldName *dotnotation.DotNotation) string {
		if len(dotPrefixs) == 0 {
			return fieldName.String()
		}
		return fieldName.AddPrefix(dotPrefixs...).String()
	}

	// convert fields into a map with a common dotPrefix
	fieldsMap := make(map[string]string)
	for _, field := range m {
		fieldName := withPrefix(field.FieldName)
		if field.FieldValueType == NumberSwagValueType {
			fieldsMap[fmt.Sprintf("%s.num()", fieldName)] = fmt.Sprintf("%f", field.FieldExample.Interface())
		} else if field.FieldValueType == IntegerSwagValueType {
			fieldsMap[fmt.Sprintf("%s.num()", fieldName)] = fmt.Sprintf("%d", field.FieldExample.Interface())
		} else if field.FieldValueType == BoolSwagValueType {
			fieldsMap[fmt.Sprintf("%s.bool()", fieldName)] = fmt.Sprintf("%t", field.FieldExample.Interface())
		} else {
			fieldsMap[fieldName] = fmt.Sprintf("%v", field.FieldExample.Interface())
		}
	}
	return jsonpath.Marshal(fieldsMap)
}
